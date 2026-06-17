package tsvm

type Program struct {
	Procedures map[string]*Procedure
}

type Procedure struct {
	Name         string
	Instructions []Instruction
	Labels       map[string]int
}

type Instruction struct {
	Op       string
	Operands []string
	Line     int
}
