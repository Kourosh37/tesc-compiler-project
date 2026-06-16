package codegen

import "fmt"

type registers struct {
	next int
	vars map[string]string
}

func newRegisters(paramNames []string) *registers {
	r := &registers{next: 1, vars: map[string]string{}}
	for _, p := range paramNames {
		r.vars[p] = fmt.Sprintf("r%d", r.next)
		r.next++
	}
	return r
}
func (r *registers) alloc() string { s := fmt.Sprintf("r%d", r.next); r.next++; return s }
func (r *registers) varReg(name string) string {
	if v, ok := r.vars[name]; ok {
		return v
	}
	v := r.alloc()
	r.vars[name] = v
	return v
}
