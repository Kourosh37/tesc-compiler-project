package semantic

import (
	"strings"
	"teslang-compiler/internal/lexer"
	"teslang-compiler/internal/parser"
	"testing"
)

func analyze(t *testing.T, src string) string {
	t.Helper()
	toks, lds := lexer.New(src).LexAll()
	if len(lds) > 0 {
		t.Fatalf("lexer: %v", lds)
	}
	p := parser.New(toks)
	prog, pds := p.ParseProgram()
	if len(pds) > 0 {
		t.Fatalf("parser: %v", pds)
	}
	ds := New().Analyze(prog)
	var b strings.Builder
	for _, d := range ds {
		b.WriteString(d.Message)
		b.WriteByte('\n')
	}
	return b.String()
}
func TestSemanticValid(t *testing.T) {
	if got := analyze(t, `funk <null> main(){ a :: int = 1; print(a); return 0; }`); got != "" {
		t.Fatal(got)
	}
}

func TestSemanticRandomBuiltin(t *testing.T) {
	if got := analyze(t, `funk <null> main(){ x :: int = random(1, 10); print(x); return 0; }`); got != "" {
		t.Fatal(got)
	}
}

func TestSemanticErrors(t *testing.T) {
	got := analyze(t, `funk <int> f(A as vector){ k :: int; x :: int; x :: vector; x = list(3); if [[ 1 ]] begin return k; endif return A[0]; } funk <null> main(){ a :: int; print(f(a, 1)); return a; }`)
	for _, want := range []string{"used before being assigned", "duplicate variable", "expected to be of type 'int' but got 'vector'", "condition expression must be bool", "expects 1 arguments but got 2", "expected argument 'A' to be of type 'vector' but got 'int'", "wrong return type"} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q in:\n%s", want, got)
		}
	}
}
