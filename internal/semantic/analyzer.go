package semantic

import (
	"fmt"

	"teslang-compiler/internal/ast"
	"teslang-compiler/internal/diagnostic"
	"teslang-compiler/internal/token"
)

type Analyzer struct {
	diagnostics []diagnostic.Diagnostic
	function    string
	returnType  ast.Type
}

func New() *Analyzer { return &Analyzer{} }
func (a *Analyzer) Analyze(p *ast.Program) []diagnostic.Diagnostic {
	global := NewScope(nil)
	a.builtins(global)
	for _, f := range p.Functions {
		a.declareFunction(global, f)
	}
	for _, f := range p.Functions {
		a.analyzeFunction(global, f)
	}
	return a.diagnostics
}
func (a *Analyzer) builtins(s *Scope) {
	add := func(n string, rt ast.Type, ps ...ParamSymbol) {
		_ = s.DefineFunction(&FunctionSymbol{Name: n, ReturnType: rt, Params: ps, Builtin: true})
	}
	add("scan", ast.TypeInt)
	add("print", ast.TypeNull, ParamSymbol{Name: "x", Type: ast.TypeUnknown})
	add("list", ast.TypeVector, ParamSymbol{Name: "n", Type: ast.TypeInt})
	add("length", ast.TypeInt, ParamSymbol{Name: "arr", Type: ast.TypeVector})
	add("random", ast.TypeInt, ParamSymbol{Name: "min", Type: ast.TypeInt}, ParamSymbol{Name: "max", Type: ast.TypeInt})
	add("exit", ast.TypeNull, ParamSymbol{Name: "n", Type: ast.TypeInt})
}
func (a *Analyzer) declareFunction(s *Scope, f *ast.FunctionDecl) {
	params := make([]ParamSymbol, len(f.Params))
	for i, p := range f.Params {
		params[i] = ParamSymbol{Name: p.Name, Type: p.Type}
	}
	if err := s.DefineFunction(&FunctionSymbol{Name: f.Name, ReturnType: f.ReturnType, Params: params, Line: f.Pos().Line, Column: f.Pos().Column}); err != nil {
		a.errAt(f.Pos(), err.Error())
	}
}
func (a *Analyzer) analyzeFunction(parent *Scope, f *ast.FunctionDecl) {
	oldFn, oldRet := a.function, a.returnType
	a.function, a.returnType = f.Name, f.ReturnType
	scope := NewScope(parent)
	if fs, ok := parent.LookupFunction(f.Name); ok {
		fs.Scope = scope
	}
	for _, p := range f.Params {
		if !validType(p.Type) {
			a.errAt(p.Position, fmt.Sprintf("wrong type '%s' found. types must be one of: int, vector, str, mstr, bool, null.", p.Type))
		}
		_ = scope.DefineVariable(&VariableSymbol{Name: p.Name, Type: p.Type, Initialized: true, Line: p.Position.Line, Column: p.Position.Column})
	}
	for _, st := range f.Body.Statements {
		if nf, ok := st.(*ast.FunctionDecl); ok {
			a.declareFunction(scope, nf)
		}
	}
	for _, st := range f.Body.Statements {
		a.stmt(scope, st)
	}
	a.function, a.returnType = oldFn, oldRet
}
func (a *Analyzer) stmt(s *Scope, st ast.Stmt) {
	switch n := st.(type) {
	case *ast.FunctionDecl:
		a.analyzeFunction(s, n)
	case *ast.VarDeclStmt:
		if !validType(n.Type) {
			a.err(n, fmt.Sprintf("wrong type '%s' found. types must be one of: int, vector, str, mstr, bool, null.", n.Type))
		}
		sym := &VariableSymbol{Name: n.Name, Type: n.Type, Initialized: n.Init != nil, Line: n.Pos().Line, Column: n.Pos().Column}
		if err := s.DefineVariable(sym); err != nil {
			a.err(n, err.Error())
		}
		if n.Init != nil {
			if got := a.expr(s, n.Init); !same(n.Type, got) {
				a.err(n, fmt.Sprintf("variable '%s' expected to be of type '%s' but got '%s'.", n.Name, n.Type, got))
			}
		}
	case *ast.ReturnStmt:
		got := a.expr(s, n.Value)
		if a.returnType == ast.TypeNull {
			if !(a.function == "main" && got == ast.TypeInt && isZero(n.Value)) && got != ast.TypeNull {
				a.err(n, fmt.Sprintf("wrong return type. expected '%s' but got '%s'.", a.returnType, got))
			}
		} else if !same(a.returnType, got) {
			a.err(n, fmt.Sprintf("wrong return type. expected '%s' but got '%s'.", a.returnType, got))
		}
	case *ast.ExprStmt:
		a.expr(s, n.Expr)
	case *ast.IfStmt:
		if t := a.expr(s, n.Cond); t != ast.TypeBool {
			a.err(n, "condition expression must be bool.")
		}
		a.stmts(NewScope(s), n.Then)
		if n.Else != nil {
			a.stmts(NewScope(s), n.Else)
		}
	case *ast.WhileStmt:
		if t := a.expr(s, n.Cond); t != ast.TypeBool {
			a.err(n, "condition expression must be bool.")
		}
		a.stmts(NewScope(s), n.Body)
	case *ast.DoWhileStmt:
		a.stmts(NewScope(s), n.Body)
		if t := a.expr(s, n.Cond); t != ast.TypeBool {
			a.err(n, "condition expression must be bool.")
		}
	case *ast.ForStmt:
		ls := NewScope(s)
		if t := a.expr(s, n.Start); t != ast.TypeInt {
			a.err(n, "for-loop start must be int.")
		}
		if t := a.expr(s, n.End); t != ast.TypeInt {
			a.err(n, "for-loop end must be int.")
		}
		if v, ok := ls.LookupVariable(n.Var); ok {
			if v.Type != ast.TypeInt {
				a.err(n, "for-loop variable must be int.")
			}
			v.Initialized = true
		} else {
			_ = ls.DefineVariable(&VariableSymbol{Name: n.Var, Type: ast.TypeInt, Initialized: true, Line: n.Pos().Line, Column: n.Pos().Column})
		}
		a.stmts(ls, n.Body)
	}
}
func (a *Analyzer) stmts(s *Scope, ss []ast.Stmt) {
	for _, st := range ss {
		a.stmt(s, st)
	}
}
func (a *Analyzer) expr(s *Scope, ex ast.Expr) ast.Type {
	switch n := ex.(type) {
	case *ast.IdentifierExpr:
		v, ok := s.LookupVariable(n.Name)
		if !ok {
			a.err(n, fmt.Sprintf("variable '%s' is not defined.", n.Name))
			return ast.TypeUnknown
		}
		if !v.Initialized {
			a.err(n, fmt.Sprintf("variable '%s' is used before being assigned.", n.Name))
		}
		return v.Type
	case *ast.NumberLiteral:
		return ast.TypeInt
	case *ast.StringLiteral:
		return ast.TypeStr
	case *ast.MultiStringLiteral:
		return ast.TypeMStr
	case *ast.BoolLiteral:
		return ast.TypeBool
	case *ast.VectorLiteral:
		for _, e := range n.Elements {
			if t := a.expr(s, e); t != ast.TypeInt {
				a.err(e, "vector literal elements must be int.")
			}
		}
		return ast.TypeVector
	case *ast.UnaryExpr:
		t := a.expr(s, n.Expr)
		if (n.Op == token.MINUS || n.Op == token.PLUS) && t == ast.TypeInt {
			return ast.TypeInt
		}
		if n.Op == token.NOT && t == ast.TypeBool {
			return ast.TypeBool
		}
		a.err(n, "invalid unary operand type.")
		return ast.TypeUnknown
	case *ast.BinaryExpr:
		return a.binary(s, n)
	case *ast.TernaryExpr:
		if t := a.expr(s, n.Cond); t != ast.TypeBool {
			a.err(n, "ternary condition must be bool.")
		}
		tt := a.expr(s, n.Then)
		et := a.expr(s, n.Else)
		if !same(tt, et) {
			a.err(n, "ternary branches must have compatible types.")
			return ast.TypeUnknown
		}
		return tt
	case *ast.AssignExpr:
		return a.assign(s, n)
	case *ast.CallExpr:
		return a.call(s, n)
	case *ast.IndexExpr:
		if t := a.expr(s, n.Target); t != ast.TypeVector {
			a.err(n, "index target must be vector.")
		}
		if t := a.expr(s, n.Index); t != ast.TypeInt {
			a.err(n, "index must be int.")
		}
		return ast.TypeInt
	}
	return ast.TypeUnknown
}
func (a *Analyzer) assign(s *Scope, n *ast.AssignExpr) ast.Type {
	got := a.expr(s, n.Value)
	switch t := n.Target.(type) {
	case *ast.IdentifierExpr:
		v, ok := s.LookupVariable(t.Name)
		if !ok {
			a.err(t, fmt.Sprintf("variable '%s' is not defined.", t.Name))
			return ast.TypeUnknown
		}
		if !same(v.Type, got) {
			a.err(t, fmt.Sprintf("variable '%s' expected to be of type '%s' but got '%s'.", t.Name, v.Type, got))
		}
		v.Initialized = true
		return v.Type
	case *ast.IndexExpr:
		a.expr(s, t)
		if got != ast.TypeInt {
			a.err(n, "vector assignment value must be int.")
		}
		return ast.TypeInt
	default:
		a.err(n, "invalid assignment target.")
		return ast.TypeUnknown
	}
}
func (a *Analyzer) call(s *Scope, n *ast.CallExpr) ast.Type {
	id, ok := n.Callee.(*ast.IdentifierExpr)
	if !ok {
		a.err(n, "function call target must be identifier.")
		return ast.TypeUnknown
	}
	f, ok := s.LookupFunction(id.Name)
	if !ok {
		a.err(n, fmt.Sprintf("function '%s' is not defined.", id.Name))
		return ast.TypeUnknown
	}
	if len(f.Params) != len(n.Args) {
		a.err(n, fmt.Sprintf("function '%s' expects %d arguments but got %d.", id.Name, len(f.Params), len(n.Args)))
	}
	for i, arg := range n.Args {
		got := a.expr(s, arg)
		if i < len(f.Params) && f.Params[i].Type != ast.TypeUnknown && !same(f.Params[i].Type, got) {
			a.err(arg, fmt.Sprintf("function '%s': expected argument '%s' to be of type '%s' but got '%s'.", id.Name, f.Params[i].Name, f.Params[i].Type, got))
		}
		if id.Name == "print" && got == ast.TypeVector {
			a.err(arg, "function 'print' does not accept vector arguments.")
		}
	}
	return f.ReturnType
}
func (a *Analyzer) binary(s *Scope, n *ast.BinaryExpr) ast.Type {
	l, r := a.expr(s, n.Left), a.expr(s, n.Right)
	switch n.Op {
	case token.PLUS:
		if l == ast.TypeInt && r == ast.TypeInt {
			return ast.TypeInt
		}
		if (l == ast.TypeStr || l == ast.TypeMStr) && (r == ast.TypeStr || r == ast.TypeMStr) {
			if l == ast.TypeMStr || r == ast.TypeMStr {
				return ast.TypeMStr
			}
			return ast.TypeStr
		}
	case token.MINUS, token.STAR, token.SLASH, token.PERCENT:
		if l == ast.TypeInt && r == ast.TypeInt {
			return ast.TypeInt
		}
	case token.LT, token.GT, token.LTE, token.GTE:
		if l == ast.TypeInt && r == ast.TypeInt {
			return ast.TypeBool
		}
	case token.EQEQ, token.NOTEQ:
		if same(l, r) {
			return ast.TypeBool
		}
	case token.AND, token.OR:
		if l == ast.TypeBool && r == ast.TypeBool {
			return ast.TypeBool
		}
	}
	a.err(n, fmt.Sprintf("invalid operand types '%s' and '%s'.", l, r))
	return ast.TypeUnknown
}
func same(a, b ast.Type) bool                  { return a == b || a == ast.TypeUnknown || b == ast.TypeUnknown }
func isZero(e ast.Expr) bool                   { n, ok := e.(*ast.NumberLiteral); return ok && n.Value == "0" }
func (a *Analyzer) err(n ast.Node, msg string) { a.errAt(n.Pos(), msg) }
func (a *Analyzer) errAt(p token.Position, msg string) {
	a.diagnostics = append(a.diagnostics, diagnostic.Diagnostic{Line: p.Line, Column: p.Column, Stage: "semantic", Function: a.function, Message: msg})
}
