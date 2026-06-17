# Language Reference

## Program Structure

A TesLang file contains zero or more function declarations.

```tes
funk <int> add(a as int, b as int)
{
    return a + b;
}

funk <null> main()
{
    print(add(10, 20));
    return 0;
}
```

There are no global variables. All variables live inside functions or nested blocks.

## Comments

Comments start with `</` and end with `/>`.

```tes
</ this is a comment />
```

Comments can span multiple lines:

```tes
</
    multi-line comment
/>
```

Comments can also be nested:

```tes
</ outer
    </ inner />
   still outer
/>
```

## Identifiers

Identifiers are names for variables and functions.

Valid:

```tes
x
total
num_list
_temp1
```

Invalid:

```tes
1x
my-name
```

## Keywords

Reserved words:

```text
funk as if else begin endif while endwhile do for to endfor return
int vector str mstr bool boolean null true false
```

`boolean` is accepted as an alias for `bool`.

## Function Syntax

Block form:

```tes
funk <int> add(a as int, b as int)
{
    return a + b;
}
```

Short return form:

```tes
funk <int> one() => return 1;
```

## Variable Declaration

Without initializer:

```tes
x :: int;
```

With initializer:

```tes
x :: int = 10;
```

Variables declared without an initializer must be assigned before being read.

## Assignment

```tes
x = 20;
```

Vector element assignment:

```tes
items[0] = 99;
```

## Return

```tes
return expression;
```

For `main`, this is accepted:

```tes
funk <null> main()
{
    return 0;
}
```

## Blocks

Function bodies use braces:

```tes
{
    statement;
}
```

Control-flow bodies use `begin` and matching end keywords:

```tes
if [[ true ]]
begin
    print(1);
endif
```

