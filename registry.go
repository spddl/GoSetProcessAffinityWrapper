package main

import (
	"log"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const ifeoPath = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Image File Execution Options`
const silentProcessExitPath = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\SilentProcessExit`

type Bits uint64

type IFEO struct {
	Executable       string
	Debugger         string
	CpuPriorityClass int
	IoPriority       int
	PagePriority     int
	CPUBits          Bits
	Boost            bool
	PassThrough      bool
	Ignore           bool
}

func (store *Store) findKeys() {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, ifeoPath, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		log.Println("Can't open registry key", err)
		return
	}
	defer k.Close()

	params, err := k.ReadSubKeyNames(-1)
	if err != nil {
		log.Printf("Can't ReadSubKeyNames %s\n", err)
		return
	}

	for _, param := range params {
		ifeo := Game{
			Executable:    param,
			PriorityClass: 2,
			IoPriority:    2,
			PagePriority:  5,
			CPUBits:       CPUMax,
		}

		subPerfOptions, err := registry.OpenKey(registry.LOCAL_MACHINE, ifeoPath+"\\"+param+"\\PerfOptions", registry.QUERY_VALUE)
		if err == nil {
			valCpuPriorityClass, _, err := subPerfOptions.GetIntegerValue("CpuPriorityClass")
			if err == nil {
				ifeo.PriorityClass = int(valCpuPriorityClass)
			}

			valIoPriority, _, err := subPerfOptions.GetIntegerValue("IoPriority")
			if err == nil {
				ifeo.IoPriority = int(valIoPriority)
			}

			valPagePriority, _, err := subPerfOptions.GetIntegerValue("PagePriority")
			if err == nil {
				ifeo.PagePriority = int(valPagePriority)
			}
		}
		subPerfOptions.Close()

		subk, err := registry.OpenKey(registry.LOCAL_MACHINE, ifeoPath+"\\"+param, registry.QUERY_VALUE)
		if err == nil {
			valDebugger, _, err := subk.GetStringValue("Debugger")
			if err == nil {
				if strings.HasPrefix(valDebugger, executableApp) {
					resulkt := parseFlags(param, strings.Split(valDebugger, " ")[1:])

					if ifeo.PriorityClass != resulkt.PriorityClass {
						log.Println("WARN PriorityClass", ifeo.PriorityClass, resulkt.PriorityClass)
						resulkt.PriorityClass = ifeo.PriorityClass
					}
					if ifeo.IoPriority != resulkt.IoPriority {
						log.Println("WARN IoPriority", ifeo.IoPriority, resulkt.IoPriority)
						resulkt.IoPriority = ifeo.IoPriority
					}
					if ifeo.PagePriority != resulkt.PagePriority {
						log.Println("WARN PagePriority", ifeo.PagePriority, resulkt.PagePriority)
						resulkt.PagePriority = ifeo.PagePriority
					}
					ifeo = resulkt
				} else {
					ifeo.Debugger = valDebugger
				}
			}
		}
		subk.Close()

		if ifeo.Debugger == "" && ifeo.PriorityClass == 2 && ifeo.IoPriority == 2 && ifeo.PagePriority == 5 {
			continue
		}
		store.Ifeo = append(store.Ifeo, ifeo)
	}
}

// find executables that only load the config on exit
func (store *Store) findProcessExitKeys() {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, silentProcessExitPath, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		log.Println("Can't open registry key", err)
		return
	}
	defer k.Close()

	params, err := k.ReadSubKeyNames(-1)
	if err != nil {
		log.Printf("Can't ReadSubKeyNames %s\n", err)
		return
	}

	for _, param := range params {
		subk, err := registry.OpenKey(registry.LOCAL_MACHINE, silentProcessExitPath+"\\"+param, registry.QUERY_VALUE)
		if err == nil {
			valMonitorProcess, _, err := subk.GetStringValue("MonitorProcess")
			if err != nil {
				subk.Close()
				continue
			}

			valReportingMode, _, err := subk.GetIntegerValue("ReportingMode")
			if err != nil {
				subk.Close()
				continue
			}

			if valMonitorProcess != "" && valReportingMode == 1 {
				var foundExecutable bool
				for i := 0; i < len(store.Ifeo); i++ {
					if store.Ifeo[i].Executable == param {
						foundExecutable = true
						subk.Close()
						break
					}
				}
				if !foundExecutable {
					store.Ifeo = append(store.Ifeo, Game{
						Executable:    param,
						PriorityClass: 2,
						IoPriority:    2,
						PagePriority:  5,
						CPUBits:       CPUMax,
					})
				}
			}
		}
	}
}

// remove the IFEO settings
func removeSettings(exe string) {
	subk, err := registry.OpenKey(registry.LOCAL_MACHINE, ifeoPath+"\\"+exe, registry.WRITE)
	if err == nil {
		subk.DeleteValue("Debugger")
	}
	subk.Close()

	subPerfOptions, err := registry.OpenKey(registry.LOCAL_MACHINE, ifeoPath+"\\"+exe+"\\PerfOptions", registry.WRITE)
	if err == nil {
		subPerfOptions.DeleteValue("CpuPriorityClass")
		subPerfOptions.DeleteValue("IoPriority")
		subPerfOptions.DeleteValue("PagePriority")

	}
	subPerfOptions.Close()
}

func applySettings(ifeo *Game) {
	subk, _, err := registry.CreateKey(registry.LOCAL_MACHINE, ifeoPath+"\\"+ifeo.Executable, registry.WRITE)
	if err == nil {
		if ifeo.Debugger != "" {
			errUnwrap(subk.SetStringValue("Debugger", ifeo.Debugger))
		} else {
			subk.DeleteValue("Debugger")
		}
	} else {
		log.Println(err)
	}
	subk.Close()

	subPerfOptions, _, err := registry.CreateKey(registry.LOCAL_MACHINE, ifeoPath+"\\"+ifeo.Executable+"\\PerfOptions", registry.ALL_ACCESS)
	if err == nil {

		if ifeo.PriorityClass == 2 { // Normal
			subPerfOptions.DeleteValue("CpuPriorityClass")
		} else {
			valCpuPriorityClass, _, err := subPerfOptions.GetIntegerValue("CpuPriorityClass")
			if err == nil {
				if int(valCpuPriorityClass) != ifeo.PriorityClass {
					errUnwrap(subPerfOptions.SetDWordValue("CpuPriorityClass", uint32(ifeo.PriorityClass)))
				}
			} else { // value not found
				errUnwrap(subPerfOptions.SetDWordValue("CpuPriorityClass", uint32(ifeo.PriorityClass)))
			}
		}

		if ifeo.IoPriority == 2 { // Normal
			subPerfOptions.DeleteValue("IoPriority")
		} else {
			valIoPriority, _, err := subPerfOptions.GetIntegerValue("IoPriority")
			if err == nil {
				if int(valIoPriority) != ifeo.IoPriority {
					errUnwrap(subPerfOptions.SetDWordValue("IoPriority", uint32(ifeo.IoPriority)))
				}
			} else {
				errUnwrap(subPerfOptions.SetDWordValue("IoPriority", uint32(ifeo.IoPriority)))
			}
		}

		if ifeo.PagePriority == 5 { // Normal
			subPerfOptions.DeleteValue("PagePriority")
		} else {
			valPagePriority, _, err := subPerfOptions.GetIntegerValue("PagePriority")
			if err == nil {
				if int(valPagePriority) != ifeo.PagePriority {
					errUnwrap(subPerfOptions.SetDWordValue("PagePriority", uint32(ifeo.PagePriority)))
				}
			} else {
				errUnwrap(subPerfOptions.SetDWordValue("PagePriority", uint32(ifeo.PagePriority)))
			}
		}

	} else {
		log.Println(err)
	}
	subPerfOptions.Close()
}

func addSilentProcessExitKey(game *Game) {
	executable := filepath.Base(game.Executable)
	log.Println("executable", executable)

	r, _, err := registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\SilentProcessExit\`+executable, registry.WRITE)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	errUnwrap(r.SetDWordValue("ReportingMode", 1))
	errUnwrap(r.SetStringValue("MonitorProcess", executableApp+" -config "+executable+".toml"))
}

func removeSilentProcessExitKey(game *Game) {
	executable := filepath.Base(game.Executable)
	log.Println("executable", executable)

	r, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\SilentProcessExit\`+executable, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return
	}
	defer r.Close()

	r.DeleteValue("ReportingMode")
	r.DeleteValue("MonitorProcess")
}

func errUnwrap(err error) {
	if err != nil {
		panic(err)
	}
}
