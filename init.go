package main

import (
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

// Computer\HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Image File Execution Options\notepad.exe
// Computer\HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\SilentProcessExit\notepad.exe
const (
	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms684863(v=vs.85).aspx
	DEBUG_PROCESS                = 0x00000001
	DEBUG_ONLY_THIS_PROCESS      = 0x00000002
	CREATE_NEW_CONSOLE           = 0x00000010
	CREATE_NEW_PROCESS_GROUP     = 0x00000200
	EXTENDED_STARTUPINFO_PRESENT = 0x00080000

	// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getpriorityclass
	ABOVE_NORMAL_PRIORITY_CLASS uint32 = 0x00008000
	BELOW_NORMAL_PRIORITY_CLASS uint32 = 0x00004000
	HIGH_PRIORITY_CLASS         uint32 = 0x00000080
	IDLE_PRIORITY_CLASS         uint32 = 0x00000040
	NORMAL_PRIORITY_CLASS       uint32 = 0x00000020
	REALTIME_PRIORITY_CLASS     uint32 = 0x00000100
)

var (
	conf Game

	ProcessStart bool

	SystemUser     bool
	executableApp  string
	executablePath string

	CPUMap     map[Bits]string
	CPUArray   []string
	CPUBits    []Bits
	CPUMax     Bits
	initConfig Game
	rnd        = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// rundll32.exe Shell32,Shell_NotifyIconA
// rundll32.exe Shell32,SHTestTokenMembership
const noop = "rundll32.exe Shell32,SHTestTokenMembership"

var PRIORITY_CLASS_Map = map[uint32]string{
	ABOVE_NORMAL_PRIORITY_CLASS: "ABOVE_NORMAL_PRIORITY_CLASS",
	BELOW_NORMAL_PRIORITY_CLASS: "BELOW_NORMAL_PRIORITY_CLASS",
	HIGH_PRIORITY_CLASS:         "HIGH_PRIORITY_CLASS",
	IDLE_PRIORITY_CLASS:         "IDLE_PRIORITY_CLASS",
	NORMAL_PRIORITY_CLASS:       "NORMAL_PRIORITY_CLASS",
	REALTIME_PRIORITY_CLASS:     "REALTIME_PRIORITY_CLASS",
}

var PRIORITY_CLASS = map[int]uint32{
	1: IDLE_PRIORITY_CLASS,
	5: BELOW_NORMAL_PRIORITY_CLASS,
	2: NORMAL_PRIORITY_CLASS,
	6: ABOVE_NORMAL_PRIORITY_CLASS,
	3: HIGH_PRIORITY_CLASS,
	4: REALTIME_PRIORITY_CLASS,
}

func init() {
	CPUMap = make(map[Bits]string, runtime.NumCPU())
	var index Bits = 1
	for i := 0; i < int(runtime.NumCPU()); i++ {
		indexString := strconv.Itoa(i)
		CPUMap[index] = indexString
		CPUArray = append(CPUArray, indexString)
		CPUBits = append(CPUBits, index)
		CPUMax = Set(CPUMax, index)
		index *= 2
	}

	var err error
	executableApp, err = os.Executable()
	if err != nil {
		panic(err)
	}
	executablePath = filepath.Dir(executableApp)

	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile) // https://ispycode.com/GO/Logging/Setting-output-flags
	SystemUser = isSystemUser()

	initConfig = parseFlags(os.Args[0], os.Args[1:])

	log.Println("initConfig", PrettyPrint(initConfig))
}
