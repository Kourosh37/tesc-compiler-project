package semantic

import "teslang-compiler/internal/ast"

type VariableSymbol struct {
	Name         string
	Type         ast.Type
	Initialized  bool
	Line, Column int
}
type ParamSymbol struct {
	Name string
	Type ast.Type
}
type FunctionSymbol struct {
	Name         string
	ReturnType   ast.Type
	Params       []ParamSymbol
	Scope        *Scope
	Line, Column int
	Builtin      bool
}
