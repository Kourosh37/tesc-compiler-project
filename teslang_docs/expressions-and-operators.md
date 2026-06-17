# Expressions and Operators

## Literals

```tes
123
"hello"
'hello'
"""
multi-line
"""
true
false
[1, 2, 3]
```

## Identifier

```tes
x
total
```

## Function Call

```tes
add(1, 2)
scan()
print(x)
```

## Index Access

```tes
items[0]
items[i]
```

The target must be `vector`, and the index must be `int`.

## Assignment

```tes
x = 10
items[0] = 99
```

Assignment is right-associative.

## Operator Precedence

Highest to lowest:

| Level | Operators / Forms |
| --- | --- |
| 1 | function call `f(args)`, index `a[i]` |
| 2 | unary `!`, `+`, `-` |
| 3 | `*`, `/`, `%` |
| 4 | `+`, `-` |
| 5 | `<`, `>`, `<=`, `>=` |
| 6 | `==`, `!=` |
| 7 | `&&` |
| 8 | `||` |
| 9 | ternary `condition ? a : b` |
| 10 | assignment `=` |

## Arithmetic

```tes
a + b
a - b
a * b
a / b
a % b
```

For integers:

```text
int + int -> int
int - int -> int
int * int -> int
int / int -> int
int % int -> int
```

String concatenation:

```text
str + str -> str
str + mstr -> mstr
mstr + str -> mstr
mstr + mstr -> mstr
```

## Comparisons

```tes
a < b
a > b
a <= b
a >= b
```

These require integer operands and return `bool`.

## Equality

```tes
a == b
a != b
```

Both sides should have compatible types.

## Logical Operators

```tes
a && b
a || b
!a
```

Operands must be `bool`.

## Ternary

```tes
result :: int = ok ? 1 : 0;
```

Rules:

- The condition must be `bool`.
- Both branches must have compatible types.

