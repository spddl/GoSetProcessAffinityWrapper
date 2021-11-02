package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/lxn/walk"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

type Store struct {
	Ifeo []Game
}

func main() {
	switch {
	case initConfig.Config != "": // ein Config Flag
		conf = ReadConfig(&initConfig)

		if conf.Executable == "" {
			panic("the parameter 'exe' is not allowed to be empty")
		}

	case initConfig.Executable != "" && initConfig.Executable != executableApp: // ein Exe Flag
		conf = initConfig
	default:
		store := new(Store)
		store.findKeys()
		store.findProcessExitKeys()
		store.loadGUI()
		return
	}

	if !SystemUser && ProcessStart {
		var pid uint32
		if conf.PassThrough { // Start Programm as Admin with DEBUG_ONLY_THIS_PROCESS Flag

			var su syscall.StartupInfo
			var pi syscall.ProcessInformation
			su.Cb = uint32(unsafe.Sizeof(su))

			exe, err := syscall.UTF16PtrFromString(conf.Executable)
			if err != nil {
				log.Println(err)
			}

			syscall.CreateProcess(exe, nil, nil, nil, false, NORMAL_PRIORITY_CLASS|DEBUG_ONLY_THIS_PROCESS, nil, nil, &su, &pi)
			err = DebugActiveProcessStop(pi.ProcessId)
			if err != nil {
				log.Println(err)
			}

			pid = pi.ProcessId

			if initConfig.Delay != "" {
				delay, err := time.ParseDuration(initConfig.Delay)
				if err != nil {
					log.Println(err)
					if !SystemUser {
						walk.MsgBox(nil, "Parse Duration Error", err.Error(), walk.MsgBoxIconError|walk.MsgBoxOK)
					}
				}
				log.Printf("time.Sleep(%s)", initConfig.Delay)
				time.Sleep(delay)
			}

		} else { // Start Programm as Admin without IFEO Registry Entry
			k, err := registry.OpenKey(registry.LOCAL_MACHINE, `Software\Microsoft\Windows NT\CurrentVersion\Image File Execution Options\`+filepath.Base(conf.Executable), registry.QUERY_VALUE|registry.SET_VALUE)
			if err != nil {
				panic(err)
			}

			DebuggerValue, _, err := k.GetStringValue("Debugger")
			if err != nil {
				log.Println(err)
			}
			log.Println("Debugger Value:", DebuggerValue)

			if err := k.DeleteValue("Debugger"); err != nil {
				log.Println(err)
			}

			var cmd *exec.Cmd
			if conf.ExecutableArgs == "" {
				log.Printf("exec.Command(%s).Start()\n", conf.Executable)
				cmd = exec.Command(conf.Executable)
			} else {
				log.Printf("exec.Command(%s, %s).Start()\n", conf.Executable, conf.ExecutableArgs)
				cmd = exec.Command(conf.Executable, strings.Split(conf.ExecutableArgs, " ")...)
			}

			err = cmd.Start()
			if err != nil {
				log.Println(err)
			}

			pid = uint32(cmd.Process.Pid)

			if initConfig.Delay != "" {
				delay, err := time.ParseDuration(initConfig.Delay)
				if err != nil {
					log.Println(err)
					if !SystemUser {
						walk.MsgBox(nil, "Parse Duration Error", err.Error(), walk.MsgBoxIconError|walk.MsgBoxOK)
					}
				}
				log.Printf("time.Sleep(%s)", initConfig.Delay)
				time.Sleep(delay)
			}

			if err := k.SetStringValue("Debugger", DebuggerValue); err != nil {
				log.Println(err)
			}

			if err := k.Close(); err != nil {
				log.Println(err)
			}
		}

		i := 0
		for {
			log.Printf("os.FindProcess(%d)", pid)
			_, exist := ProcessPidExist(int(pid))
			if exist {
				log.Println("process gefunden, break")
				break
			}
			time.Sleep(time.Second * 5)
			if i > 60 {
				// 5 minute break
				os.Exit(0)
			}
			i += 1
		}

		t := Task{name: "GoSetProcessAffinity-" + strconv.FormatUint(rnd.Uint64(), 16)} // RandomName
		defer func() {
			for {
				err := t.deleteTask()
				if err == nil {
					break
				}
			}
		}()

		t.createTask(strconv.Itoa(int(pid)), initConfig.Executable)
		t.runTask()

	} else if ProcessStart { // SystemKontext
		if initConfig.ProcessID == 0 {
			log.Println("no PID exit")
			return
		}

		i := 0
		for {
			p, exist := ProcessPidExist(initConfig.ProcessID)
			if exist {
				if p.Pid == initConfig.ProcessID {
					break
				}
			}

			pid, e := ProcessNameExist(filepath.Base(initConfig.Executable))
			if e != nil {
				log.Println(e)
			} else {
				log.Println("pid", pid)
				initConfig.ProcessID = int(pid)
			}

			time.Sleep(time.Second * 5)

			if i > 60 {
				// 5 minute break
				// os.Exit(0)
				return
			}
			i += 1
		}

		pHndl, err := windows.OpenProcess(ProcessSetIinformation|ProcessQueryInformation, false, uint32(initConfig.ProcessID))
		defer windows.CloseHandle(pHndl)
		if err != nil {
			log.Println(err)
		}
		if pHndl == 0 {
			log.Println("no handle")
			return
		}

		if initConfig.Boost {
			SetProcessPriorityBoost(pHndl, false) // If the parameter is FALSE, dynamic boosting is enabled.
		}

		if initConfig.CPUBits != CPUMax {
			err := SetProcessAffinityMask(pHndl, initConfig.CPUBits)
			if err != nil {
				log.Println("setProcessAffinityMask: ", err)
			}
		}

		if initConfig.PriorityClass != -1 {
			if getCpuPriorityClass(pHndl) != PRIORITY_CLASS[initConfig.PriorityClass] {
				err := windows.SetPriorityClass(pHndl, PRIORITY_CLASS[initConfig.PriorityClass])
				if err != nil {
					log.Println("SetPriorityClass: ", err)
				}
			}
		}

		if initConfig.IoPriority != -1 {
			var IoPriorityByte = uint32(initConfig.IoPriority)
			if getProcessIoPriority(pHndl) != IoPriorityByte {
				ntStatus := NtSetInformationProcess(pHndl, ProcessIoPriority, &IoPriorityByte, 4)
				log.Println("NtStatus (ProcessIoPriority)", ntStatus)
			}
		}

		if initConfig.PagePriority != -1 {
			var PagePriorityuint32 = uint32(initConfig.PagePriority)
			if getProcessPagePriority(pHndl) != PagePriorityuint32 {
				ntStatus := NtSetInformationProcess(pHndl, ProcessPagePriority, &PagePriorityuint32, 4)
				log.Println("NtStatus (ProcessPagePriority)", ntStatus)
			}
		}
	}

	var foundSystemPostScripts bool
	for _, script := range conf.PostScripts {
		if script.System {
			foundSystemPostScripts = true
			break
		}
	}

	if !SystemUser && !ProcessStart && foundSystemPostScripts {
		t := Task{name: "GoSetProcessAffinity-" + strconv.FormatUint(rnd.Uint64(), 16)} // RandomName
		defer func() {
			for {
				err := t.deleteTask()
				if err == nil {
					break
				}
			}
		}()

		t.createTask(strconv.Itoa(0), "")
		t.runTask()
	}

	if ProcessStart {
		startScripts(conf.PreScripts)
	} else {
		startScripts(conf.PostScripts)
	}
}

func startScripts(scripts []Scripts) {
	for _, script := range scripts {
		if script.System != SystemUser {
			continue
		}
		var cmd_instance *exec.Cmd
		if script.Args == "" {
			log.Printf("exec.Command(%s)\n", script.Name)
			cmd_instance = exec.Command(script.Name)
		} else {
			log.Printf("exec.Command(%s, %s)\n", script.Name, script.Args)
			cmd_instance = exec.Command(script.Name, strings.Split(script.Args, " ")...)
		}
		if script.HideWindow {
			cmd_instance.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		}

		cmd_instance.Start()
		// output, err := cmd_instance.CombinedOutput()
		// if err != nil {
		// 	log.Println(err)
		// }
		// log.Println("output", string(output))
	}
}

func ProcessPidExist(pid int) (*os.Process, bool) {
	p, err := os.FindProcess(pid)
	return p, err == nil
}

const processEntrySize = 568 // unsafe.Sizeof(windows.ProcessEntry32{})

func ProcessNameExist(name string) (uint32, error) {
	h, e := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if e != nil {
		return 0, e
	}
	p := windows.ProcessEntry32{Size: processEntrySize}
	for {
		e := windows.Process32Next(h, &p)
		if e != nil {
			return 0, e
		}
		if windows.UTF16ToString(p.ExeFile[:]) == name {
			return p.ProcessID, nil
		}
	}
}

func PrettyPrint(data interface{}) string {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return ""
	}
	return string(val)
}
