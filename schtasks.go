package main

import (
	"log"
	"os/exec"
)

type Task struct {
	name string
}

func (t *Task) createTask(pid, exe string) {
	out, err := exec.Command("schtasks", "/create", "/tn", t.name, "/tr", executableApp+" -pid "+pid+" -exe "+exe, "/sc", "ONCE", "/sd", "01/01/2337", "/st", "00:00", "/ru", "SYSTEM", "/rl", "HIGHEST", "/f").CombinedOutput()
	if err != nil {
		log.Printf("createTask() err %s\n", out)
	}
}

func (t *Task) runTask() {
	out, err := exec.Command("schtasks", "/run", "/tn", t.name).CombinedOutput()
	if err != nil {
		log.Printf("runTask() err %+v %s\n", err, out)
	}
	log.Printf("runTask() %+v %s\n", err, out)
}

func (t *Task) deleteTask() error {
	cmd := exec.Command("schtasks", "/delete", "/tn", t.name, "/f")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("deleteTask() err %+v %s\n", err, out)
	}
	return nil
}
