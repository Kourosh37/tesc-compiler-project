package parser

import (
	"teslang-compiler/internal/ast"
	"teslang-compiler/internal/lexer"
	"testing"
)

func parse(t *testing.T, src string) *ast.Program {
	t.Helper()
	toks, lds := lexer.New(src).LexAll()
	if len(lds) > 0 {
		t.Fatalf("lexer: %v", lds)
	}
	p := New(toks)
	prog, ds := p.ParseProgram()
	if len(ds) > 0 {
		t.Fatalf("parser: %v", ds)
	}
	return prog
}
func TestParserFunctionAndStatements(t *testing.T) {
	prog := parse(t, `funk <int> f(a as int){ x :: int = 1; if [[ true ]] begin x = x + a; else x = 0; endif while [[ true ]] begin x = x - 1; endwhile do begin x = x + 1; while [[ false ]] endwhile for (i = 0 to 3) begin x = x + i; endfor return x; }`)
	if len(prog.Functions) != 1 || len(prog.Functions[0].Body.Statements) != 6 {
		t.Fatalf("bad program: %#v", prog.Functions[0].Body.Statements)
	}
}
func TestParserNestedFunctionAndPrecedence(t *testing.T) {
	prog := parse(t, `funk <int> outer(){ funk <int> inner()=> return 1 + 2 * 3; return inner(); }`)
	if _, ok := prog.Functions[0].Body.Statements[0].(*ast.FunctionDecl); !ok {
		t.Fatal("expected nested function")
	}
}
