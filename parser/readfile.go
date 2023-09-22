package parser

import (
	"errors"
	"strings"
)

func ParseSub(filename string) ([]SubLine, error) {
	switch strings.ToLower(filename[strings.LastIndex(filename, "."):]) {
		case ".ass", ".ssa":
			return ParseAss(filename)
		case ".srt":
			return ParseSrt(filename)
		default:
			return nil, errors.New("Unrecognized filetype")
	}
}
