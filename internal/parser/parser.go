package parser

import (
	"fmt"

	"teslang-compiler/internal/ast"
	"teslang-compiler/internal/diagnostic"
	"teslang-compiler/internal/token"
)

type Parser struct {
	toks        []token.Token
	i           int
	diagnostics []diagnostic.Diagnostic
}

func New(toks []token.Token) *Parser                   { return &Parser{toks: toks} }
func (p *Parser) Diagnostics() []diagnostic.Diagnostic { return p.diagnostics }

func (p *Parser) ParseProgram() (*ast.Program, []diagnostic.Diagnostic) {
	prog := &ast.Program{}
	for !p.at(token.EOF) {
		if p.at(token.FUNK) {
			prog.Functions = append(prog.Functions, p.function())
		} else {
			p.err(p.peek(), "expected function declaration")
			p.advance()
		}
	}
	return prog, p.diagnostics
}

func (p *Parser) function() *ast.FunctionDecl {
	start := p.expect(token.FUNK, "expected 'funk'")
	p.expect(token.LT, "expected '<'")
	rt := p.parseType()
	p.expect(token.GT, "expected '>'")
	name := p.expect(token.ID, "expected function name")
	p.expect(token.LPAREN, "expected '('")
	var params []ast.Param
	if !p.at(token.RPAREN) {
		for {
			n := p.expect(token.ID, "expected parameter name")
			p.expect(token.AS, "expected 'as'")
			params = append(params, ast.Param{Name: n.Lexeme, Type: p.parseType(), Position: n.Pos()})
			if !p.match(token.COMMA) {
				break
			}
		}
	}
	p.expect(token.RPAREN, "expected ')'")
	var body *ast.BlockStmt
	if p.match(token.ARROW) {
		p.expect(token.RETURN, "expected return")
		ex := p.expression(0)
		p.expect(token.SEMICOLON, "expected ';'")
		body = &ast.BlockStmt{Base: ast.Base{Position: start.Pos()}, Statements: []ast.Stmt{&ast.ReturnStmt{Base: ast.Base{Position: start.Pos()}, Value: ex}}}
	} else {
		body = p.block()
	}
	return &ast.FunctionDecl{Base: ast.Base{Position: start.Pos()}, Name: name.Lexeme, ReturnType: rt, Params: params, Body: body}
}

func (p *Parser) block() *ast.BlockStmt {
	t := p.expect(token.LBRACE, "expected '{'")
	var ss []ast.Stmt
	for !p.at(token.RBRACE) && !p.at(token.EOF) {
		ss = append(ss, p.statement())
	}
	p.expect(token.RBRACE, "expected '}'")
	return &ast.BlockStmt{Base: ast.Base{Position: t.Pos()}, Statements: ss}
}

func (p *Parser) statement() ast.Stmt {
	switch {
	case p.at(token.FUNK):
		return p.function()
	case p.at(token.RETURN):
		return p.returnStmt()
	case p.at(token.IF):
		return p.ifStmt()
	case p.at(token.WHILE):
		return p.whileStmt()
	case p.at(token.DO):
		return p.doWhileStmt()
	case p.at(token.FOR):
		return p.forStmt()
	case p.at(token.ID) && p.peekN(1).Type == token.DOUBLE_COLON:
		return p.varDecl()
	default:
		ex := p.expression(0)
		p.expect(token.SEMICOLON, "expected ';'")
		return &ast.ExprStmt{Base: ast.Base{Position: ex.Pos()}, Expr: ex}
	}
}

func (p *Parser) varDecl() ast.Stmt {
	n := p.expect(token.ID, "expected variable name")
	p.expect(token.DOUBLE_COLON, "expected '::'")
	typ := p.parseType()
	var init ast.Expr
	if p.match(token.ASSIGN) {
		init = p.expression(0)
	}
	p.expect(token.SEMICOLON, "expected ';'")
	return &ast.VarDeclStmt{Base: ast.Base{Position: n.Pos()}, Name: n.Lexeme, Type: typ, Init: init}
}
func (p *Parser) returnStmt() ast.Stmt {
	t := p.advance()
	ex := p.expression(0)
	p.expect(token.SEMICOLON, "expected ';'")
	return &ast.ReturnStmt{Base: ast.Base{Position: t.Pos()}, Value: ex}
}
func (p *Parser) ifStmt() ast.Stmt {
	t := p.advance()
	p.expect(token.DOUBLE_LBRACKET, "expected '[['")
	c := p.expression(0)
	p.expect(token.DOUBLE_RBRACKET, "expected ']]'")
	p.expect(token.BEGIN, "expected begin")
	then := p.until(token.ELSE, token.ENDIF)
	var els []ast.Stmt
	if p.match(token.ELSE) {
		els = p.until(token.ENDIF)
	}
	p.expect(token.ENDIF, "expected endif")
	return &ast.IfStmt{Base: ast.Base{Position: t.Pos()}, Cond: c, Then: then, Else: els}
}
func (p *Parser) whileStmt() ast.Stmt {
	t := p.advance()
	p.expect(token.DOUBLE_LBRACKET, "expected '[['")
	c := p.expression(0)
	p.expect(token.DOUBLE_RBRACKET, "expected ']]'")
	p.expect(token.BEGIN, "expected begin")
	b := p.until(token.ENDWHILE)
	p.expect(token.ENDWHILE, "expected endwhile")
	return &ast.WhileStmt{Base: ast.Base{Position: t.Pos()}, Cond: c, Body: b}
}
func (p *Parser) doWhileStmt() ast.Stmt {
	t := p.advance()
	p.expect(token.BEGIN, "expected begin")
	b := p.until(token.WHILE)
	p.expect(token.WHILE, "expected while")
	p.expect(token.DOUBLE_LBRACKET, "expected '[['")
	c := p.expression(0)
	p.expect(token.DOUBLE_RBRACKET, "expected ']]'")
	p.expect(token.ENDWHILE, "expected endwhile")
	return &ast.DoWhileStmt{Base: ast.Base{Position: t.Pos()}, Body: b, Cond: c}
}
func (p *Parser) forStmt() ast.Stmt {
	t := p.advance()
	p.expect(token.LPAREN, "expected '('")
	n := p.expect(token.ID, "expected loop variable")
	p.expect(token.ASSIGN, "expected '='")
	s := p.expression(0)
	p.expect(token.TO, "expected to")
	e := p.expression(0)
	p.expect(token.RPAREN, "expected ')'")
	p.expect(token.BEGIN, "expected begin")
	b := p.until(token.ENDFOR)
	p.expect(token.ENDFOR, "expected endfor")
	return &ast.ForStmt{Base: ast.Base{Position: t.Pos()}, Var: n.Lexeme, Start: s, End: e, Body: b}
}
func (p *Parser) until(ends ...token.TokenType) []ast.Stmt {
	var ss []ast.Stmt
	for !p.at(token.EOF) && !p.atAny(ends...) {
		ss = append(ss, p.statement())
	}
	return ss
}

func (p *Parser) parseType() ast.Type {
	t := p.advance()
	switch t.Type {
	case token.INT:
		return ast.TypeInt
	case token.VECTOR:
		return ast.TypeVector
	case token.STR:
		return ast.TypeStr
	case token.MSTR:
		return ast.TypeMStr
	case token.BOOL:
		return ast.TypeBool
	case token.NULL:
		return ast.TypeNull
	default:
		p.err(t, fmt.Sprintf("wrong type '%s' found. types must be one of: int, vector, str, mstr, bool, null.", t.Lexeme))
		return ast.TypeUnknown
	}
}

const (
	precAssign  = 1
	precTernary = 2
	precOr      = 3
	precAnd     = 4
	precEq      = 5
	precCmp     = 6
	precAdd     = 7
	precMul     = 8
	precUnary   = 9
	precPostfix = 10
)

func (p *Parser) expression(min int) ast.Expr {
	left := p.prefix()
	for {
		t := p.peek()
		if t.Type == token.LPAREN && precPostfix >= min {
			left = p.call(left)
			continue
		}
		if t.Type == token.LBRACKET && precPostfix >= min {
			left = p.index(left)
			continue
		}
		if t.Type == token.QUESTION && precTernary >= min {
			p.advance()
			a := p.expression(0)
			p.expect(token.COLON, "expected ':'")
			b := p.expression(precTernary)
			left = &ast.TernaryExpr{Base: ast.Base{Position: left.Pos()}, Cond: left, Then: a, Else: b}
			continue
		}
		if t.Type == token.ASSIGN && precAssign >= min {
			p.advance()
			right := p.expression(precAssign)
			left = &ast.AssignExpr{Base: ast.Base{Position: left.Pos()}, Target: left, Value: right}
			continue
		}
		prec := infixPrec(t.Type)
		if prec < min || prec == 0 {
			break
		}
		op := p.advance()
		right := p.expression(prec + 1)
		left = &ast.BinaryExpr{Base: ast.Base{Position: left.Pos()}, Op: op.Type, Left: left, Right: right}
	}
	return left
}
func (p *Parser) prefix() ast.Expr {
	t := p.advance()
	switch t.Type {
	case token.ID:
		return &ast.IdentifierExpr{Base: ast.Base{Position: t.Pos()}, Name: t.Lexeme}
	case token.NUMBER:
		return &ast.NumberLiteral{Base: ast.Base{Position: t.Pos()}, Value: t.Lexeme}
	case token.STRING:
		return &ast.StringLiteral{Base: ast.Base{Position: t.Pos()}, Value: t.Lexeme}
	case token.MSTRING:
		return &ast.MultiStringLiteral{Base: ast.Base{Position: t.Pos()}, Value: t.Lexeme}
	case token.TRUE, token.FALSE:
		return &ast.BoolLiteral{Base: ast.Base{Position: t.Pos()}, Value: t.Type == token.TRUE}
	case token.NOT, token.PLUS, token.MINUS:
		return &ast.UnaryExpr{Base: ast.Base{Position: t.Pos()}, Op: t.Type, Expr: p.expression(precUnary)}
	case token.LPAREN:
		ex := p.expression(0)
		p.expect(token.RPAREN, "expected ')'")
		return ex
	case token.LBRACKET:
		return p.vector(t)
	default:
		p.err(t, "expected expression")
		return &ast.NumberLiteral{Base: ast.Base{Position: t.Pos()}, Value: "0"}
	}
}
func (p *Parser) vector(t token.Token) ast.Expr {
	var es []ast.Expr
	if !p.at(token.RBRACKET) {
		for {
			es = append(es, p.expression(0))
			if !p.match(token.COMMA) {
				break
			}
		}
	}
	p.expect(token.RBRACKET, "expected ']'")
	return &ast.VectorLiteral{Base: ast.Base{Position: t.Pos()}, Elements: es}
}
func (p *Parser) call(c ast.Expr) ast.Expr {
	p.advance()
	var args []ast.Expr
	if !p.at(token.RPAREN) {
		for {
			args = append(args, p.expression(0))
			if !p.match(token.COMMA) {
				break
			}
		}
	}
	p.expect(token.RPAREN, "expected ')'")
	return &ast.CallExpr{Base: ast.Base{Position: c.Pos()}, Callee: c, Args: args}
}
func (p *Parser) index(c ast.Expr) ast.Expr {
	p.advance()
	idx := p.expression(0)
	p.expect(token.RBRACKET, "expected ']'")
	return &ast.IndexExpr{Base: ast.Base{Position: c.Pos()}, Target: c, Index: idx}
}
func infixPrec(t token.TokenType) int {
	switch t {
	case token.OR:
		return precOr
	case token.AND:
		return precAnd
	case token.EQEQ, token.NOTEQ:
		return precEq
	case token.LT, token.GT, token.LTE, token.GTE:
		return precCmp
	case token.PLUS, token.MINUS:
		return precAdd
	case token.STAR, token.SLASH, token.PERCENT:
		return precMul
	}
	return 0
}
func (p *Parser) at(t token.TokenType) bool { return p.peek().Type == t }
func (p *Parser) atAny(ts ...token.TokenType) bool {
	for _, t := range ts {
		if p.at(t) {
			return true
		}
	}
	return false
}
func (p *Parser) match(t token.TokenType) bool {
	if p.at(t) {
		p.advance()
		return true
	}
	return false
}
func (p *Parser) expect(t token.TokenType, msg string) token.Token {
	if p.at(t) {
		return p.advance()
	}
	tok := p.peek()
	p.err(tok, msg)
	return tok
}
func (p *Parser) advance() token.Token {
	t := p.peek()
	if !p.at(token.EOF) {
		p.i++
	}
	return t
}
func (p *Parser) peek() token.Token { return p.peekN(0) }
func (p *Parser) peekN(n int) token.Token {
	if p.i+n >= len(p.toks) {
		return token.Token{Type: token.EOF, Line: 1, Column: 1}
	}
	return p.toks[p.i+n]
}
func (p *Parser) err(t token.Token, msg string) {
	p.diagnostics = append(p.diagnostics, diagnostic.Diagnostic{Line: t.Line, Column: t.Column, Stage: "parser", Message: msg})
}
