package parser

import (
	"bufio"
	"os"
	"strings"
)

func ParseSrt(filename string) ([]SubLine, error) {
	file, err := os.Open(filename)
	if err != nil {
		return []SubLine{}, err
	}
	defer file.Close()

	subtitles := []SubLine{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if scanner.Text() == "" {
			continue
		}
		event := SubLine{
			Index: scanner.Text(),
		}

		scanner.Scan()
		times := strings.Split(scanner.Text(), " --> ")
		if len(times) != 2 { break }
		event.StartTime = times[0]
		event.EndTime = times[1]

		for scanner.Scan() && scanner.Text() != "" {
			event.Text += scanner.Text() + " "
		}
		event.Text = event.Text[:len(event.Text) - 1] // Removing extra space

		subtitles = append(subtitles, event)
	}

	return subtitles, nil
}
