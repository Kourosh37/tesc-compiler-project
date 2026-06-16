package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"teslang-compiler/internal/codegen"
	"teslang-compiler/internal/diagnostic"
	"teslang-compiler/internal/lexer"
	"teslang-compiler/internal/parser"
	"teslang-compiler/internal/semantic"
	"teslang-compiler/internal/token"
)

func main() {
	tokensFlag := flag.Bool("tokens", false, "tokenize only")
	checkFlag := flag.Bool("check", false, "check syntax and semantics")
	emitFlag := flag.Bool("emit-tsvm", false, "emit TSVM")
	flag.Parse()
	if !*tokensFlag && !*checkFlag && !*emitFlag {
		*emitFlag = true
	}
	src, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	l := lexer.New(string(src))
	toks, ds := l.LexAll()
	if *tokensFlag {
		fmt.Println("Line | Column | Token | Value")
		for _, t := range toks {
			if t.Type != token.EOF {
				fmt.Printf("%d | %d | %s | %s\n", t.Line, t.Column, t.Type, t.Lexeme)
			}
		}
		if len(ds) > 0 {
			fmt.Println(diagnostic.FormatAll(ds))
			os.Exit(1)
		}
		return
	}
	if len(ds) > 0 {
		fmt.Println(diagnostic.FormatAll(ds))
		os.Exit(1)
	}
	p := parser.New(toks)
	prog, pds := p.ParseProgram()
	if len(pds) > 0 {
		fmt.Println(diagnostic.FormatAll(pds))
		os.Exit(1)
	}
	sds := semantic.New().Analyze(prog)
	if len(sds) > 0 {
		fmt.Println(diagnostic.FormatAll(sds))
		os.Exit(1)
	}
	if *checkFlag {
		fmt.Println("OK")
		return
	}
	out, cds := codegen.New().Generate(prog)
	if len(cds) > 0 {
		fmt.Println(diagnostic.FormatAll(cds))
		os.Exit(1)
	}
	fmt.Print(out)
}
