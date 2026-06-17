# Functions and Scope

## Function Declaration

```tes
funk <int> add(a as int, b as int)
{
    return a + b;
}
```

Parts:

```text
funk <return-type> name(parameter-list)
{
    statements
}
```

## Parameters

Parameters use `name as type`:

```tes
funk <int> square(x as int)
{
    return x * x;
}
```

Multiple parameters:

```tes
funk <int> add(a as int, b as int)
{
    return a + b;
}
```

## Calling Functions

```tes
result :: int = add(10, 20);
```

Argument count and argument types must match the function declaration.

## Return Type

The returned expression must match the declared return type:

```tes
funk <int> f()
{
    return 1;
}
```

This is invalid:

```tes
funk <vector> f()
{
    return 1;
}
```

## main Function

`main` is the normal program entry point:

```tes
funk <null> main()
{
    print(1);
    return 0;
}
```

## Nested Functions

Functions can be declared inside other functions:

```tes
funk <int> outer()
{
    funk <int> inner()
    {
        return 10;
    }

    return inner();
}
```

Nested functions are visible inside their enclosing function.

Current code generation uses name mangling for nested functions. Full closure conversion is not implemented, so nested functions should avoid depending on captured outer variables for portable behavior.

## Scope Rules

Variables are resolved from the current scope outward.

```tes
funk <null> main()
{
    x :: int = 1;

    if [[ true ]]
    begin
        y :: int = x + 1;
        print(y);
    endif

    return 0;
}
```

Variables cannot be redeclared in the same scope:

```tes
x :: int;
x :: vector; </ error />
```

Parameters are initialized automatically. Local variables without initializers are not.

