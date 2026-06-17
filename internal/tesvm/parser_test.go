package tesvm

import (
	"strings"
	"testing"
)

func TestParseProgram(t *testing.T) {
	src := `
proc main
  mov r1, 10
  label loop
  jmp loop
`
	prog, err := Parse(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	main := prog.Procedures["main"]
	if main == nil || len(main.Instructions) != 2 || main.Labels["loop"] != 1 {
		t.Fatalf("bad program: %#v", prog)
	}
}

func TestSplitOperandsWithQuotedComma(t *testing.T) {
	ops, err := splitOperands(`r1, "a,b", r2`)
	if err != nil {
		t.Fatal(err)
	}
	if len(ops) != 3 || ops[1] != `"a,b"` {
		t.Fatalf("bad operands: %#v", ops)
	}
}
