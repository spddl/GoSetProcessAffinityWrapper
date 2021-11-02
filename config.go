package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/lxn/walk"
)

type Game struct {
	Boost          bool      `toml:"boost"`
	PassThrough    bool      `toml:"passthrough"`
	Ignore         bool      `toml:"ignore"` // Maybe obsolete
	PriorityClass  int       `toml:"priorityClass"`
	IoPriority     int       `toml:"ioPriority"`
	PagePriority   int       `toml:"pagePriority "`
	CPU            string    `toml:"cpu"`
	CPUBits        Bits      `toml:"-"`
	Executable     string    `toml:"exe"`
	ExecutableArgs string    `toml:"-"`
	Debugger       string    `toml:"-"`
	ProcessID      int       `toml:"-"`
	Delay          string    `toml:"delay"`
	Config         string    `toml:"-"`
	PreScripts     []Scripts `toml:"preScripts"`
	PostScripts    []Scripts `toml:"postScripts"`
}

type Scripts struct {
	Name       string `toml:"name"`
	Args       string `toml:"args"`
	HideWindow bool   `toml:"hideWindow,omitempty"`
	System     bool   `toml:"system,omitempty"`
}

func ReadConfig(ifeo *Game) Game {
	tomlData, err := os.ReadFile(filepath.Join(executablePath, ifeo.Config))
	if err != nil {
		log.Println(err)
	}

	if _, err := toml.Decode(string(tomlData), &ifeo); err != nil {
		log.Println(err)
		if !SystemUser {
			walk.MsgBox(nil, "toml Decode Error", err.Error(), walk.MsgBoxIconError|walk.MsgBoxOK)
		}
	}

	return *ifeo
}

func SaveConfig(ifeo *Game) {
	f, err := os.Create(filepath.Join(executablePath, ifeo.Executable+".toml"))
	if err != nil {
		// failed to create/open the file
		log.Fatal(err)
	}
	if err := toml.NewEncoder(f).Encode(ifeo); err != nil {
		// failed to encode
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		// failed to close the file
		log.Fatal(err)
	}
}
