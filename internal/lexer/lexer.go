package lexer

import (
	"strings"
	"unicode"

	"teslang-compiler/internal/diagnostic"
	"teslang-compiler/internal/token"
)

type Lexer struct {
	src          []rune
	i, line, col int
	diagnostics  []diagnostic.Diagnostic
}

func New(input string) *Lexer { return &Lexer{src: []rune(input), line: 1, col: 1} }

func (l *Lexer) Diagnostics() []diagnostic.Diagnostic { return l.diagnostics }

func (l *Lexer) LexAll() ([]token.Token, []diagnostic.Diagnostic) {
	var toks []token.Token
	for {
		t := l.Next()
		toks = append(toks, t)
		if t.Type == token.EOF {
			break
		}
	}
	return toks, l.diagnostics
}

func (l *Lexer) Next() token.Token {
	for {
		l.skipWhitespace()
		if l.matchString("</") {
			l.skipComment()
			continue
		}
		break
	}
	startLine, startCol := l.line, l.col
	if l.eof() {
		return token.Token{Type: token.EOF, Line: startLine, Column: startCol}
	}
	ch := l.peek()
	if isIdentStart(ch) {
		lex := l.takeWhile(isIdentPart)
		if typ, ok := token.Keywords[lex]; ok {
			return token.Token{Type: typ, Lexeme: lex, Line: startLine, Column: startCol}
		}
		return token.Token{Type: token.ID, Lexeme: lex, Line: startLine, Column: startCol}
	}
	if unicode.IsDigit(ch) {
		return token.Token{Type: token.NUMBER, Lexeme: l.takeWhile(unicode.IsDigit), Line: startLine, Column: startCol}
	}
	if l.matchString(`"""`) {
		return l.multiString(startLine, startCol)
	}
	if ch == '"' || ch == '\'' {
		return l.string(startLine, startCol, ch)
	}
	if tok, ok := l.operator(startLine, startCol); ok {
		return tok
	}
	l.advance()
	l.error(startLine, startCol, "unexpected character '"+string(ch)+"'")
	return token.Token{Type: token.ILLEGAL, Lexeme: string(ch), Line: startLine, Column: startCol}
}

func (l *Lexer) skipWhitespace() {
	for !l.eof() {
		switch l.peek() {
		case ' ', '\t', '\r', '\n':
			l.advance()
		default:
			return
		}
	}
}

func (l *Lexer) skipComment() {
	line, col := l.line, l.col
	l.advanceN(2)
	depth := 1
	for !l.eof() && depth > 0 {
		if l.matchString("</") {
			depth++
			l.advanceN(2)
		} else if l.matchString("/>") {
			depth--
			l.advanceN(2)
		} else {
			l.advance()
		}
	}
	if depth > 0 {
		l.error(line, col, "unterminated comment")
	}
}

func (l *Lexer) string(line, col int, quote rune) token.Token {
	l.advance()
	var b strings.Builder
	for !l.eof() && l.peek() != quote {
		if l.peek() == '\n' {
			l.error(line, col, "unterminated string")
			return token.Token{Type: token.STRING, Lexeme: b.String(), Line: line, Column: col}
		}
		if l.peek() == '\\' {
			l.advance()
			if l.eof() {
				break
			}
			switch l.peek() {
			case 'n':
				b.WriteRune('\n')
			case 't':
				b.WriteRune('\t')
			case '\\':
				b.WriteRune('\\')
			case '"':
				b.WriteRune('"')
			case '\'':
				b.WriteRune('\'')
			default:
				b.WriteRune(l.peek())
			}
			l.advance()
			continue
		}
		b.WriteRune(l.advance())
	}
	if l.eof() {
		l.error(line, col, "unterminated string")
	} else {
		l.advance()
	}
	return token.Token{Type: token.STRING, Lexeme: b.String(), Line: line, Column: col}
}

func (l *Lexer) multiString(line, col int) token.Token {
	l.advanceN(3)
	start := l.i
	for !l.eof() && !l.matchString(`"""`) {
		l.advance()
	}
	lex := string(l.src[start:l.i])
	if l.eof() {
		l.error(line, col, "unterminated multi-line string")
	} else {
		l.advanceN(3)
	}
	return token.Token{Type: token.MSTRING, Lexeme: lex, Line: line, Column: col}
}

func (l *Lexer) operator(line, col int) (token.Token, bool) {
	ops := []struct {
		s string
		t token.TokenType
	}{
		{"[[", token.DOUBLE_LBRACKET}, {"]]", token.DOUBLE_RBRACKET}, {"<=", token.LTE}, {">=", token.GTE},
		{"==", token.EQEQ}, {"!=", token.NOTEQ}, {"&&", token.AND}, {"||", token.OR}, {"::", token.DOUBLE_COLON}, {"=>", token.ARROW},
		{"<", token.LT}, {">", token.GT}, {"=", token.ASSIGN}, {"+", token.PLUS}, {"-", token.MINUS}, {"*", token.STAR},
		{"/", token.SLASH}, {"%", token.PERCENT}, {"!", token.NOT}, {"?", token.QUESTION}, {":", token.COLON}, {";", token.SEMICOLON},
		{",", token.COMMA}, {"(", token.LPAREN}, {")", token.RPAREN}, {"{", token.LBRACE}, {"}", token.RBRACE},
		{"[", token.LBRACKET}, {"]", token.RBRACKET},
	}
	for _, op := range ops {
		if l.matchString(op.s) {
			l.advanceN(len([]rune(op.s)))
			return token.Token{Type: op.t, Lexeme: op.s, Line: line, Column: col}, true
		}
	}
	return token.Token{}, false
}

func (l *Lexer) takeWhile(f func(rune) bool) string {
	start := l.i
	for !l.eof() && f(l.peek()) {
		l.advance()
	}
	return string(l.src[start:l.i])
}
func (l *Lexer) eof() bool { return l.i >= len(l.src) }
func (l *Lexer) peek() rune {
	if l.eof() {
		return 0
	}
	return l.src[l.i]
}
func (l *Lexer) advanceN(n int) {
	for j := 0; j < n; j++ {
		l.advance()
	}
}
func (l *Lexer) advance() rune {
	ch := l.src[l.i]
	l.i++
	if ch == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return ch
}
func (l *Lexer) matchString(s string) bool {
	r := []rune(s)
	if l.i+len(r) > len(l.src) {
		return false
	}
	for j := range r {
		if l.src[l.i+j] != r[j] {
			return false
		}
	}
	return true
}
func (l *Lexer) error(line, col int, msg string) {
	l.diagnostics = append(l.diagnostics, diagnostic.Diagnostic{Line: line, Column: col, Stage: "lexer", Message: msg})
}
func isIdentStart(r rune) bool { return unicode.IsLetter(r) || r == '_' }
func isIdentPart(r rune) bool  { return isIdentStart(r) || unicode.IsDigit(r) }
