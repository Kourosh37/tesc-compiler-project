package codegen

import (
	"fmt"
	"strings"

	"teslang-compiler/internal/ast"
	"teslang-compiler/internal/diagnostic"
	"teslang-compiler/internal/token"
)

type Generator struct {
	out         []string
	labels      labels
	diagnostics []diagnostic.Diagnostic
	prefix      string
}

func New() *Generator { return &Generator{} }
func (g *Generator) Generate(p *ast.Program) (string, []diagnostic.Diagnostic) {
	for _, f := range p.Functions {
		g.function("", f)
	}
	return strings.Join(g.out, "\n"), g.diagnostics
}
func (g *Generator) emit(s string, args ...any) { g.out = append(g.out, "  "+fmt.Sprintf(s, args...)) }
func (g *Generator) raw(s string, args ...any)  { g.out = append(g.out, fmt.Sprintf(s, args...)) }
func (g *Generator) function(prefix string, f *ast.FunctionDecl) {
	name := f.Name
	if prefix != "" {
		name = prefix + "__" + f.Name
	}
	g.raw("proc %s", name)
	ps := make([]string, len(f.Params))
	for i, p := range f.Params {
		ps[i] = p.Name
	}
	r := newRegisters(ps)
	returned := false
	for _, st := range f.Body.Statements {
		if _, ok := st.(*ast.FunctionDecl); !ok {
			if g.stmt(r, name, st) {
				returned = true
				break
			}
		}
	}
	if !returned {
		g.emit("ret")
	}
	g.raw("")
	for _, st := range f.Body.Statements {
		if nf, ok := st.(*ast.FunctionDecl); ok {
			g.function(name, nf)
		}
	}
}
func (g *Generator) stmt(r *registers, fname string, st ast.Stmt) bool {
	switch n := st.(type) {
	case *ast.VarDeclStmt:
		dst := r.varReg(n.Name)
		if n.Init != nil {
			g.exprInto(r, n.Init, dst)
		}
	case *ast.ExprStmt:
		g.expr(r, n.Expr)
	case *ast.ReturnStmt:
		g.emit("mov r0, %s", g.expr(r, n.Value))
		g.emit("ret")
		return true
	case *ast.IfStmt:
		elseL, endL := g.labels.new("else"), g.labels.new("endif")
		c := g.expr(r, n.Cond)
		g.emit("jz %s, %s", c, elseL)
		for _, s := range n.Then {
			g.stmt(r, fname, s)
		}
		g.emit("jmp %s", endL)
		g.raw("label %s", elseL)
		for _, s := range n.Else {
			g.stmt(r, fname, s)
		}
		g.raw("label %s", endL)
	case *ast.WhileStmt:
		start, end := g.labels.new("while"), g.labels.new("endwhile")
		g.raw("label %s", start)
		c := g.expr(r, n.Cond)
		g.emit("jz %s, %s", c, end)
		for _, s := range n.Body {
			g.stmt(r, fname, s)
		}
		g.emit("jmp %s", start)
		g.raw("label %s", end)
	case *ast.DoWhileStmt:
		start := g.labels.new("do")
		g.raw("label %s", start)
		for _, s := range n.Body {
			g.stmt(r, fname, s)
		}
		c := g.expr(r, n.Cond)
		g.emit("jnz %s, %s", c, start)
	case *ast.ForStmt:
		vr := r.varReg(n.Var)
		g.emit("mov %s, %s", vr, g.expr(r, n.Start))
		endv := g.expr(r, n.End)
		start, end := g.labels.new("for"), g.labels.new("endfor")
		g.raw("label %s", start)
		cond := r.alloc()
		g.emit("lt %s, %s, %s", cond, vr, endv)
		g.emit("jz %s, %s", cond, end)
		for _, s := range n.Body {
			g.stmt(r, fname, s)
		}
		one := r.alloc()
		g.emit("mov %s, 1", one)
		g.emit("add %s, %s, %s", vr, vr, one)
		g.emit("jmp %s", start)
		g.raw("label %s", end)
	}
	return false
}

func (g *Generator) exprInto(r *registers, ex ast.Expr, dst string) string {
	switch n := ex.(type) {
	case *ast.BinaryExpr:
		l, rr := g.expr(r, n.Left), g.expr(r, n.Right)
		op := map[token.TokenType]string{token.PLUS: "add", token.MINUS: "sub", token.STAR: "mul", token.SLASH: "div", token.PERCENT: "mod", token.LT: "lt", token.GT: "gt", token.LTE: "le", token.GTE: "ge", token.EQEQ: "eq", token.NOTEQ: "ne", token.AND: "and", token.OR: "or"}[n.Op]
		g.emit("%s %s, %s, %s", op, dst, l, rr)
	case *ast.CallExpr:
		id := n.Callee.(*ast.IdentifierExpr)
		args := []string{}
		for _, a := range n.Args {
			args = append(args, g.expr(r, a))
		}
		switch id.Name {
		case "scan":
			g.emit("call read, %s", dst)
		case "list":
			g.emit("call list, %s%s", dst, argSuffix(args))
		case "length":
			g.emit("call length, %s%s", dst, argSuffix(args))
		default:
			g.emit("call %s, %s%s", id.Name, dst, argSuffix(args))
		}
	default:
		g.emit("mov %s, %s", dst, g.expr(r, ex))
	}
	return dst
}
func (g *Generator) expr(r *registers, ex ast.Expr) string {
	switch n := ex.(type) {
	case *ast.IdentifierExpr:
		return r.varReg(n.Name)
	case *ast.NumberLiteral:
		return n.Value
	case *ast.StringLiteral:
		return fmt.Sprintf("%q", n.Value)
	case *ast.MultiStringLiteral:
		return fmt.Sprintf("%q", n.Value)
	case *ast.BoolLiteral:
		if n.Value {
			return "1"
		}
		return "0"
	case *ast.VectorLiteral:
		dst := r.alloc()
		args := []string{}
		for _, e := range n.Elements {
			args = append(args, g.expr(r, e))
		}
		g.emit("call vector, %s, %s", dst, strings.Join(args, ", "))
		return dst
	case *ast.UnaryExpr:
		v := g.expr(r, n.Expr)
		dst := r.alloc()
		if n.Op == token.NOT {
			g.emit("not %s, %s", dst, v)
		} else if n.Op == token.MINUS {
			g.emit("sub %s, 0, %s", dst, v)
		} else {
			g.emit("mov %s, %s", dst, v)
		}
		return dst
	case *ast.BinaryExpr:
		l, rr, dst := g.expr(r, n.Left), g.expr(r, n.Right), r.alloc()
		op := map[token.TokenType]string{token.PLUS: "add", token.MINUS: "sub", token.STAR: "mul", token.SLASH: "div", token.PERCENT: "mod", token.LT: "lt", token.GT: "gt", token.LTE: "le", token.GTE: "ge", token.EQEQ: "eq", token.NOTEQ: "ne", token.AND: "and", token.OR: "or"}[n.Op]
		g.emit("%s %s, %s, %s", op, dst, l, rr)
		return dst
	case *ast.TernaryExpr:
		dst := r.alloc()
		elseL, endL := g.labels.new("tern_else"), g.labels.new("tern_end")
		c := g.expr(r, n.Cond)
		g.emit("jz %s, %s", c, elseL)
		g.emit("mov %s, %s", dst, g.expr(r, n.Then))
		g.emit("jmp %s", endL)
		g.raw("label %s", elseL)
		g.emit("mov %s, %s", dst, g.expr(r, n.Else))
		g.raw("label %s", endL)
		return dst
	case *ast.AssignExpr:
		switch t := n.Target.(type) {
		case *ast.IdentifierExpr:
			dst := r.varReg(t.Name)
			return g.exprInto(r, n.Value, dst)
		case *ast.IndexExpr:
			val := g.expr(r, n.Value)
			arr := g.expr(r, t.Target)
			idx := g.expr(r, t.Index)
			g.emit("storeidx %s, %s, %s", arr, idx, val)
			return val
		}
	case *ast.CallExpr:
		id := n.Callee.(*ast.IdentifierExpr)
		args := []string{}
		for _, a := range n.Args {
			args = append(args, g.expr(r, a))
		}
		if id.Name == "scan" {
			dst := r.alloc()
			g.emit("call read, %s", dst)
			return dst
		}
		if id.Name == "print" {
			if len(args) > 0 {
				g.emit("call log, %s", args[0])
			}
			return "0"
		}
		if id.Name == "exit" {
			if len(args) > 0 {
				g.emit("call exit, %s", args[0])
			}
			return "0"
		}
		dst := r.alloc()
		g.emit("call %s, %s%s", id.Name, dst, argSuffix(args))
		return dst
	case *ast.IndexExpr:
		dst := r.alloc()
		g.emit("loadidx %s, %s, %s", dst, g.expr(r, n.Target), g.expr(r, n.Index))
		return dst
	}
	return "0"
}
func argSuffix(args []string) string {
	if len(args) == 0 {
		return ""
	}
	return ", " + strings.Join(args, ", ")
}
