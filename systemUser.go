package main

import (
	"log"
	"os/user"
)

func isSystemUser() bool {
	userCurrent, err := user.Current()
	if err != nil {
		log.Panic(err)
	}
	log.Println("user.Uid", userCurrent.Uid)
	return userCurrent.Uid == "S-1-5-18"
}

// func amIAdmin() bool {
// 	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
// 	admin := err == nil
// 	return admin
// }

// func runMeAdmin() {
// 	verb := "runas"
// 	exe, _ := os.Executable()
// 	args := strings.Join(os.Args[1:], " ")
// 	cwd, _ := os.Getwd()

// 	verbPtr, _ := syscall.UTF16PtrFromString(verb)
// 	exePtr, _ := syscall.UTF16PtrFromString(exe)
// 	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
// 	argPtr, _ := syscall.UTF16PtrFromString(args)

// 	var showCmd int32 = 1 // SW_NORMAL

// 	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
