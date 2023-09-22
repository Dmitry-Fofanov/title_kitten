package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2/widget"
	"github.com/dmitry-fofanov/title_kitten/parser"
	appStrings "github.com/dmitry-fofanov/title_kitten/strings"
	"github.com/ncruces/zenity"
)

const (
	numberOfWorkers = 16
)

var (
	// resultsData = []string{}
	// resultsData = [][4]string{}
	resultsData = make(map[string][]string)

	statusLabel   *widget.Label
	filePathInput = widget.NewEntry()
	queryInput    = widget.NewEntry()
	progressBar   *widget.ProgressBar
	// resultsWidget *widget.List
	// resultsWidget *widget.Table
	resultsWidget *widget.Tree

	ignoreCase         = false
	isRunning          = false
	isInFileDialog     = false
	supportedFiletypes = map[string]bool{
		".ass": true,
		".ssa": true,
		".srt": true,
	}
)

func fileOpenDialog() {
	if isInFileDialog {
		return
	}
	isInFileDialog = true
	newPath, err := zenity.SelectFile(
		zenity.Filename(filePathInput.Text),
		zenity.Directory(),
		zenity.Title(appStrings.AppName),
	)
	isInFileDialog = false
	if err != nil {
		return
	}
	filePathInput.Text = newPath
	filePathInput.Refresh()
}

func searchWorker(query string, ignoreCase bool, files <-chan string, results chan<- []string) {
	var fileData []string
	for filename := range files {
		fileData = []string{}
		subtitles, err := parser.ParseSub(filename)
		if err != nil {
			results <- fileData
			continue
		}
		for _, event := range subtitles {
			text := event.Text
			if ignoreCase {
				text = strings.ToLower(text)
			}
			if strings.Contains(text, query) {
				fileData = append(
					fileData,
					event.Format(),
				)
			}
		}
		if len(fileData) > 0 {
			fileData = append([]string{filename}, fileData...)
			// fileData = append(fileData, [4]string{"", "", "", ""})
		}
		results <- fileData
	}
}

func mainSearchFunction() {
	if isRunning {
		zenity.Warning(
			appStrings.BusyWarning,
			zenity.Title(appStrings.WarningTitle),
			zenity.WarningIcon,
		)
		return
	}

	isRunning = true
	defer func() { isRunning = false }()

	query := queryInput.Text

	if query == "" {
		zenity.Error(
			appStrings.EmptyQueryError,
			zenity.Title(appStrings.ErrorTitle),
			zenity.ErrorIcon,
		)
		return
	}

	currentIgnoreCase := ignoreCase
	if currentIgnoreCase {
		query = strings.ToLower(query)
	}
	root := filePathInput.Text
	rootInfo, err := os.Stat(root)

	if err != nil {
		zenity.Error(
			appStrings.BrokenPathError,
			zenity.Title(appStrings.ErrorTitle),
			zenity.ErrorIcon,
		)
		return
	}

	var (
		fileList  []string
		fileCount int
		increment float64
	)

	if rootInfo.IsDir() {
		filepath.WalkDir(
			root,
			func(path string, d fs.DirEntry, _ error) error {
				if supportedFiletypes[strings.ToLower(filepath.Ext(d.Name()))] {
					fileList = append(fileList, path)
					fileCount++
				}
				return nil
			},
		)
		increment = 1.0 / float64(fileCount)
	} else {
		fileCount = 1
		fileList = []string{root}
		increment = 1.0
	}

	if fileCount == 0 {
		zenity.Error(
			appStrings.NoFilesError,
			zenity.Title(appStrings.ErrorTitle),
			zenity.ErrorIcon,
		)
		statusLabel.SetText(appStrings.NoFilesText)
		return
	}

	statusLabel.SetText(appStrings.SearchingText)

	progressBar.SetValue(0)
	progressBar.Refresh()
	// resultsData = []string{}
	resultsData = make(map[string][]string)
	resultsWidget.Refresh()

	files := make(chan string, fileCount)
	results := make(chan []string, fileCount)

	startTime := time.Now()

	for w := 1; w <= numberOfWorkers; w++ {
		go searchWorker(query, currentIgnoreCase, files, results)
	}

	for _, filename := range fileList {
		files <- filename
	}

	for i := 1; i <= fileCount; i++ {
		fileResults := <-results
		if len(fileResults) > 0 {
			// resultsData = append(resultsData, fileResults...)
			resultsData[fileResults[0]] = fileResults[1:]
			resultsData[""] = append(resultsData[""], fileResults[0])
			resultsWidget.Refresh()
		}
		progressBar.SetValue(progressBar.Value + increment)
		progressBar.Refresh()
	}

	close(files)

	processDuration := time.Since(startTime)
	statusLabel.SetText(appStrings.FinishedText + processDuration.String())
}

func main() {
	setUpGui()
}
