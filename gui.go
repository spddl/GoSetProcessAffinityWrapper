package main

import (
	"log"

	"github.com/lxn/walk"
	//lint:ignore ST1001 standard behavior lxn/walk
	. "github.com/lxn/walk/declarative"
)

type MyMainWindow struct {
	*walk.MainWindow
	tv    *walk.TableView
	model *Model
}

func (store *Store) loadGUI() {
	mw := &MyMainWindow{
		model: &Model{items: store.Ifeo},
		tv:    &walk.TableView{},
	}

	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "GoSetProcessAffinityWrapper",
		// MinSize: Size{
		// 	Width:  240,
		// 	Height: 320,
		// },
		Size: Size{
			Width:  550,
			Height: 600,
		},
		Layout: VBox{
			MarginsZero: true,
			SpacingZero: true,
		},
		Children: []Widget{
			Composite{
				Layout: VBox{},
				Children: []Widget{

					TableView{
						OnItemActivated:  mw.lb_ItemActivated,
						Name:             "tableView", // Name is needed for settings persistence
						AlternatingRowBG: true,
						ColumnsOrderable: true,
						Columns: []TableViewColumn{
							{
								Name: "Executable",
							},
							{
								Title: "Priority Class",
								Name:  "PriorityClass",
								FormatFunc: func(value interface{}) string {
									switch value.(int) {
									case 1:
										return "Idle"
									case 2:
										return "Normal"
									case 3:
										return "High"
									case 4:
										return "RealTime"
									case 5:
										return "Below Normal"
									case 6:
										return "Above Normal"
									default:
										return "Normal"
									}
								},
							},
							{
								Title: "I/O Priority",
								Name:  "IoPriority",
								FormatFunc: func(value interface{}) string {
									switch value.(int) {
									case 0:
										return "Very Low"
									case 1:
										return "Low"
									case 2:
										return "Normal"
									case 3:
										return "High"
									default:
										return "Normal"
									}
								},
							},
							{
								Title: "Page Priority",
								Name:  "PagePriority",
								FormatFunc: func(value interface{}) string {
									switch value.(int) {
									case 0:
										return "Idle"
									case 1:
										return "Very Low"
									case 2:
										return "Low"
									case 3:
										return "Medium"
									case 4:
										return "Below Normal"
									case 5:
										return "Normal"
									default:
										return "Normal"
									}
								},
							},
						},
						Model:    mw.model,
						AssignTo: &mw.tv,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text: "Add",
						OnClicked: func() {
							newGame := &Game{
								Executable:    "Game.exe",
								PriorityClass: 2,
								IoPriority:    2,
								PagePriority:  5,
								CPUBits:       CPUMax,
							}
							result, _ := OpenDialog(mw, newGame)
							if result == 1 {
								mw.model.items = append(mw.model.items, *newGame)
								mw.model.PublishRowsReset()
								applySettings(newGame)
							}
						},
					},

					PushButton{
						Text: "Remove",
						OnClicked: func() {
							removeSettings(mw.model.items[mw.tv.CurrentIndex()].Executable)
							mw.model.items = removeGame(mw.model.items, mw.tv.CurrentIndex())
							mw.model.PublishRowsReset()
						},
					},
				},
			},
		},
	}).Create(); err != nil {
		log.Println(err)
		return
	}

	var maxDeviceDesc int
	for i := range store.Ifeo {
		newDeviceDesc := mw.TextWidthSize(store.Ifeo[i].Executable)
		if maxDeviceDesc < newDeviceDesc {
			maxDeviceDesc = newDeviceDesc
		}
	}
	if maxDeviceDesc < 150 {
		mw.tv.Columns().At(0).SetWidth(maxDeviceDesc)
	}

	mw.Show()
	mw.Run()
}

func (mw *MyMainWindow) lb_ItemActivated() {
	newItem := &mw.model.items[mw.tv.CurrentIndex()]
	result, _ := OpenDialog(mw, newItem)
	if result == 1 { // OK
		applySettings(newItem)
	}
}

func (mw *MyMainWindow) TextWidthSize(text string) int {
	canvas, err := (*mw.tv).CreateCanvas()
	if err != nil {
		return 0
	}
	defer canvas.Dispose()

	bounds, _, err := canvas.MeasureTextPixels(text, (*mw.tv).Font(), walk.Rectangle{Width: 9999999}, walk.TextCalcRect)
	if err != nil {
		return 0
	}

	return bounds.Size().Width
}

type Model struct {
	walk.SortedReflectTableModelBase
	items []Game
}

func (m *Model) Items() interface{} {
	return m.items
}

func removeGame(slice []Game, s int) []Game {
	return append(slice[:s], slice[s+1:]...)
}
