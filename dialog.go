package main

import (
	"github.com/lxn/walk"
	//lint:ignore ST1001 standard behavior lxn/walk
	. "github.com/lxn/walk/declarative"
)

type chain struct {
	Value int
	Name  string
}

func CpuPriorityClassCombo() []*chain {
	return []*chain{
		{1, "Idle"},
		{2, "Normal"},
		{3, "High"},
		{4, "RealTime"},
		{5, "Below Normal"},
		{6, "Above Normal"},
	}
}

func IoPriorityCombo() []*chain {
	return []*chain{
		{0, "Very Low"},
		{1, "Low"},
		{2, "Normal"},
		{3, "High"},
	}
}

func PagePriorityCombo() []*chain {
	return []*chain{
		{0, "Idle"},
		{1, "Very Low"},
		{2, "Low"},
		{3, "Medium"},
		{4, "Below Normal"},
		{5, "Normal"},
	}
}

var AdminLabel *walk.Label

type MyDialog struct {
	*walk.Dialog
	startScripts *ListModel
	endScripts   *ListModel
	startLB      *walk.ListBox
	endLB        *walk.ListBox
}

type ListModel struct {
	walk.ListModelBase
	items []Scripts
}

func OpenDialog(owner walk.Form, ifeo *Game) (int, error) {
	dlg := &MyDialog{
		startScripts: &ListModel{items: ifeo.PreScripts},
		endScripts:   &ListModel{items: ifeo.PostScripts},
	}

	var ExecutableLE, DelayLE *walk.LineEdit
	var acceptPB, cancelPB *walk.PushButton
	var ignoreCB, passthroughCB *walk.CheckBox
	var CPUArrayComposite, StartScriptsComposite, EndScriptsComposite *walk.Composite
	var cpuPriorityClassCB, ioPriorityCB, pagePriorityCB *walk.ComboBox

	return Dialog{
		AssignTo:      &dlg.Dialog,
		Title:         "Edit Executable",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize: Size{
			Width:  300,
			Height: 300,
		},
		Layout: VBox{},
		Children: []Widget{

			GroupBox{
				Title:  "Process Settings",
				Layout: Grid{Columns: 1},
				Children: []Widget{

					Composite{
						Layout: Grid{Columns: 2},
						Children: []Widget{
							Composite{
								Layout: Grid{
									Columns: 2,
								},
								Children: []Widget{
									Label{
										Text: "Executable:",
									},
									LineEdit{
										AssignTo:  &ExecutableLE,
										CueBanner: "Game.exe",
										Text:      ifeo.Executable,
										OnTextChanged: func() {
											ifeo.Executable = ExecutableLE.Text()
										},
									},

									Label{
										Text: "Cpu Priority Class:",
									},
									ComboBox{
										AssignTo:      &cpuPriorityClassCB,
										Value:         ifeo.PriorityClass,
										BindingMember: "Value",
										DisplayMember: "Name",
										Model:         CpuPriorityClassCombo(),
										OnCurrentIndexChanged: func() { // BUG: calls 3 times
											ifeo.PriorityClass = cpuPriorityClassCB.CurrentIndex() + 1 // Start at 1
											checkAdminRights(ifeo)
										},
									},

									Label{
										Text: "I/O Priority:",
									},
									ComboBox{
										AssignTo:      &ioPriorityCB,
										Value:         ifeo.IoPriority,
										BindingMember: "Value",
										DisplayMember: "Name",
										Model:         IoPriorityCombo(),
										OnCurrentIndexChanged: func() { // BUG: calls 3 times
											ifeo.IoPriority = ioPriorityCB.CurrentIndex()
											checkAdminRights(ifeo)
										},
									},

									Label{
										Text: "Page Priority:",
									},
									ComboBox{
										AssignTo:      &pagePriorityCB,
										Value:         ifeo.PagePriority,
										BindingMember: "Value",
										DisplayMember: "Name",
										Model:         PagePriorityCombo(),
										OnCurrentIndexChanged: func() { // BUG: calls 3 times
											ifeo.PagePriority = pagePriorityCB.CurrentIndex()
											checkAdminRights(ifeo)
										},
									},

									Label{
										Text: "Delay:",
									},
									LineEdit{
										AssignTo:  &DelayLE,
										CueBanner: "1m10s", // https://pkg.go.dev/time#ParseDuration
										Text:      ifeo.Delay,
										OnTextChanged: func() {
											ifeo.Delay = DelayLE.Text()
											checkAdminRights(ifeo)
										},
									},

									CheckBox{
										AssignTo:       &ignoreCB,
										Name:           "Ignore",
										Text:           "Block Process:",
										TextOnLeftSide: true,
										ColumnSpan:     2,
										Checked:        ifeo.Debugger == `"`+noop+`"`,
										OnClicked: func() {
											if ignoreCB.Checked() {
												ifeo.Debugger = noop

												cpuPriorityClassCB.SetEnabled(false)
												ioPriorityCB.SetEnabled(false)
												pagePriorityCB.SetEnabled(false)
												DelayLE.SetEnabled(false)
												passthroughCB.SetEnabled(false)
												CPUArrayComposite.SetEnabled(false)
												StartScriptsComposite.SetEnabled(false)
												EndScriptsComposite.SetEnabled(false)
											} else {
												ifeo.Debugger = ""

												cpuPriorityClassCB.SetEnabled(true)
												ioPriorityCB.SetEnabled(true)
												pagePriorityCB.SetEnabled(true)
												DelayLE.SetEnabled(true)
												passthroughCB.SetEnabled(true)
												CPUArrayComposite.SetEnabled(true)
												StartScriptsComposite.SetEnabled(true)
												EndScriptsComposite.SetEnabled(true)
											}
											checkAdminRights(ifeo)
										},
									},

									CheckBox{
										AssignTo:       &passthroughCB,
										Name:           "PassThrough",
										ToolTipText:    "don't use it on games with anticheat",
										Text:           "Pass through:",
										ColumnSpan:     2,
										TextOnLeftSide: true,
										Checked:        ifeo.PassThrough,
										OnClicked: func() {
											if passthroughCB.Checked() {
												ifeo.PassThrough = true
											} else {
												ifeo.PassThrough = false
											}
											checkAdminRights(ifeo)
										},
									},
								},
							},
						},
					},
					Composite{
						AssignTo:  &CPUArrayComposite,
						Alignment: AlignHCenterVCenter,
						Layout:    Grid{Columns: 2},
						Children:  CheckBoxList(CPUArray, ifeo),
					},
				},
			},

			Composite{
				AssignTo: &StartScriptsComposite,
				Layout: Grid{
					Columns:     3,
					MarginsZero: true,
				},
				Children: []Widget{
					Label{
						Text:       "Start scripts:",
						ColumnSpan: 3,
					},
					ListBox{
						AssignTo:   &dlg.startLB,
						ColumnSpan: 3,
						Model:      dlg.startScripts,
						OnItemActivated: func() {
							index := dlg.startLB.CurrentIndex()
							if index < 0 {
								return
							}
							result, data := dialogScriptsCommands(dlg.Dialog, dlg.startScripts.items[index])
							if result == 1 {
								dlg.startScripts.items[index] = data
								dlg.startScripts.ListModelBase.PublishItemsReset()

								ifeo.PreScripts[index] = data
							}

						},
					},
					PushButton{
						Text: "Add",
						OnClicked: func() {
							result, data := dialogScriptsCommands(dlg.Dialog, Scripts{}) // 1 => OK
							if result == 1 {
								dlg.startScripts.items = append(dlg.startScripts.items, data)
								dlg.startScripts.ListModelBase.PublishItemsReset()

								ifeo.PreScripts = append(ifeo.PreScripts, data)
							}

						},
					},
					PushButton{
						Text: "Remove",
						OnClicked: func() {
							index := dlg.startLB.CurrentIndex()
							if index < 0 {
								return
							}
							dlg.startScripts.items = removeScript(dlg.startScripts.items, index)
							dlg.startScripts.ListModelBase.PublishItemsReset()

							ifeo.PreScripts = removeScript(ifeo.PreScripts, index)
						},
					},

					VSpacer{},
				},
			},

			Composite{
				AssignTo: &EndScriptsComposite,
				Layout: Grid{
					Columns:     3,
					MarginsZero: true,
				},
				Children: []Widget{

					Label{
						Text:       "End scripts:",
						ColumnSpan: 3,
					},
					ListBox{
						AssignTo:   &dlg.endLB,
						ColumnSpan: 3,
						Model:      dlg.endScripts,
						OnItemActivated: func() {
							index := dlg.endLB.CurrentIndex()
							if index < 0 {
								return
							}
							result, data := dialogScriptsCommands(dlg.Dialog, dlg.endScripts.items[index])
							if result == 1 {
								dlg.endScripts.items[index] = data
								dlg.endScripts.ListModelBase.PublishItemsReset()

								ifeo.PostScripts[index] = data
							}
						},
					},
					PushButton{
						Text: "Add",
						OnClicked: func() {
							result, data := dialogScriptsCommands(dlg.Dialog, Scripts{}) // 1 => OK
							if result == 1 {
								dlg.endScripts.items = append(dlg.endScripts.items, data)
								dlg.endScripts.ListModelBase.PublishItemsReset()

								ifeo.PostScripts = append(ifeo.PostScripts, data)
							}
						},
					},
					PushButton{
						Text: "Remove",
						OnClicked: func() {
							index := dlg.endLB.CurrentIndex()
							if index < 0 {
								return
							}
							dlg.endScripts.items = removeScript(dlg.endScripts.items, index)
							dlg.endScripts.ListModelBase.PublishItemsReset()

							ifeo.PostScripts = removeScript(ifeo.PostScripts, index)
						},
					},
					VSpacer{},
				},
			},

			HSpacer{},
			Label{
				AssignTo:           &AdminLabel,
				Text:               "Admin rights required to start " + ifeo.Executable,
				TextColor:          walk.RGB(0xff, 0, 0),
				Background:         SolidColorBrush{Color: walk.RGB(0xf0, 0xf0, 0xf0)},
				AlwaysConsumeSpace: true,
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							if ifeo.Debugger != noop {
								ifeo.Debugger = createWrapperString(ifeo)
							}
							dlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)
}

func CheckBoxList(names []string, ifeo *Game) []Widget {
	bs := make([]*walk.CheckBox, len(names))
	children := []Widget{}
	for i, name := range names {
		bs[i] = new(walk.CheckBox)
		local_CPUBits := CPUBits[i]
		children = append(children, CheckBox{
			AssignTo: &bs[i],
			Text:     "CPU " + name,
			Checked:  Has(ifeo.CPUBits, local_CPUBits),
			OnClicked: func() {
				ifeo.CPUBits = Toggle(ifeo.CPUBits, local_CPUBits)
				checkAdminRights(ifeo)
			},
		})
	}
	return children
}

func checkAdminRights(ifeo *Game) {
	if createWrapperString(ifeo) == "" {
		AdminLabel.SetVisible(false)
	} else {
		AdminLabel.SetVisible(true)
	}
}

func removeScript(slice []Scripts, s int) []Scripts {
	return append(slice[:s], slice[s+1:]...)
}

func (m *ListModel) ItemCount() int {
	return len(m.items)
}

func (m *ListModel) Value(index int) interface{} {
	if m.items[index].Args == "" {
		return m.items[index].Name
	}
	return m.items[index].Name + " " + m.items[index].Args
}
