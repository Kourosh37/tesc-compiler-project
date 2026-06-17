# Examples

## Hello Number

```tes
funk <null> main()
{
    print(123);
    return 0;
}
```

## Add Two Numbers

```tes
funk <int> add(a as int, b as int)
{
    return a + b;
}

funk <null> main()
{
    x :: int = 10;
    y :: int = 20;
    print(add(x, y));
    return 0;
}
```

Output:

```text
30
```

## Read Two Numbers

```tes
funk <int> add(a as int, b as int)
{
    return a + b;
}

funk <null> main()
{
    a :: int;
    b :: int;

    a = scan();
    b = scan();

    print(add(a, b));
    return 0;
}
```

Run:

```powershell
"5`n7`n" | .\bin\tes.exe .\add.tes
```

Output:

```text
12
```

## Sum Vector

```tes
funk <int> sum(items as vector)
{
    result :: int = 0;

    for (i = 0 to length(items))
    begin
        result = result + items[i];
    endfor

    return result;
}

funk <null> main()
{
    values :: vector = list(3);
    values[0] = 10;
    values[1] = 20;
    values[2] = 30;

    print(sum(values));
    return 0;
}
```

Output:

```text
60
```

## if / else

```tes
funk <null> main()
{
    x :: int = 10;

    if [[ x > 5 ]]
    begin
        print(1);
    else
        print(0);
    endif

    return 0;
}
```

## while Loop

```tes
funk <null> main()
{
    i :: int = 0;

    while [[ i < 3 ]]
    begin
        print(i);
        i = i + 1;
    endwhile

    return 0;
}
```

Output:

```text
0
1
2
```

## Ternary

```tes
funk <null> main()
{
    x :: int = 10;
    result :: int = x > 5 ? 100 : 200;
    print(result);
    return 0;
}
```

