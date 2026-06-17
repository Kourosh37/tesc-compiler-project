# Types and Values

TesLang has these types:

```text
int
bool
str
mstr
vector
null
```

## int

Integer values.

```tes
x :: int = 42;
y :: int = -1;
```

Supported arithmetic:

```tes
a + b
a - b
a * b
a / b
a % b
```

## bool

Boolean values:

```tes
ok :: bool = true;
done :: bool = false;
```

`boolean` is accepted as an alias in type positions:

```tes
flag :: boolean = true;
```

## str

Single-line string:

```tes
name :: str = "TesLang";
```

Single quotes are also supported:

```tes
name :: str = 'TesLang';
```

Escapes:

```text
\"   quote
\'   quote
\\   backslash
\n   newline
\t   tab
```

## mstr

Multi-line string using triple quotes:

```tes
text :: mstr = """
line one
line two
""";
```

## vector

`vector` is an array of integers.

Literal:

```tes
nums :: vector = [1, 2, 3];
```

Created with `list`:

```tes
nums :: vector = list(3);
nums[0] = 10;
nums[1] = 20;
nums[2] = 30;
```

Indexing returns an `int`:

```tes
x :: int = nums[1];
```

## null

`null` is the void-like function return type.

```tes
funk <null> logNumber(x as int)
{
    print(x);
    return 0;
}
```

`null` is used for functions that do not produce a meaningful value.

