package main

import (
	"log"

	"github.com/lxn/walk"

	//lint:ignore ST1001 standard behavior lxn/walk
	. "github.com/lxn/walk/declarative"
)

func dialogScriptsCommands(owner walk.Form, data Scripts) (int, Scripts) {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var ProgramLE, ArgumentsLE *walk.LineEdit
	var HideCB, SystemCB *walk.CheckBox

	err := Dialog{
		AssignTo:      &dlg,
		Title:         "",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize: Size{
			Width:  700,
			Height: 250,
		},
		Layout: VBox{},
		Children: []Widget{
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text: "Powershell",
						OnClicked: func() {
							ProgramLE.SetText("Powershell")
							ArgumentsLE.SetText("-Command &{ Add-Content -Path " + executablePath + "\\logging.log -Value \"start $(Get-Date)\" }")
						},
					},
					PushButton{
						Text: "PowerCfg (High Performance)",
						OnClicked: func() {
							ProgramLE.SetText("powercfg")
							ArgumentsLE.SetText("/setactive 8c5e7fda-e8bf-4a96-9a85-a6e23a8c635c")
						},
					},
					PushButton{
						Text: "PowerCfg (Balanced)",
						OnClicked: func() {
							ProgramLE.SetText("powercfg")
							ArgumentsLE.SetText("/setactive 381b4222-f694-41f0-9685-ff5bb260df2e")
						},
					},
					HSpacer{},
				},
			},

			Composite{
				Layout: Grid{
					Columns: 2,
				},
				Children: []Widget{
					Label{
						Text: "Program",
					},
					LineEdit{
						AssignTo: &ProgramLE,
						Text:     data.Name,
						OnTextChanged: func() {
							data.Name = ProgramLE.Text()
						},
					},

					Label{
						Text: "Arguments",
					},
					LineEdit{
						AssignTo: &ArgumentsLE,
						Text:     data.Args,
						OnTextChanged: func() {
							data.Args = ArgumentsLE.Text()
						},
					},

					CheckBox{
						ColumnSpan: 2,
						AssignTo:   &HideCB,
						Text:       "Hide Window",
						Checked:    data.HideWindow,
						OnClicked: func() {
							data.HideWindow = HideCB.Checked()
						},
					},
					CheckBox{
						ColumnSpan: 2,
						AssignTo:   &SystemCB,
						Text:       "Start as System",
						Checked:    data.System,
						OnClicked: func() {
							data.System = SystemCB.Checked()
						},
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							dlg.Accept()
						},
					},
					PushButton{
						AssignTo: &cancelPB,
						Text:     "Cancel",
						OnClicked: func() {
							dlg.Cancel()
						},
					},
					HSpacer{},
				},
			},
		},
	}.Create(owner)
	if err != nil {
		log.Println(err)
	}
	return dlg.Run(), data
}
