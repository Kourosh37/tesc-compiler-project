# TesLang Documentation

TesLang is a small educational programming language designed for compiler design projects. This documentation explains how to write TesLang programs, compile them, and run them with the TESVM runtime.

## Documentation Map

- [Getting Started](getting-started.md): write, compile, and run your first program.
- [Language Reference](language-reference.md): complete syntax and language rules.
- [Types and Values](types-and-values.md): `int`, `bool`, `str`, `mstr`, `vector`, and `null`.
- [Functions and Scope](functions-and-scope.md): functions, parameters, returns, nested functions, and scope.
- [Statements and Control Flow](statements-and-control-flow.md): declarations, assignments, `if`, loops, and returns.
- [Expressions and Operators](expressions-and-operators.md): precedence, calls, indexing, ternary expressions, and assignment.
- [Built-in Functions](builtins.md): `scan`, `print`, `list`, `length`, `random`, and `exit`.
- [Examples](examples.md): complete sample programs.
- [Diagnostics and Common Mistakes](diagnostics.md): common compiler errors and fixes.
- [Toolchain](toolchain.md): `tes`, `tesc`, `tesvm`, build scripts, and generated files.

## Minimal Program

```tes
funk <null> main()
{
    print(123);
    return 0;
}
```

Run it:

```powershell
.\bin\tes.exe .\hello.tes
```

Expected output:

```text
123
```
