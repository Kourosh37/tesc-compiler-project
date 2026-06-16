package ast

import "teslang-compiler/internal/token"

type Node interface{ Pos() token.Position }
type Expr interface {
	Node
	exprNode()
}
type Stmt interface {
	Node
	stmtNode()
}
type Decl interface {
	Node
	declNode()
}

type Base struct{ Position token.Position }

func (b Base) Pos() token.Position { return b.Position }

type Type string

const (
	TypeInt     Type = "int"
	TypeVector  Type = "vector"
	TypeStr     Type = "str"
	TypeMStr    Type = "mstr"
	TypeBool    Type = "bool"
	TypeNull    Type = "null"
	TypeUnknown Type = "unknown"
)

type Program struct{ Functions []*FunctionDecl }

func (p *Program) Pos() token.Position {
	if len(p.Functions) > 0 {
		return p.Functions[0].Pos()
	}
	return token.Position{Line: 1, Column: 1}
}

type Param struct {
	Name     string
	Type     Type
	Position token.Position
}

type FunctionDecl struct {
	Base
	Name       string
	ReturnType Type
	Params     []Param
	Body       *BlockStmt
}

func (*FunctionDecl) declNode() {}
func (*FunctionDecl) stmtNode() {}

type BlockStmt struct {
	Base
	Statements []Stmt
}

func (*BlockStmt) stmtNode() {}

type VarDeclStmt struct {
	Base
	Name string
	Type Type
	Init Expr
}

func (*VarDeclStmt) stmtNode() {}

type ReturnStmt struct {
	Base
	Value Expr
}

func (*ReturnStmt) stmtNode() {}

type IfStmt struct {
	Base
	Cond       Expr
	Then, Else []Stmt
}

func (*IfStmt) stmtNode() {}

type WhileStmt struct {
	Base
	Cond Expr
	Body []Stmt
}

func (*WhileStmt) stmtNode() {}

type DoWhileStmt struct {
	Base
	Body []Stmt
	Cond Expr
}

func (*DoWhileStmt) stmtNode() {}

type ForStmt struct {
	Base
	Var        string
	Start, End Expr
	Body       []Stmt
}

func (*ForStmt) stmtNode() {}

type ExprStmt struct {
	Base
	Expr Expr
}

func (*ExprStmt) stmtNode() {}

type IdentifierExpr struct {
	Base
	Name string
}

func (*IdentifierExpr) exprNode() {}

type NumberLiteral struct {
	Base
	Value string
}

func (*NumberLiteral) exprNode() {}

type StringLiteral struct {
	Base
	Value string
}

func (*StringLiteral) exprNode() {}

type MultiStringLiteral struct {
	Base
	Value string
}

func (*MultiStringLiteral) exprNode() {}

type BoolLiteral struct {
	Base
	Value bool
}

func (*BoolLiteral) exprNode() {}

type VectorLiteral struct {
	Base
	Elements []Expr
}

func (*VectorLiteral) exprNode() {}

type AssignExpr struct {
	Base
	Target Expr
	Value  Expr
}

func (*AssignExpr) exprNode() {}

type BinaryExpr struct {
	Base
	Op          token.TokenType
	Left, Right Expr
}

func (*BinaryExpr) exprNode() {}

type UnaryExpr struct {
	Base
	Op   token.TokenType
	Expr Expr
}

func (*UnaryExpr) exprNode() {}

type TernaryExpr struct {
	Base
	Cond, Then, Else Expr
}

func (*TernaryExpr) exprNode() {}

type CallExpr struct {
	Base
	Callee Expr
	Args   []Expr
}

func (*CallExpr) exprNode() {}

type IndexExpr struct {
	Base
	Target, Index Expr
}

func (*IndexExpr) exprNode() {}
