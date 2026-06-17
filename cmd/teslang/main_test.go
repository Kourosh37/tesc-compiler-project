package main

import "testing"

func TestOutputPathDefaultsToTargetDirectory(t *testing.T) {
	got, err := outputPath("samples/hello.tes", options{mode: modeEmitTSVM})
	if err != nil {
		t.Fatal(err)
	}
	if got != "target\\tsvm\\samples\\hello.tsvm" && got != "target/tsvm/samples/hello.tsvm" {
		t.Fatalf("unexpected output path: %q", got)
	}
}

func TestOutputPathUsesCustomOutputDirectory(t *testing.T) {
	got, err := outputPath("samples/hello.tes", options{mode: modeEmitTSVM, outDir: "out"})
	if err != nil {
		t.Fatal(err)
	}
	if got != "out\\samples\\hello.tsvm" && got != "out/samples/hello.tsvm" {
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
