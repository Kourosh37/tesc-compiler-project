# TesLang Compiler Samples

This folder contains small `.tes` files used to demonstrate compiler stages.

## Samples

- `lexer_sample.tes`: tokenization sample with nested comments, loops, function calls, and vector indexing.
- `codegen_sample.tes`: valid compile-and-run sample that reads two integers and prints their sum.
- `semantic_errors.tes`: intentionally invalid sample for checking semantic diagnostics.

## Run

Tokenize:

```powershell
.\bin\tesc.exe --tokens .\examples\compiler\lexer_sample.tes
```

Check semantic errors:

```powershell
.\bin\tesc.exe --check .\examples\compiler\semantic_errors.tes
```

Compile and run:

```powershell
"5`n7`n" | .\bin\tes.exe .\examples\compiler\codegen_sample.tes
```

