package codegen

import (
	"strings"
	"teslang-compiler/internal/lexer"
	"teslang-compiler/internal/parser"
	"teslang-compiler/internal/semantic"
	"testing"
)

func TestCodegenSample(t *testing.T) {
	src := `funk <int> add(a as int, b as int){ result :: int = a + b; return result; } funk <null> main(){ a :: int; b :: int; a = scan(); b = scan(); print(add(a,b)); return 0; }`
	toks, lds := lexer.New(src).LexAll()
	if len(lds) > 0 {
		t.Fatal(lds)
	}
	p := parser.New(toks)
	prog, pds := p.ParseProgram()
	if len(pds) > 0 {
		t.Fatal(pds)
	}
	if ds := semantic.New().Analyze(prog); len(ds) > 0 {
		t.Fatal(ds)
	}
	out, ds := New().Generate(prog)
	if len(ds) > 0 {
		t.Fatal(ds)
	}
	for _, want := range []string{"proc add", "add r3, r1, r2", "call read", "call log", "mov r0, 0"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in:\n%s", want, out)
		}
	}
}
