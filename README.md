# TesLang Compiler

TesLang Compiler is a small educational compiler written in Go. It can compile `.tes` files into text TSVM-style intermediate code, print tokens, and run syntax and semantic checks.

## Language Features

- Functions with `funk <type> name(...) { ... }`
- Types: `int`, `vector`, `str`, `mstr`, `bool`, `null`
- Variables declared with `name :: type`
- `if`, `while`, `do while`, and `for`
- Nested functions
- Integer arithmetic, comparisons, logical operators, ternary expressions
- Vector literals and indexing
- Built-ins: `scan`, `print`, `list`, `length`, `exit`

## Compiler Pipeline

1. Manual lexer with nested `</ ... />` comments and line/column tracking.
2. Recursive descent parser for declarations and statements.
3. Pratt parser for expressions.
4. AST construction.
5. Semantic analyzer with nested symbol tables.
6. Text TSVM code generator.

## Run

PowerShell:

```powershell
go run ./cmd/teslang .\testdata\codegen_sample.tes
Get-Content .\testdata\lexer_sample.tes | go run ./cmd/teslang --tokens
Get-Content .\testdata\semantic_errors.tes | go run ./cmd/teslang --check
Get-Content .\testdata\codegen_sample.tes | go run ./cmd/teslang --emit-tsvm
```

Linux, macOS, Command Prompt, and Git Bash:

```sh
go run ./cmd/teslang testdata/codegen_sample.tes
go run ./cmd/teslang --tokens < testdata/lexer_sample.tes
go run ./cmd/teslang --check < testdata/semantic_errors.tes
go run ./cmd/teslang --emit-tsvm < testdata/codegen_sample.tes
```

Default mode is `--emit-tsvm`. When an input file is provided, the compiler writes a `.tsvm` file next to it.

Useful file-based commands:

```sh
tesc hello.tes
tesc -o build/hello.tsvm hello.tes
tesc --out-dir build src/a.tes src/b.tes
tesc --stdout hello.tes
tesc --check src/a.tes src/b.tes
tesc --tokens hello.tes
```

For detailed build, run, cross-compilation, and output-file instructions, see [RUNBOOK.md](RUNBOOK.md).

## Build

Build the compiler for the current operating system:

```sh
go build -o bin/tesc ./cmd/teslang
```

On Windows PowerShell:

```powershell
New-Item -ItemType Directory -Force .\bin | Out-Null
go build -o .\bin\tesc.exe .\cmd\teslang
```

Generated binaries are written to `bin/`. Cross-compiled release binaries can be written to `dist/`; both directories are ignored by Git.

## Test

```powershell
go test ./...
```

## Example Token Output

```text
Line | Column | Token | Value
1 | 1 | FUNK | funk
1 | 6 | LT | <
1 | 7 | INT | int
```

## Example Semantic Error

```text
Error [semantic] line 24, column 5 in function 'main': variable 'A' expected to be of type 'int' but got 'vector'.
```

## Example TSVM Output

```text
proc add
  add r3, r1, r2
  mov r0, r3
  ret

proc main
  call read, r1
  call read, r2
  call add, r3, r1, r2
  call log, r3
  mov r0, 0
  ret
```

## Design Decisions

- Implemented the lexer manually.
- Used recursive descent parsing for program structure.
- Used a Pratt parser for expressions.
- Used nested symbol tables for lexical scope.
- Used `null` as the void-like type.
- Treats `vector` as `array<int>`.
- Emits a readable text-based TSVM intermediate representation.
- Nested function code generation uses name mangling such as `outer__inner`.

## Limitations

- Nested functions are emitted with name mangling, but full closure conversion is not implemented.
- TSVM vector operations use pseudo instructions such as `vector`, `loadidx`, and `storeidx`.
- Logical operators currently emit eager `and`/`or` instructions.
