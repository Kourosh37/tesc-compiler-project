package tsvm

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

func Parse(r io.Reader) (*Program, error) {
	prog := &Program{Procedures: map[string]*Procedure{}}
	var current *Procedure
	scanner := bufio.NewScanner(r)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := stripComment(scanner.Text())
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		switch fields[0] {
		case "proc":
			if len(fields) != 2 {
				return nil, fmt.Errorf("line %d: proc requires a name", lineNo)
			}
			if _, exists := prog.Procedures[fields[1]]; exists {
				return nil, fmt.Errorf("line %d: duplicate procedure %q", lineNo, fields[1])
			}
			current = &Procedure{Name: fields[1], Labels: map[string]int{}}
			prog.Procedures[current.Name] = current
		case "label":
			if current == nil {
				return nil, fmt.Errorf("line %d: label outside procedure", lineNo)
			}
			if len(fields) != 2 {
				return nil, fmt.Errorf("line %d: label requires a name", lineNo)
			}
			current.Labels[fields[1]] = len(current.Instructions)
		default:
			if current == nil {
				return nil, fmt.Errorf("line %d: instruction outside procedure", lineNo)
			}
			op, rest, _ := strings.Cut(line, " ")
			ops, err := splitOperands(rest)
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", lineNo, err)
			}
			current.Instructions = append(current.Instructions, Instruction{Op: op, Operands: ops, Line: lineNo})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return prog, nil
}

func stripComment(line string) string {
	inQuote := false
	escaped := false
	for i, r := range line {
		if escaped {
			escaped = false
			continue
		}
		if r == '\\' && inQuote {
			escaped = true
			continue
		}
		if r == '"' {
			inQuote = !inQuote
			continue
		}
		if r == '#' && !inQuote {
			return line[:i]
		}
	}
	return line
}

func splitOperands(s string) ([]string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	var out []string
	var b strings.Builder
	inQuote := false
	escaped := false
	for _, r := range s {
		if escaped {
			b.WriteRune(r)
			escaped = false
			continue
		}
		if r == '\\' && inQuote {
			b.WriteRune(r)
			escaped = true
			continue
		}
		if r == '"' {
			inQuote = !inQuote
			b.WriteRune(r)
			continue
		}
		if r == ',' && !inQuote {
			out = append(out, strings.TrimSpace(b.String()))
			b.Reset()
			continue
		}
		b.WriteRune(r)
	}
	if inQuote {
		return nil, fmt.Errorf("unterminated quoted operand")
	}
	if strings.TrimSpace(b.String()) != "" {
		out = append(out, strings.TrimSpace(b.String()))
	}
	for _, op := range out {
		if op == "" || strings.IndexFunc(op, unicode.IsSpace) == 0 {
			return nil, fmt.Errorf("empty operand")
		}
	}
	return out, nil
}
