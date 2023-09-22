package parser

type SubLine struct {
	Index string
	StartTime string
	EndTime   string
	Text      string
}

func (line SubLine) Format() string {
	return line.Index + " - " + line.StartTime + " --> " + line.EndTime + ": " + line.Text
	// return [4]string{line.Index, line.StartTime, line.EndTime, line.Text}
}
