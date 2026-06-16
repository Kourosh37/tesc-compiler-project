package diagnostic

import "fmt"

type Diagnostic struct {
	Line     int
	Column   int
	Stage    string
	Function string
	Message  string
}

func (d Diagnostic) Error() string {
	loc := fmt.Sprintf("Error [%s] line %d, column %d", d.Stage, d.Line, d.Column)
	if d.Function != "" {
		loc += fmt.Sprintf(" in function '%s'", d.Function)
	}
	return loc + ": " + d.Message
}

func FormatAll(ds []Diagnostic) string {
	out := ""
	for i, d := range ds {
		if i > 0 {
			out += "\n"
		}
		out += d.Error()
	}
	return out
}
