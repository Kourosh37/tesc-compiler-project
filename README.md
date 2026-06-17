# TesLang Compiler

TesLang Compiler is a small educational compiler written in Go. It can compile `.tes` files into text TESVM-style intermediate code, print tokens, run syntax and semantic checks, and execute generated `.tesvm` files with the bundled VM.

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
6. Text TESVM code generator.
7. TESVM parser and virtual machine runtime.

## Run

PowerShell:

```powershell
go run ./cmd/teslang .\testdata\codegen_sample.tes
Get-Content .\testdata\lexer_sample.tes | go run ./cmd/teslang --tokens
Get-Content .\testdata\semantic_errors.tes | go run ./cmd/teslang --check
Get-Content .\testdata\codegen_sample.tes | go run ./cmd/teslang --emit-tesvm
```

Linux, macOS, Command Prompt, and Git Bash:

```sh
go run ./cmd/teslang testdata/codegen_sample.tes
go run ./cmd/teslang --tokens < testdata/lexer_sample.tes
go run ./cmd/teslang --check < testdata/semantic_errors.tes
go run ./cmd/teslang --emit-tesvm < testdata/codegen_sample.tes
```

Default mode is `--emit-tesvm`. When an input file is provided, the compiler writes generated TESVM under `target/tesvm`.

Useful file-based commands:

```sh
tes hello.tes
tes --force hello.tes
tes --trace hello.tes
tesc hello.tes
tesvm target/tesvm/hello.tesvm
tesc -o build/hello.tesvm hello.tes
tesc --out-dir build src/a.tes src/b.tes
tesc --stdout hello.tes
tesc --check src/a.tes src/b.tes
tesc --tokens hello.tes
```

For detailed build, run, cross-compilation, and output-file instructions, see [RUNBOOK.md](RUNBOOK.md).

For full language documentation, see [teslang_docs](teslang_docs/README.md).

## Build

Build with scripts:

```powershell
.\scripts\windows\build-all.ps1
.\scripts\windows\build-tes.ps1
.\scripts\windows\build-tesc.ps1
.\scripts\windows\build-tesvm.ps1
```

```sh
./scripts/linux/build-all.sh
./scripts/linux/build-tes.sh
./scripts/linux/build-tesc.sh
./scripts/linux/build-tesvm.sh

./scripts/macos/build-all.sh
./scripts/macos/build-tes.sh
./scripts/macos/build-tesc.sh
./scripts/macos/build-tesvm.sh
```

Or build manually:

```sh
go build -o bin/tes ./cmd/tes
go build -o bin/tesc ./cmd/teslang
go build -o bin/tesvm ./cmd/tesvm
```

On Windows PowerShell:

```powershell
New-Item -ItemType Directory -Force .\bin | Out-Null
go build -o .\bin\tes.exe .\cmd\tes
go build -o .\bin\tesc.exe .\cmd\teslang
go build -o .\bin\tesvm.exe .\cmd\tesvm
```

Generated binaries are written to `bin/`. Generated TESVM files are written to `target/tesvm` by default. Cross-compiled release binaries can be written to `dist/`; these generated directories are ignored by Git.

The main build scripts create `tes`, `tesc`, and `tesvm`. `tes` compiles a `.tes` source when needed and immediately runs the generated `.tesvm`.

`tes` recompiles when the generated `.tesvm` file is missing or older than the `.tes` source. Use `--force` to recompile unconditionally.

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

## Example TESVM Output

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
- Emits a readable text-based TESVM intermediate representation.
- Executes TESVM with a small register-based VM.
- Nested function code generation uses name mangling such as `outer__inner`.

## Limitations

- Nested functions are emitted with name mangling, but full closure conversion is not implemented.
- TESVM vector operations use VM-supported pseudo instructions such as `vector`, `loadidx`, and `storeidx`.
- Logical operators currently emit eager `and`/`or` instructions.
