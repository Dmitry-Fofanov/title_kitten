package parser

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

func ParseAssEvent(line string) (SubLine, error) {
	splitLine := strings.SplitN(line, ",", 10)
	if len(splitLine) < 10 {
		return SubLine{}, errors.New("Bad event string")
	}
	return SubLine{
		StartTime: splitLine[1],
		EndTime: splitLine[2],
		Text: splitLine[9],
	}, nil
}

func ParseAss(filename string) ([]SubLine, error) {
	file, err := os.Open(filename)
	if err != nil {
		return []SubLine{}, err
	}
	defer file.Close()

	subtitles := []SubLine{}
	lineNumber := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() && scanner.Text() != "[Events]" {} // Skip metadata
	scanner.Scan() // Skip formatting line

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "Dialogue") {
			lineNumber++
			event, err := ParseAssEvent(scanner.Text())
			if err != nil { continue }
			event.Index = strconv.Itoa(lineNumber)
			subtitles = append(subtitles, event)
		}
	}

	return subtitles, nil
}
