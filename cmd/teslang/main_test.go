package main

import "testing"

func TestOutputPathDefaultsToTargetDirectory(t *testing.T) {
	got, err := outputPath("samples/hello.tes", options{mode: modeEmitTESVM})
	if err != nil {
		t.Fatal(err)
	}
	if got != "target\\tesvm\\samples\\hello.tesvm" && got != "target/tesvm/samples/hello.tesvm" {
		t.Fatalf("unexpected output path: %q", got)
	}
}

func TestOutputPathUsesCustomOutputDirectory(t *testing.T) {
	got, err := outputPath("samples/hello.tes", options{mode: modeEmitTESVM, outDir: "out"})
	if err != nil {
		t.Fatal(err)
	}
	if got != "out\\samples\\hello.tesvm" && got != "out/samples/hello.tesvm" {
		t.Fatalf("unexpected output path: %q", got)
	}
}

func TestValidateOptionsRejectsInvalidOutputSelection(t *testing.T) {
	err := validateOptions(options{mode: modeEmitTESVM, output: "a.tesvm", outDir: "out"}, []string{"a.tes"})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidateOptionsRequiresSingleInputForOutputFile(t *testing.T) {
	err := validateOptions(options{mode: modeEmitTESVM, output: "out.tesvm"}, []string{"a.tes", "b.tes"})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
