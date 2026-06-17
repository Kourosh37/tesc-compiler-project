package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"teslang-compiler/internal/codegen"
	"teslang-compiler/internal/diagnostic"
	"teslang-compiler/internal/lexer"
	"teslang-compiler/internal/parser"
	"teslang-compiler/internal/semantic"
	"teslang-compiler/internal/tesvm"
)

const defaultOutputDir = "target/tesvm"

type options struct {
	entry   string
	outDir  string
	force   bool
	trace   bool
	verbose bool
}

func main() {
	opts := options{}
	flag.StringVar(&opts.entry, "entry", "main", "entry procedure")
	flag.StringVar(&opts.outDir, "out-dir", defaultOutputDir, "TESVM output directory")
	flag.BoolVar(&opts.force, "force", false, "force recompilation before running")
	flag.BoolVar(&opts.trace, "trace", false, "trace VM instruction execution")
	flag.BoolVar(&opts.verbose, "v", false, "print compile/cache decisions")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: tes [--entry main] [--out-dir target/tesvm] [--force] [--trace] [-v] file.tes")
		os.Exit(2)
	}
	code, err := runFile(flag.Arg(0), opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(code)
}

func runFile(input string, opts options) (int, error) {
	outPath := outputPath(input, opts.outDir)
	needs, err := needsCompile(input, outPath, opts.force)
	if err != nil {
		return 1, err
	}
	if needs {
		if opts.verbose {
			fmt.Fprintf(os.Stderr, "compiling %s -> %s\n", input, outPath)
		}
		if err := compileFile(input, outPath); err != nil {
			return 1, err
		}
	} else if opts.verbose {
		fmt.Fprintf(os.Stderr, "using cached %s\n", outPath)
	}
	return executeFile(outPath, opts)
}

func needsCompile(sourcePath, outputPath string, force bool) (bool, error) {
	if force {
		return true, nil
	}
	srcInfo, err := os.Stat(sourcePath)
	if err != nil {
		return false, err
	}
	outInfo, err := os.Stat(outputPath)
	if os.IsNotExist(err) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return srcInfo.ModTime().After(outInfo.ModTime()), nil
}

func compileFile(input, outPath string) error {
	src, err := os.ReadFile(input)
	if err != nil {
		return err
	}
	out, err := compileSource(string(src))
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(outPath, []byte(out), 0644)
}

func compileSource(src string) (string, error) {
	toks, lds := lexer.New(src).LexAll()
	if len(lds) > 0 {
		return "", errors.New(diagnostic.FormatAll(lds))
	}
	prog, pds := parser.New(toks).ParseProgram()
	if len(pds) > 0 {
		return "", errors.New(diagnostic.FormatAll(pds))
	}
	sds := semantic.New().Analyze(prog)
	if len(sds) > 0 {
		return "", errors.New(diagnostic.FormatAll(sds))
	}
	out, cds := codegen.New().Generate(prog)
	if len(cds) > 0 {
		return "", errors.New(diagnostic.FormatAll(cds))
	}
	return out, nil
}

func executeFile(path string, opts options) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 1, err
	}
	defer f.Close()
	prog, err := tesvm.Parse(f)
	if err != nil {
		return 1, fmt.Errorf("%s: %w", path, err)
	}
	return tesvm.NewVM(prog, tesvm.WithInput(os.Stdin), tesvm.WithOutput(os.Stdout), tesvm.WithTrace(opts.trace)).Run(opts.entry)
}

func outputPath(input, outDir string) string {
	return filepath.Join(outDir, outputName(input))
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
