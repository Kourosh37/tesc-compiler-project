package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"teslang-compiler/internal/codegen"
	"teslang-compiler/internal/diagnostic"
	"teslang-compiler/internal/lexer"
	"teslang-compiler/internal/parser"
	"teslang-compiler/internal/semantic"
	"teslang-compiler/internal/token"
)

const defaultOutputDir = "target/tesvm"

type mode int

const (
	modeTokens mode = iota
	modeCheck
	modeEmitTESVM
)

type options struct {
	mode   mode
	output string
	outDir string
	stdout bool
}

func main() {
	tokensFlag := flag.Bool("tokens", false, "tokenize only")
	checkFlag := flag.Bool("check", false, "check syntax and semantics")
	emitFlag := flag.Bool("emit-tesvm", false, "emit TESVM")
	outputFlag := flag.String("o", "", "output file path; only valid with one input file")
	outDirFlag := flag.String("out-dir", "", "output directory for generated files")
	stdoutFlag := flag.Bool("stdout", false, "write generated output to stdout")
	flag.Parse()

	selectedModes := 0
	for _, selected := range []bool{*tokensFlag, *checkFlag, *emitFlag} {
		if selected {
			selectedModes++
		}
	}
	if selectedModes > 1 {
		fmt.Fprintln(os.Stderr, "use only one of --tokens, --check, or --emit-tesvm")
		os.Exit(2)
	}

	opts := options{mode: modeEmitTESVM, output: *outputFlag, outDir: *outDirFlag, stdout: *stdoutFlag}
	if *tokensFlag {
		opts.mode = modeTokens
	}
	if *checkFlag {
		opts.mode = modeCheck
	}
	if !*tokensFlag && !*checkFlag && !*emitFlag {
		opts.mode = modeEmitTESVM
	}
	if *emitFlag {
		opts.mode = modeEmitTESVM
	}
	if err := validateOptions(opts, flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if err := run(opts, flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func validateOptions(opts options, inputs []string) error {
	selectedOutputs := 0
	if opts.output != "" {
		selectedOutputs++
	}
	if opts.outDir != "" {
		selectedOutputs++
	}
	if opts.stdout {
		selectedOutputs++
	}
	if selectedOutputs > 1 {
		return fmt.Errorf("use only one of -o, --out-dir, or --stdout")
	}
	if opts.output != "" && len(inputs) != 1 {
		return fmt.Errorf("-o requires exactly one input file")
	}
	if (opts.output != "" || opts.outDir != "") && opts.mode != modeEmitTESVM {
		return fmt.Errorf("-o and --out-dir are only valid with --emit-tesvm")
	}
	return nil
}

func run(opts options, inputs []string) error {
	if len(inputs) == 0 {
		src, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		return process("<stdin>", string(src), "", opts, true)
	}
	failed := false
	for _, input := range inputs {
		src, err := os.ReadFile(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", input, err)
			failed = true
			continue
		}
		outPath, err := outputPath(input, opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", input, err)
			failed = true
			continue
		}
		if err := process(input, string(src), outPath, opts, len(inputs) == 1); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", input, err)
			failed = true
		}
	}
	if failed {
		return fmt.Errorf("one or more inputs failed")
	}
	return nil
}

func outputPath(input string, opts options) (string, error) {
	if opts.mode != modeEmitTESVM || opts.stdout {
		return "", nil
	}
	if opts.output != "" {
		return opts.output, nil
	}
	name := outputName(input)
	if opts.outDir != "" {
		if err := os.MkdirAll(opts.outDir, 0755); err != nil {
			return "", err
		}
		return filepath.Join(opts.outDir, name), nil
	}
	return filepath.Join(defaultOutputDir, name), nil
}

func outputName(input string) string {
	clean := filepath.Clean(input)
	if rel, err := filepath.Rel(".", clean); err == nil && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && rel != ".." && !filepath.IsAbs(rel) {
		clean = rel
	} else {
		clean = filepath.Base(clean)
	}
	ext := filepath.Ext(clean)
	return strings.TrimSuffix(clean, ext) + ".tesvm"
}

func process(name, src, outPath string, opts options, single bool) error {
	toks, ds := lexer.New(src).LexAll()
	if opts.mode == modeTokens {
		if !single {
			fmt.Printf("== %s ==\n", name)
		}
		printTokens(toks)
		if len(ds) > 0 {
			printDiagnostics(ds)
			return fmt.Errorf("lexing failed")
		}
		return nil
	}
	if len(ds) > 0 {
		printDiagnostics(ds)
		return fmt.Errorf("lexing failed")
	}
	prog, pds := parser.New(toks).ParseProgram()
	if len(pds) > 0 {
		printDiagnostics(pds)
		return fmt.Errorf("parsing failed")
	}
	sds := semantic.New().Analyze(prog)
	if len(sds) > 0 {
		printDiagnostics(sds)
		return fmt.Errorf("semantic analysis failed")
	}
	if opts.mode == modeCheck {
		if single {
			fmt.Println("OK")
		} else {
			fmt.Printf("%s: OK\n", name)
		}
		return nil
	}
	out, cds := codegen.New().Generate(prog)
	if len(cds) > 0 {
		printDiagnostics(cds)
		return fmt.Errorf("code generation failed")
	}
	if opts.stdout || outPath == "" {
		fmt.Print(out)
		if !strings.HasSuffix(out, "\n") {
			fmt.Println()
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(outPath, []byte(out), 0644); err != nil {
		return err
	}
	fmt.Printf("%s -> %s\n", name, outPath)
	return nil
}

func printTokens(toks []token.Token) {
	fmt.Println("Line | Column | Token | Value")
	for _, t := range toks {
		if t.Type != token.EOF {
			fmt.Printf("%d | %d | %s | %s\n", t.Line, t.Column, t.Type, t.Lexeme)
		}
	}
}

func printDiagnostics(ds []diagnostic.Diagnostic) {
	fmt.Println(diagnostic.FormatAll(ds))
}
