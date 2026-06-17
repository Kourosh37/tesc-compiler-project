# Statements and Control Flow

## Expression Statement

```tes
print(1);
x = 10;
```

## Variable Declaration

```tes
x :: int;
y :: int = 10;
```

## Return Statement

```tes
return x;
```

## if / else

Syntax:

```tes
if [[ condition ]]
begin
    statements
endif
```

With `else`:

```tes
if [[ x > 0 ]]
begin
    print(1);
else
    print(0);
endif
```

The condition must be `bool`.

## while

```tes
while [[ i < 10 ]]
begin
    print(i);
    i = i + 1;
endwhile
```

The condition must be `bool`.

## do while

```tes
do
begin
    print(i);
    i = i + 1;
while [[ i < 10 ]]
endwhile
```

The body runs before the condition is checked.

## for

```tes
for (i = 0 to 10)
begin
    print(i);
endfor
```

Rules:

- The loop variable is an `int`.
- The start expression must be `int`.
- The end expression must be `int`.
- The generated loop runs while `i < end`.

Example:

```tes
for (i = 0 to 3)
begin
    print(i);
endfor
```

Output:

```text
0
1
2
```

