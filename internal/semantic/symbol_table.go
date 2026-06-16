package semantic

import "fmt"

type Scope struct {
	Parent    *Scope
	Variables map[string]*VariableSymbol
	Functions map[string]*FunctionSymbol
}

func NewScope(parent *Scope) *Scope {
	return &Scope{Parent: parent, Variables: map[string]*VariableSymbol{}, Functions: map[string]*FunctionSymbol{}}
}
func (s *Scope) DefineVariable(v *VariableSymbol) error {
	if _, ok := s.Variables[v.Name]; ok {
		return fmt.Errorf("duplicate variable declaration '%s'", v.Name)
	}
	s.Variables[v.Name] = v
	return nil
}
func (s *Scope) LookupVariable(name string) (*VariableSymbol, bool) {
	for sc := s; sc != nil; sc = sc.Parent {
		if v, ok := sc.Variables[name]; ok {
			return v, true
		}
	}
	return nil, false
}
func (s *Scope) DefineFunction(f *FunctionSymbol) error {
	if _, ok := s.Functions[f.Name]; ok {
		return fmt.Errorf("duplicate function declaration '%s'", f.Name)
	}
	s.Functions[f.Name] = f
	return nil
}
func (s *Scope) LookupFunction(name string) (*FunctionSymbol, bool) {
	for sc := s; sc != nil; sc = sc.Parent {
		if f, ok := sc.Functions[name]; ok {
			return f, true
		}
	}
	return nil, false
}
