package lexer

import (
	"teslang-compiler/internal/token"
	"testing"
)

func TestLexerNestedCommentsAndPositions(t *testing.T) {
	l := New("a = 10; </ x </ y /> z />\nb = 20;")
	toks, ds := l.LexAll()
	if len(ds) != 0 {
		t.Fatalf("diagnostics: %v", ds)
	}
	if toks[0].Lexeme != "a" || toks[4].Lexeme != "b" || toks[4].Line != 2 || toks[4].Column != 1 {
		t.Fatalf("bad tokens: %#v", toks[:6])
	}
}
func TestLexerStringsAndMultiStrings(t *testing.T) {
	l := New("'a\\n' \"b\" \"\"\"x\ny\"\"\"")
	toks, ds := l.LexAll()
	if len(ds) != 0 {
		t.Fatalf("diagnostics: %v", ds)
	}
	if toks[0].Type != token.STRING || toks[1].Type != token.STRING || toks[2].Type != token.MSTRING {
		t.Fatalf("bad string tokens: %#v", toks[:3])
	}
}
func TestLexerUnterminated(t *testing.T) {
	_, ds := New("</ no end").LexAll()
	if len(ds) == 0 {
		t.Fatal("expected unterminated comment")
	}
	_, ds = New("'no end").LexAll()
	if len(ds) == 0 {
		t.Fatal("expected unterminated string")
	}
}
