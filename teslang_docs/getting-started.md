# Getting Started

## 1. Write a Program

Create `hello.tes`:

```tes
funk <null> main()
{
    print(123);
    return 0;
}
```

Every executable TesLang program should define a `main` function.

## 2. Build the Tools

Windows:

```powershell
.\scripts\windows\build-all.ps1
```

Linux:

```sh
./scripts/linux/build-all.sh
```

macOS:

```sh
./scripts/macos/build-all.sh
```

This creates:

```text
bin/tes     or bin/tes.exe      compile-if-needed and run
bin/tesc    or bin/tesc.exe     compile only
bin/tesvm   or bin/tesvm.exe    run TESVM files only
```

## 3. Run Directly

The easiest workflow is `tes`:

```powershell
.\bin\tes.exe .\hello.tes
```

On Linux/macOS:

```sh
./bin/tes ./hello.tes
```

`tes` compiles the file if needed, writes generated TESVM code under `target/tesvm`, and then runs it.

## 4. Compile Only

```powershell
.\bin\tesc.exe .\hello.tes
```

This writes:

```text
target/tesvm/hello.tesvm
```

## 5. Run Compiled TESVM

```powershell
.\bin\tesvm.exe .\target\tesvm\hello.tesvm
```

## 6. Check Without Running

```powershell
.\bin\tesc.exe --check .\hello.tes
```

Valid programs print:

```text
OK
```

## 7. Print Tokens

```powershell
.\bin\tesc.exe --tokens .\hello.tes
```

This is useful when debugging lexer behavior.

