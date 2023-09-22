package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dmitry-fofanov/title_kitten/icon"
	appStrings "github.com/dmitry-fofanov/title_kitten/strings"
)

type selectableText struct {
	widget.Entry;
	// text *widget.RichText;
	// dirty bool;
}

func newSelectableText() *selectableText {
	textObject := &selectableText{}
	textObject.ExtendBaseWidget(textObject)
	textObject.Wrapping = fyne.TextWrapOff
	// text.Wrapping = fyne.TextWrapWord
	// text.MultiLine = true
	textObject.SetMinRowsVisible(1)
	// textObject.Scroll = container.ScrollNone

	return textObject
}

func (st *selectableText) TypedRune(r rune) {}

func setUpGui() {
	myApp := app.New()
	myApp.SetIcon(icon.AppIcon)
	mainWindow := myApp.NewWindow(appStrings.AppName + " " + appStrings.AppVersion)

	statusLabel = widget.NewLabel(appStrings.AppName + "!")
	ignoreCaseCheckbox := widget.NewCheck(
		appStrings.IgnoreCase,
		func(state bool) {
			ignoreCase = state
		},
	)
	progressBar = widget.NewProgressBar()

	// resultsWidget = widget.NewList(
	// 	func() int {
	// 		return len(resultsData)
	// 	},
	// 	func() fyne.CanvasObject {
	// 		return newSelectableText()
	// 	},
	// 	func(i widget.ListItemID, o fyne.CanvasObject) {
	// 		o.(*selectableText).SetText(resultsData[i])
	// 	})
	// resultsWidget = widget.NewTable(
	// 	func() (int, int) {
	// 		return len(resultsData), len(resultsData[0])
	// 	},
	// 	func() fyne.CanvasObject {
	// 		return newSelectableText()
	// 	},
	// 	func(i widget.TableCellID, o fyne.CanvasObject) {
	// 		o.(*selectableText).SetText(resultsData[i.Row][i.Col])
	// 	})
	// resultsWidget.SetColumnWidth(0, 50)
	// resultsWidget.SetColumnWidth(1, 100)
	// resultsWidget.SetColumnWidth(2, 100)
	// resultsWidget.SetColumnWidth(3, 2000)
	resultsWidget = widget.NewTree(
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			return resultsData[id]
		},
		func(id widget.TreeNodeID) bool {
			_, isBranch := resultsData[id]
			return isBranch
		},
		func(_ bool) fyne.CanvasObject {
			return newSelectableText()
		},
		func(id widget.TreeNodeID, _ bool, object fyne.CanvasObject) {
			object.(*selectableText).SetText(id)
		},
	)

	content := container.NewBorder(
		container.NewVBox(
			statusLabel,
			container.NewBorder(
				nil,
				nil,
				container.NewVBox(
					widget.NewLabel(appStrings.FilePathLabel),
					widget.NewLabel(appStrings.QueryLabel),
				),
				container.NewVBox(
					widget.NewButton(
						appStrings.FileButtonLabel,
						func() { go fileOpenDialog() },
					),
					ignoreCaseCheckbox,
				),
				container.NewVBox(
					filePathInput,
					queryInput,
				),
			),
			widget.NewButton(appStrings.SearchButtonLabel, func() { go mainSearchFunction() }),
			progressBar,
			&widget.Separator{},
		),
		nil,
		nil,
		nil,
		resultsWidget,
	)

	mainWindow.SetContent(content)

	mainWindow.Resize(fyne.Size{Height: 675, Width: 1200})

	mainWindow.ShowAndRun()
}
