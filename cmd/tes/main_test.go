package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestOutputPathUsesTargetTree(t *testing.T) {
	got := outputPath("src/main.tes", "target/tesvm")
	if got != "target\\tesvm\\src\\main.tesvm" && got != "target/tesvm/src/main.tesvm" {
		t.Fatalf("unexpected output path: %q", got)
	}
}

func TestNeedsCompileMissingOutput(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "main.tes")
	if err := os.WriteFile(src, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}
	needs, err := needsCompile(src, filepath.Join(dir, "main.tesvm"), false)
	if err != nil {
		t.Fatal(err)
	}
	if !needs {
		t.Fatal("expected missing output to require compile")
	}
}

func TestNeedsCompileStaleOutput(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "main.tes")
	out := filepath.Join(dir, "main.tesvm")
	if err := os.WriteFile(out, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Millisecond)
	if err := os.WriteFile(src, []byte("new"), 0644); err != nil {
		t.Fatal(err)
	}
	needs, err := needsCompile(src, out, false)
	if err != nil {
		t.Fatal(err)
	}
	if !needs {
		t.Fatal("expected stale output to require compile")
	}
}

func TestCompileSource(t *testing.T) {
	out, err := compileSource(`funk <null> main(){ print(1 + 2); return 0; }`)
	if err != nil {
		t.Fatal(err)
	}
	if out == "" {
		t.Fatal("expected generated code")
	}
}
