# Built-in Functions

TesLang provides a small set of built-in functions.

## scan

Reads an integer from input.

```tes
x :: int = scan();
```

Type:

```text
scan() -> int
```

## print

Prints a value.

```tes
print(123);
print("hello");
print(true);
```

Type:

```text
print(x) -> null
```

Supported argument types:

```text
int
str
mstr
bool
```

Vectors are rejected by semantic analysis unless compiler support is extended.

## list

Creates a vector of integer zeros.

```tes
items :: vector = list(3);
```

Type:

```text
list(n: int) -> vector
```

Example:

```tes
items :: vector = list(3);
items[0] = 10;
items[1] = 20;
items[2] = 30;
```

## length

Returns vector length.

```tes
n :: int = length(items);
```

Type:

```text
length(arr: vector) -> int
```

## random

Returns a random integer in an inclusive range.

```tes
secret :: int = random(1, 100);
```

Type:

```text
random(min: int, max: int) -> int
```

Rules:

- `min` must be `int`.
- `max` must be `int`.
- `min` must be less than or equal to `max`.
- The result can be equal to `min` or `max`.

## exit

Stops the VM with an exit code.

```tes
exit(1);
```

Type:

```text
exit(n: int) -> null
```
