# Diagnostics and Common Mistakes

The compiler reports diagnostics with stage, line, and column.

Example:

```text
Error [semantic] line 8, column 13 in function 'find': variable 'k' is used before being assigned.
```

## Undefined Variable

Bad:

```tes
funk <null> main()
{
    print(x);
    return 0;
}
```

Fix:

```tes
funk <null> main()
{
    x :: int = 1;
    print(x);
    return 0;
}
```

## Uninitialized Variable

Bad:

```tes
x :: int;
print(x);
```

Fix:

```tes
x :: int = 0;
print(x);
```

or:

```tes
x :: int;
x = 0;
print(x);
```

## Duplicate Variable

Bad:

```tes
x :: int;
x :: vector;
```

Fix: use a different name or remove the duplicate declaration.

## Assignment Type Mismatch

Bad:

```tes
x :: int;
x = list(3);
```

Fix:

```tes
x :: vector = list(3);
```

## Wrong Argument Count

Bad:

```tes
funk <int> add(a as int, b as int)
{
    return a + b;
}

funk <null> main()
{
    print(add(1));
    return 0;
}
```

Fix:

```tes
print(add(1, 2));
```

## Wrong Argument Type

Bad:

```tes
funk <int> first(items as vector)
{
    return items[0];
}

funk <null> main()
{
    x :: int = 1;
    print(first(x));
    return 0;
}
```

Fix:

```tes
items :: vector = [1, 2, 3];
print(first(items));
```

## Wrong Return Type

Bad:

```tes
funk <vector> f()
{
    return 1;
}
```

Fix:

```tes
funk <int> f()
{
    return 1;
}
```

## Non-bool Condition

Bad:

```tes
if [[ 1 ]]
begin
    print(1);
endif
```

Fix:

```tes
if [[ 1 == 1 ]]
begin
    print(1);
endif
```

## Invalid Vector Indexing

Bad:

```tes
x :: int = 1;
print(x[0]);
```

Fix:

```tes
items :: vector = [1, 2, 3];
print(items[0]);
```

## Unterminated Comment

Bad:

```tes
</ missing end
```

Fix:

```tes
</ closed />
```

## Unterminated String

Bad:

```tes
"hello
```

Fix:

```tes
"hello"
```

