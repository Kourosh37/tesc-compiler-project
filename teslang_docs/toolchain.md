# Toolchain

TesLang has three command-line tools.

## tes

Compile-if-needed and run.

```powershell
.\bin\tes.exe .\program.tes
```

Behavior:

1. Computes the generated path under `target/tesvm`.
2. If the `.tesvm` file is missing, compiles the `.tes` source.
3. If the `.tes` source is newer than the `.tesvm` file, recompiles.
4. Runs the generated `.tesvm` file with the TESVM runtime.

Useful flags:

```text
--force             always recompile before running
--trace             print VM instruction trace
--entry <name>      choose entry procedure, default main
--out-dir <dir>     choose generated TESVM directory
-v                  print cache/compile decisions
```

Examples:

```powershell
.\bin\tes.exe .\program.tes
.\bin\tes.exe --force .\program.tes
.\bin\tes.exe --trace .\program.tes
.\bin\tes.exe --out-dir .\build .\program.tes
```

## tesc

Compile only.

```powershell
.\bin\tesc.exe .\program.tes
```

Default output:

```text
target/tesvm/program.tesvm
```

Useful flags:

```text
--tokens        print lexer tokens
--check         syntax and semantic check only
--emit-tesvm    generate TESVM output
--stdout        print generated TESVM to stdout
-o <file>       write exact output file
--out-dir <dir> write generated output under a directory
```

## tesvm

Run generated TESVM code.

```powershell
.\bin\tesvm.exe .\target\tesvm\program.tesvm
```

Useful flags:

```text
--entry <name>  choose entry procedure, default main
--trace         print VM instruction trace
```

## Build Scripts

Windows:

```powershell
.\scripts\windows\build-all.ps1
.\scripts\windows\build-tes.ps1
.\scripts\windows\build-tesc.ps1
.\scripts\windows\build-tesvm.ps1
```

Linux:

```sh
./scripts/linux/build-all.sh
./scripts/linux/build-tes.sh
./scripts/linux/build-tesc.sh
./scripts/linux/build-tesvm.sh
```

macOS:

```sh
./scripts/macos/build-all.sh
./scripts/macos/build-tes.sh
./scripts/macos/build-tesc.sh
./scripts/macos/build-tesvm.sh
```

## Generated Files

```text
target/tesvm/**/*.tesvm
bin/
dist/
```

These are generated artifacts and are ignored by Git.

