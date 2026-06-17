package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"teslang-compiler/internal/tsvm"
)

func main() {
	entry := flag.String("entry", "main", "entry procedure")
	trace := flag.Bool("trace", false, "print executed instructions")
	flag.Parse()

	var r io.Reader
	name := "<stdin>"
	if flag.NArg() == 0 {
		r = os.Stdin
	} else if flag.NArg() == 1 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
		r = f
		name = flag.Arg(0)
	} else {
		fmt.Fprintln(os.Stderr, "usage: tsvm [--entry main] [--trace] [file.tsvm]")
		os.Exit(2)
	}

	prog, err := tsvm.Parse(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
		os.Exit(1)
	}
	code, err := tsvm.NewVM(prog, tsvm.WithInput(os.Stdin), tsvm.WithOutput(os.Stdout), tsvm.WithTrace(*trace)).Run(*entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
		os.Exit(1)
	}
	os.Exit(code)
}
