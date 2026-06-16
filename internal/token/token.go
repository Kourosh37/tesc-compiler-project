package token

import "fmt"

type Position struct {
	Line   int
	Column int
}

type TokenType string

const (
	EOF     TokenType = "EOF"
	ILLEGAL TokenType = "ILLEGAL"
	ID      TokenType = "ID"
	NUMBER  TokenType = "NUMBER"
	STRING  TokenType = "STRING"
	MSTRING TokenType = "MSTRING"

	FUNK     TokenType = "FUNK"
	AS       TokenType = "AS"
	IF       TokenType = "IF"
	ELSE     TokenType = "ELSE"
	BEGIN    TokenType = "BEGIN"
	ENDIF    TokenType = "ENDIF"
	WHILE    TokenType = "WHILE"
	ENDWHILE TokenType = "ENDWHILE"
	DO       TokenType = "DO"
	FOR      TokenType = "FOR"
	TO       TokenType = "TO"
	ENDFOR   TokenType = "ENDFOR"
	RETURN   TokenType = "RETURN"
	INT      TokenType = "INT"
	VECTOR   TokenType = "VECTOR"
	STR      TokenType = "STR"
	MSTR     TokenType = "MSTR"
	BOOL     TokenType = "BOOL"
	NULL     TokenType = "NULL"
	TRUE     TokenType = "TRUE"
	FALSE    TokenType = "FALSE"

	LPAREN          TokenType = "LPAREN"
	RPAREN          TokenType = "RPAREN"
	LBRACE          TokenType = "LBRACE"
	RBRACE          TokenType = "RBRACE"
	LBRACKET        TokenType = "LBRACKET"
	RBRACKET        TokenType = "RBRACKET"
	DOUBLE_LBRACKET TokenType = "DOUBLE_LBRACKET"
	DOUBLE_RBRACKET TokenType = "DOUBLE_RBRACKET"
	SEMICOLON       TokenType = "SEMICOLON"
	COMMA           TokenType = "COMMA"
	ASSIGN          TokenType = "ASSIGN"
	DOUBLE_COLON    TokenType = "DOUBLE_COLON"
	ARROW           TokenType = "ARROW"
	PLUS            TokenType = "PLUS"
	MINUS           TokenType = "MINUS"
	STAR            TokenType = "STAR"
	SLASH           TokenType = "SLASH"
	PERCENT         TokenType = "PERCENT"
	LT              TokenType = "LT"
	GT              TokenType = "GT"
	LTE             TokenType = "LTE"
	GTE             TokenType = "GTE"
	EQEQ            TokenType = "EQEQ"
	NOTEQ           TokenType = "NOTEQ"
	AND             TokenType = "AND"
	OR              TokenType = "OR"
	NOT             TokenType = "NOT"
	QUESTION        TokenType = "QUESTION"
	COLON           TokenType = "COLON"
)

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
	Column int
}

func (t Token) Pos() Position { return Position{Line: t.Line, Column: t.Column} }

func (t Token) String() string {
	return fmt.Sprintf("%d:%d %s %q", t.Line, t.Column, t.Type, t.Lexeme)
}

var Keywords = map[string]TokenType{
	"funk": FUNK, "as": AS, "if": IF, "else": ELSE, "begin": BEGIN, "endif": ENDIF,
	"while": WHILE, "endwhile": ENDWHILE, "do": DO, "for": FOR, "to": TO, "endfor": ENDFOR,
	"return": RETURN, "int": INT, "vector": VECTOR, "str": STR, "mstr": MSTR,
	"bool": BOOL, "boolean": BOOL, "null": NULL, "true": TRUE, "false": FALSE,
}
