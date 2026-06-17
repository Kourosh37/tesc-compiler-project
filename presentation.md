# ارائه پروژه کامپایلر TesLang

**درس اصول طراحی کامپایلر**

پیاده‌سازی یک کامپایلر برای زبان **TesLang** در سه مرحله:

**تحلیل لغوی ← تحلیل نحوی و معنایی ← تولید کد میانی TESVM**

---

## صورت مسئله

در داکیومنت پروژه، کامپایلر باید سه کار اصلی انجام دهد:

۱. **گام اول:** خواندن کد TesLang، حذف کامنت‌ها، استخراج توکن‌ها، و نمایش خط و ستون هر توکن.

۲. **گام دوم:** بررسی ساختار برنامه با Parser و پیدا کردن خطاهای معنایی با Symbol Table.

۳. **گام سوم:** تبدیل برنامه معتبر TesLang به کد میانی قابل اجرا روی TSVM/TESVM.

**ارجاع کد:**

```text
cmd/teslang/main.go:163-218
```

---

## مسیر کلی کامپایلر

کامپایلر به شکل یک pipeline ساده پیاده‌سازی شده است:

```text
TesLang Source
   ↓
Lexer → Tokenها
   ↓
Parser → AST
   ↓
Semantic Analyzer → خطا یا تأیید
   ↓
Code Generator → TESVM IR
```

در ابزار اصلی، همین ترتیب مرحله‌به‌مرحله اجرا می‌شود.

**ارجاع کد:**

```text
cmd/teslang/main.go:163-198
cmd/tes/main.go:102-115
```

---

## ابزارهای اجرایی پروژه

پروژه سه ابزار اصلی دارد:

| ابزار | کاربرد |
|---|---|
| `teslang` | کامپایلر اصلی؛ چاپ توکن، بررسی semantic، تولید TESVM |
| `tes` | کامپایل و اجرای مستقیم برنامه |
| `tesvm` | اجرای فایل TESVM |

در ارائه، تمرکز اصلی روی `teslang` است؛ چون دقیقاً مراحل خواسته‌شده PDF را اجرا می‌کند.

**ارجاع کد:**

```text
cmd/teslang/main.go:36-43
cmd/tes/main.go:102-132
cmd/tesvm/main.go
```

---

## گام اول: تحلیلگر لغوی

طبق PDF، Lexer باید از بین کدهای TesLang فقط توکن‌های واقعی را بیرون بکشد.

کارهای انجام‌شده در پروژه:

- خواندن کل ورودی
- حذف فاصله‌ها و کامنت‌ها
- تولید توکن‌ها
- ذخیره خط و ستون شروع هر توکن
- گزارش خطاهای لغوی مثل کاراکتر ناشناخته یا رشته ناتمام

**ارجاع کد:**

```text
internal/lexer/lexer.go:11-31
internal/lexer/lexer.go:33-69
cmd/teslang/main.go:220-227
```

---

## خط و ستون توکن‌ها

PDF خواسته بود هر توکن همراه با **Line** و **Column** چاپ شود.

در این پروژه:

- هر `Token` فیلدهای `Line` و `Column` دارد.
- Lexer قبل از خواندن هر توکن، خط و ستون شروع آن را ذخیره می‌کند.
- با هر `\n` شماره خط زیاد می‌شود و ستون دوباره از ۱ شروع می‌شود.

**ارجاع کد:**

```text
internal/token/token.go:73-80
internal/lexer/lexer.go:42-45
internal/lexer/lexer.go:197-205
cmd/teslang/main.go:220-227
```

---

## حذف کامنت‌های تو در تو

طبق PDF، کامنت‌ها با `</` شروع و با `/>` تمام می‌شوند و می‌توانند تو در تو باشند.

راه‌حل پروژه خیلی ساده است:

```text
وقتی </ دیده شد: depth = 1
هر </ جدید: depth++
هر /> جدید: depth--
تا وقتی depth صفر نشده، هیچ توکنی تولید نمی‌شود.
```

**ارجاع کد:**

```text
internal/lexer/lexer.go:36-38
internal/lexer/lexer.go:82-100
examples/compiler/lexer_sample.tes:24-33
```

---

## توکن‌ها، کلمات کلیدی و رشته‌ها

Lexer موارد اصلی زبان TesLang را تشخیص می‌دهد:

- شناسه‌ها و keywordها مثل `funk`, `for`, `return`
- عددها
- رشته‌های معمولی با `'` یا `"`
- رشته‌های چندخطی با `"""`
- عملگرها مثل `+`, `==`, `&&`, `::`, `=>`

**ارجاع کد:**

```text
internal/token/token.go:10-71
internal/token/token.go:86-91
internal/lexer/lexer.go:47-63
internal/lexer/lexer.go:102-155
internal/lexer/lexer.go:157-176
```

---

## گام دوم: Parser

PDF گفته بود تجزیه نحوی الزامی است.

در پروژه، Parser توکن‌ها را می‌گیرد و یک درخت AST می‌سازد. بعداً Semantic Analyzer و Codegen همین AST را بررسی می‌کنند.

**ارجاع کد:**

```text
internal/parser/parser.go:11-20
internal/parser/parser.go:20-31
internal/ast/ast.go:5-17
internal/ast/ast.go:35-58
```

---

## Parser: تابع‌ها و بدنه برنامه

طبق گرامر PDF، برنامه از چند تابع تشکیل می‌شود.

Parser این موارد را پشتیبانی می‌کند:

- تعریف تابع با `funk <type> name(...)`
- پارامترها با قالب `name as type`
- تابع با بدنه `{ ... }`
- تابع کوتاه با `=> return expr`

**ارجاع کد:**

```text
internal/parser/parser.go:20-31
internal/parser/parser.go:33-62
internal/parser/parser.go:64-72
```

---

## Parser: دستورهای زبان

دستورهای اصلی TesLang طبق PDF در Parser پیاده‌سازی شده‌اند:

| دستور | ارجاع کد |
|---|---|
| تعریف متغیر | `internal/parser/parser.go:97-107` |
| `return` | `internal/parser/parser.go:108-113` |
| `if / else` | `internal/parser/parser.go:114-127` |
| `while` | `internal/parser/parser.go:128-137` |
| `do while` | `internal/parser/parser.go:138-148` |
| `for` | `internal/parser/parser.go:149-162` |
| expression statement | `internal/parser/parser.go:90-94` |

---

## Parser: مدیریت اولویت عملگرها

در PDF گفته شده گرامر expression مبهم است و باید اولویت عملگرها مشخص شود.

در پروژه از **Pratt Parser** استفاده شده است. یعنی Parser برای هر عملگر یک precedence دارد و expressionها را با ترتیب درست می‌سازد.

مثلاً:

```text
a + b * c
```

به صورت `a + (b * c)` تحلیل می‌شود، چون `*` اولویت بیشتری دارد.

**ارجاع کد:**

```text
internal/parser/parser.go:192-203
internal/parser/parser.go:205-240
internal/parser/parser.go:300-316
```

---

## گام دوم: تحلیل معنایی

بعد از Parser، برنامه از نظر معنایی بررسی می‌شود.

این مرحله خطاهایی را پیدا می‌کند که ممکن است از نظر syntax درست باشند، ولی از نظر منطق زبان غلط هستند؛ مثل:

- استفاده از متغیر تعریف‌نشده
- استفاده از متغیر مقداردهی‌نشده
- نوع اشتباه در assignment
- تعداد یا نوع اشتباه آرگومان‌های تابع
- return با نوع اشتباه

**ارجاع کد:**

```text
internal/semantic/analyzer.go:18-27
internal/semantic/analyzer.go:72-136
internal/semantic/analyzer.go:142-255
```

---

## Symbol Table و Scope

PDF خواسته بود برای Semantic Analyzer از Symbol Table با scopeهای تو در تو استفاده شود.

در پروژه هر scope شامل این موارد است:

```text
Parent    → scope بالادستی
Variables → متغیرهای همان scope
Functions → تابع‌های همان scope
```

اگر نامی در scope فعلی پیدا نشود، در scopeهای پدر جست‌وجو می‌شود.

**ارجاع کد:**

```text
internal/semantic/symbol_table.go:5-13
internal/semantic/symbol_table.go:14-28
internal/semantic/symbol_table.go:29-43
internal/semantic/symbols.go:5-22
```

---

## ثبت تابع‌ها و Built-inها

برای اینکه تابع‌ها قبل از استفاده شناخته شوند، Semantic Analyzer دو مرحله دارد:

۱. اول همه تابع‌ها را در Symbol Table ثبت می‌کند.

۲. بعد بدنه تابع‌ها را تحلیل می‌کند.

همچنین تابع‌های داخلی مثل `scan`, `print`, `list`, `length`, `exit` ثبت شده‌اند.

**ارجاع کد:**

```text
internal/semantic/analyzer.go:18-27
internal/semantic/analyzer.go:29-39
internal/semantic/analyzer.go:40-48
internal/semantic/analyzer.go:49-71
```

---

## خطاهای خواسته‌شده در PDF

| خطای خواسته‌شده | محل پیاده‌سازی |
|---|---|
| متغیر تعریف‌نشده | `internal/semantic/analyzer.go:145-149` |
| متغیر مقداردهی‌نشده | `internal/semantic/analyzer.go:150-152` |
| نوع اشتباه در تعریف یا assignment | `internal/semantic/analyzer.go:76-87`, `207-220` |
| return با نوع اشتباه | `internal/semantic/analyzer.go:89-97` |
| تعداد آرگومان اشتباه | `internal/semantic/analyzer.go:243-245` |
| نوع آرگومان اشتباه | `internal/semantic/analyzer.go:246-250` |
| index روی غیر-vector | `internal/semantic/analyzer.go:196-203` |
| شرط غیر-bool | `internal/semantic/analyzer.go:100-117` |

---

## نمونه خطاهای معنایی

فایل `semantic_errors.tes` برای نمایش خطاهای مرحله دوم ساخته شده است.

چند نمونه از خطاهای داخل آن:

- `k` و `j` تعریف شده‌اند ولی قبل از مقداردهی خوانده می‌شوند.
- `x` اصلاً تعریف نشده است.
- `A` از نوع `int` است، ولی مقدار `list(3)` یعنی `vector` می‌گیرد.
- `find(A)` تعداد آرگومان اشتباه دارد.
- `find(a, A)` نوع آرگومان‌ها را اشتباه می‌فرستد.
- `main` از نوع `null` است ولی `A` را return می‌کند.

**ارجاع کد:**

```text
examples/compiler/semantic_errors.tes:22-65
```

---

## فرمت گزارش خطا

طبق PDF، خطاها باید مکان داشته باشند تا بتوان سریع آن‌ها را پیدا کرد.

در پروژه، هر خطا شامل این موارد است:

```text
مرحله خطا + خط + ستون + نام تابع + پیام خطا
```

نمونه:

```text
Error [semantic] line 30, column 20 in function 'find':
variable 'k' is used before being assigned.
```

**ارجاع کد:**

```text
internal/diagnostic/diagnostic.go:5-19
internal/diagnostic/diagnostic.go:21-30
internal/semantic/analyzer.go:292-295
```

---

## گام سوم: تولید کد میانی TESVM

طبق PDF، بعد از تأیید برنامه باید کد میانی تولید شود.

Generator خروجی متنی TESVM می‌سازد. هر تابع TesLang به یک `proc` تبدیل می‌شود.

نمونه:

```tesvm
proc add
  add r3, r1, r2
  mov r0, r3
  ret
```

**ارجاع کد:**

```text
internal/codegen/tesvm.go:20-25
internal/codegen/tesvm.go:28-57
```

---

## قرارداد registerها در Codegen

برای تولید کد، پروژه از register استفاده می‌کند:

- `r0` برای مقدار برگشتی تابع است.
- پارامترهای تابع از `r1` شروع می‌شوند.
- متغیرهای محلی و مقدارهای موقت register جدا می‌گیرند.

این مدل شبیه نمونه کد میانی داخل PDF است.

**ارجاع کد:**

```text
internal/codegen/registers.go:5-26
internal/codegen/tesvm.go:33-39
internal/codegen/tesvm.go:67-70
```

---

## تولید کد برای دستورها

Generator برای statementهای اصلی کد میانی تولید می‌کند:

| دستور TesLang | خروجی TESVM |
|---|---|
| تعریف متغیر | گرفتن register و تولید مقدار اولیه |
| `return` | `mov r0, value` و `ret` |
| `if / else` | `jz`, `jmp`, `label` |
| `while` | label شروع، شرط، پرش به پایان |
| `do while` | اجرای بدنه و سپس شرط |
| `for` | مقدار اولیه، شرط، افزایش اندیس |

**ارجاع کد:**

```text
internal/codegen/tesvm.go:58-121
```

---

## تولید کد برای expressionها

Expressionها به register، literal، یا دستور TESVM تبدیل می‌شوند.

نمونه‌ها:

- `a + b` → دستور `add`
- `x = expr` → تولید مقدار و ذخیره در register متغیر
- `v[i]` → دستور `loadidx`
- `v[i] = x` → دستور `storeidx`
- `print(x)` → دستور `call log, x`
- `scan()` → دستور `call read, dst`

**ارجاع کد:**

```text
internal/codegen/tesvm.go:123-149
internal/codegen/tesvm.go:150-241
```

---

## نمونه کامل گام سوم

در نمونه زیر، برنامه دو عدد را می‌خواند، جمع می‌کند و چاپ می‌کند:

```teslang
funk <int> add(a as int, b as int) {
  result :: int = a + b;
  return result;
}
```

این بخش به TESVM تبدیل می‌شود و دستور `add` تولید می‌کند.

**ارجاع کد:**

```text
examples/compiler/codegen_sample.tes:11-16
examples/compiler/codegen_sample.tes:19-33
internal/codegen/tesvm.go:123-149
internal/codegen/tesvm.go:184-188
```

---

## اجرای پروژه برای ارائه

برای نمایش هر مرحله می‌توان این دستورها را اجرا کرد:

```powershell
Get-Content .\examples\compiler\lexer_sample.tes |
  go run .\cmd\teslang --tokens
```

```powershell
Get-Content .\examples\compiler\semantic_errors.tes |
  go run .\cmd\teslang --check
```

```powershell
Get-Content .\examples\compiler\codegen_sample.tes |
  go run .\cmd\teslang --emit-tesvm
```

**ارجاع کد:**

```text
README.md:32-34
cmd/teslang/main.go:36-68
```

---

## جمع‌بندی انطباق با PDF

| خواسته PDF | وضعیت در پروژه |
|---|---|
| استخراج توکن‌ها | انجام شده در `internal/lexer` |
| گزارش خط و ستون | انجام شده در `token.Token` و خروجی `--tokens` |
| حذف کامنت تو در تو | انجام شده با `depth` در `skipComment` |
| Parser | انجام شده در `internal/parser` |
| AST | انجام شده در `internal/ast` |
| Symbol Table تو در تو | انجام شده در `internal/semantic` |
| تشخیص خطاهای معنایی | انجام شده در `semantic.Analyzer` |
| تولید کد میانی | انجام شده در `internal/codegen` |

---

## پایان

این پروژه طبق خواسته PDF از ورودی TesLang تا خروجی TESVM پیش می‌رود:

```text
کد TesLang
→ توکن‌ها
→ AST
→ بررسی معنایی
→ کد میانی TESVM
```

نقطه قوت اصلی پروژه این است که هر مرحله جدا، قابل بررسی، و دارای ارجاع دقیق در کد است.

