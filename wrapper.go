package main

import (
	"flag"
	"path/filepath"
	"strconv"
	"strings"
)

func createWrapperString(ifeo *Game) string {
	var args []string

	if len(ifeo.PreScripts) != 0 || len(ifeo.PostScripts) != 0 {
		SaveConfig(ifeo) // speichert die Config
		if len(ifeo.PostScripts) == 0 {
			removeSilentProcessExitKey(ifeo) // löscht den RegKey
		} else {
			addSilentProcessExitKey(ifeo) // setzt den RegKey
		}
		return executableApp + " -config " + ifeo.Executable + ".toml -exe"
	}

	if ifeo.Ignore { // noop
		return ""
	}

	if ifeo.PriorityClass == 4 { // 4 is Realtime
		args = append(args, "-priorityClass 4")
	}

	if ifeo.IoPriority == 3 { // 3 ist High
		args = append(args, "-ioPriority 3")
	}

	if ifeo.Boost {
		args = append(args, "-boost")
	}

	if ifeo.PassThrough {
		args = append(args, "-passThrough")
	}

	if ifeo.Delay != "" {
		args = append(args, "-delay "+ifeo.Delay)
	}

	if ifeo.CPUBits != CPUMax {
		var CPUArray []string
		for indexCPU, bit := range CPUBits {
			if Has(ifeo.CPUBits, bit) {
				CPUArray = append(CPUArray, strconv.Itoa(indexCPU))
			}
		}
		args = append(args, "-cpu "+strings.Join(CPUArray, ","))
	}

	if len(args) == 0 {
		return ""
	}

	return executableApp + " " + strings.Join(args, " ") + " -exe"
}

func parseFlags(exe string, args []string) Game {
	fs := flag.NewFlagSet(exe, flag.ContinueOnError)

	var exeFlag string
	fs.StringVar(&exeFlag, "exe", "", "")

	var priorityClass int
	fs.IntVar(&priorityClass, "priorityClass", -1, "")

	var ioPriority int
	fs.IntVar(&ioPriority, "ioPriority", -1, "")

	var pagePriority int
	fs.IntVar(&pagePriority, "pagePriority", -1, "")

	var logging bool
	fs.BoolVar(&logging, "logging", false, "")

	var boost bool
	fs.BoolVar(&boost, "boost", false, "")

	var passThrough bool
	fs.BoolVar(&passThrough, "passThrough", false, "")

	var ignore bool
	fs.BoolVar(&ignore, "ignore", false, "")

	var cpuFlag string
	fs.StringVar(&cpuFlag, "cpu", "", "")

	var ProcessID int
	fs.IntVar(&ProcessID, "pid", -1, "")

	var delayFlag string
	fs.StringVar(&delayFlag, "delay", "", "")

	var configFlag string
	fs.StringVar(&configFlag, "config", "", "")

	fs.Parse(args)

	var result Game
	if configFlag == "" {

		if exeFlag == "" && exe != executableApp && exe != filepath.Base(executableApp) {
			exeFlag = exe
		}
		result = Game{
			Executable:    exeFlag,
			PriorityClass: priorityClass,
			IoPriority:    ioPriority,
			PagePriority:  pagePriority,
			Boost:         boost,
			PassThrough:   passThrough,
			Ignore:        ignore,
			ProcessID:     ProcessID,
			Delay:         delayFlag,
			Debugger:      strings.Join(args, " "),
			Config:        configFlag,
		}

	} else {
		if exeFlag != "" {
			ProcessStart = true
		}
		result = ReadConfig(&Game{
			Config:     configFlag,
			ProcessID:  ProcessID,
			Executable: exeFlag,
		})
	}

	var executableArgs string
	if len(fs.Args()) != 0 {
		executableArgs = strings.Join(fs.Args(), " ")
	}
	result.ExecutableArgs = executableArgs

	if cpuFlag != "" {
		for _, val := range strings.Split(cpuFlag, ",") {
			i, err := strconv.Atoi(val)
			if err != nil {
				panic(err)
			}
			result.CPUBits = Set(result.CPUBits, CPUBits[i])
		}
	} else {
		result.CPUBits = CPUMax
	}

	return result
}
