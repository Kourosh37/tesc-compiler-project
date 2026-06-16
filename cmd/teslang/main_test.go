package main

import "testing"

func TestOutputPathDefaultsNextToInput(t *testing.T) {
	got, err := outputPath("samples/hello.tes", options{mode: modeEmitTSVM})
	if err != nil {
		t.Fatal(err)
	}
	if got != "samples\\hello.tsvm" && got != "samples/hello.tsvm" {
		t.Fatalf("unexpected output path: %q", got)
	}
}

func TestValidateOptionsRejectsInvalidOutputSelection(t *testing.T) {
	err := validateOptions(options{mode: modeEmitTSVM, output: "a.tsvm", outDir: "out"}, []string{"a.tes"})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidateOptionsRequiresSingleInputForOutputFile(t *testing.T) {
	err := validateOptions(options{mode: modeEmitTSVM, output: "out.tsvm"}, []string{"a.tes", "b.tes"})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
