package tesvm

import (
	"bytes"
	"strings"
	"testing"
)

func runProgram(t *testing.T, src string, input string) (string, int) {
	t.Helper()
	prog, err := Parse(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	code, err := NewVM(prog, WithInput(strings.NewReader(input)), WithOutput(&out)).Run("main")
	if err != nil {
		t.Fatal(err)
	}
	return out.String(), code
}

func TestVMArithmeticAndCall(t *testing.T) {
	out, code := runProgram(t, `
proc add
  add r3, r1, r2
  mov r0, r3
  ret

proc main
  mov r1, 10
  mov r2, 20
  call add, r3, r1, r2
  call log, r3
  mov r0, 0
  ret
`, "")
	if code != 0 || out != "30\n" {
		t.Fatalf("code=%d out=%q", code, out)
	}
}

func TestVMBranchesAndLoop(t *testing.T) {
	out, _ := runProgram(t, `
proc main
  mov r1, 0
  mov r2, 3
label loop
  lt r3, r1, r2
  jz r3, end
  mov r4, 1
  add r1, r1, r4
  jmp loop
label end
  call log, r1
  ret
`, "")
	if out != "3\n" {
		t.Fatalf("out=%q", out)
	}
}

func TestVMVectorOps(t *testing.T) {
	out, _ := runProgram(t, `
proc main
  call list, r1, 3
  storeidx r1, 0, 7
  loadidx r2, r1, 0
  call length, r3, r1
  add r4, r2, r3
  call log, r4
  ret
`, "")
	if out != "10\n" {
		t.Fatalf("out=%q", out)
	}
}

func TestVMRead(t *testing.T) {
	out, _ := runProgram(t, `
proc main
  call read, r1
  call log, r1
  ret
`, "42\n")
	if out != "42\n" {
		t.Fatalf("out=%q", out)
	}
}

func TestVMRandom(t *testing.T) {
	out, _ := runProgram(t, `
proc main
  call random, r1, 1, 1
  call log, r1
  ret
`, "")
	if out != "1\n" {
		t.Fatalf("out=%q", out)
	}
}
