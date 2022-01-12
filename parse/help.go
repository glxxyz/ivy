// Code generated by "go generate robpike.io/ivy/parse"; DO NOT EDIT.

package parse

var helpLines = []string{
	"Ivy is an interpreter for an APL-like language. It is a plaything and a work in",
	"progress.",
	"",
	"Unlike APL, the input is ASCII and the results are exact (but see the next paragraph).",
	"It uses exact rational arithmetic so it can handle arbitrary precision. Values to be",
	"input may be integers (3, -1), rationals (1/3, -45/67), complex numbers (0j1, 1/2j3) or",
	"floating point values (1e3, -1.5 (representing 1000 and -3/2)).",
	"",
	"Some functions such as sqrt are irrational. When ivy evaluates an irrational",
	"function, the result is stored in a high-precision floating-point number (default",
	"256 bits of mantissa). Thus when using irrational functions, the values have high",
	"precision but are not exact.",
	"",
	"Unlike in most other languages, operators always have the same precedence and",
	"expressions are evaluated in right-associative order. That is, unary operators",
	"apply to everything to the right, and binary operators apply to the operand",
	"immediately to the left and to everything to the right.  Thus, 3*4+5 is 27 (it",
	"groups as 3*(4+5)) and iota 3+2 is 1 2 3 4 5 while 3+iota 2 is 4 5. A vector",
	"is a single operand, so 1 2 3 + 3 + 3 4 5 is (1 2 3) + 3 + (3 4 5), or 7 9 11.",
	"",
	"As a special but important case, note that 1/3, with no intervening spaces, is a",
	"single rational number, not the expression 1 divided by 3. This can affect precedence:",
	"3/6*4 is 2 while 3 / 6*4 is 1/8 since the spacing turns the / into a division",
	"operator. Use parentheses or spaces to disambiguate: 3/(6*4) or 3 /6*4.",
	"",
	"Indexing uses [] notation: x[1], x[1; 2], and so on. Indexing by a vector",
	"selects multiple elements: x[1 2] creates a new item from x[1] and x[2].",
	"",
	"Only a subset of APL's functionality is implemented, but the intention is to",
	"have most numerical operations supported eventually.",
	"",
	"Semicolons separate multiple statements on a line. Variables are alphanumeric and are",
	"assigned with the = operator. Assignment is an expression.",
	"",
	"After each successful expression evaluation, the result is stored in the variable",
	"called _ (underscore) so it can be used in the next expression.",
	"",
	"The APL operators, adapted from https://en.wikipedia.org/wiki/APL_syntax_and_symbols,",
	"and their correspondence are listed here. The correspondence is incomplete and inexact.",
	"",
	"Unary operators",
	"",
	"\tName              APL   Ivy     Meaning",
	"\tRoll              ?B    ?       One integer selected randomly from the first B integers",
	"\tCeiling           ⌈B    ceil    Least integer greater than or equal to B",
	"\tFloor             ⌊B    floor   Greatest integer less than or equal to B",
	"\tShape             ⍴B    rho     Number of components in each dimension of B",
	"\tNot               ∼B    not     Logical: not 1 is 0, not 0 is 1",
	"\tAbsolute value    ∣B    abs     Magnitude of B",
	"\tIndex generator   ⍳B    iota    Vector of the first B integers",
	"\tExponential       ⋆B    **      e to the B power",
	"\tNegation          −B    -       Changes sign of B",
	"\tIdentity          +B    +       No change to B",
	"\tSignum            ×B    sgn     ¯1 if B<0; 0 if B=0; 1 if B>0",
	"\tReciprocal        ÷B    /       1 divided by B",
	"\tRavel             ,B    ,       Reshapes B into a vector",
	"\tMatrix inverse    ⌹B            Inverse of matrix B",
	"\tPi times          ○B            Multiply by π",
	"\tLogarithm         ⍟B    log     Natural logarithm of B",
	"\tReversal          ⌽B    rot     Reverse elements of B along last axis",
	"\tReversal          ⊖B    flip    Reverse elements of B along first axis",
	"\tGrade up          ⍋B    up      Indices of B which will arrange B in ascending order",
	"\tGrade down        ⍒B    down    Indices of B which will arrange B in descending order",
	"\tExecute           ⍎B    ivy     Execute an APL (ivy) expression",
	"\tMonadic format    ⍕B    text    A character representation of B",
	"\tMonadic transpose ⍉B    transp  Reverse the axes of B",
	"\tFactorial         !B    !       Product of integers 1 to B",
	"\tBitwise not             ^       Bitwise complement of B (integer only)",
	"\tSquare root       B⋆.5  sqrt    Square root of B.",
	"\tSine              1○B   sin     sin(B)",
	"\tCosine            2○B   cos     cos(B)",
	"\tTangent           3○B   tan     tan(B)",
	"\tInverse sine      ¯1○B  asin    arcsin(B)",
	"\tInverse cosine    ¯2○B  acos    arccos(B)",
	"\tInverse tangent   ¯3○B  atan    arctan(B)",
	"\tReal part         9○B   real    Real part of a complex number.",
	"\tImaginary part    11○B  imag    Imaginary part of a complex number.",
	"\tPhase angle       12○B  phase   Phase angle (argument) of a complex number.",
	"",
	"Binary operators",
	"",
	"\tName                  APL   Ivy     Meaning",
	"\tAdd                   A+B   +       Sum of A and B",
	"\tSubtract              A−B   -       A minus B",
	"\tMultiply              A×B   *       A multiplied by B",
	"\tDivide                A÷B   /       A divided by B (exact rational division)",
	"\t                            div     A divided by B (Euclidean)",
	"\t                            idiv    A divided by B (Go)",
	"\tExponentiation        A⋆B   **      A raised to the B power",
	"\tDeal                  A?B   ?       A distinct integers selected randomly from the first B integers",
	"\tMembership            A∈B   in      1 for elements of A present in B; 0 where not.",
	"\tMaximum               A⌈B   max     The greater value of A or B",
	"\tMinimum               A⌊B   min     The smaller value of A or B",
	"\tReshape               A⍴B   rho     Array of shape A with data B",
	"\tTake                  A↑B   take    Select the first (or last) A elements of B according to ×A",
	"\tDrop                  A↓B   drop    Remove the first (or last) A elements of B according to ×A",
	"\tDecode                A⊥B   decode  Value of a polynomial whose coefficients are B at A",
	"\tEncode                A⊤B   encode  Base-A representation of the value of B",
	"\tResidue               A∣B           B modulo A",
	"\t                            mod     A modulo B (Euclidean)",
	"\t                            imod    A modulo B (Go)",
	"\tCatenation            A,B   ,       Elements of B appended to the elements of A",
	"\tExpansion             A\\B   fill    Insert zeros (or blanks) in B corresponding to zeros in A",
	"\t                                    In ivy: abs(A) gives count, A <= 0 inserts zero (or blank)",
	"\tCompression           A/B   sel     Select elements in B corresponding to ones in A",
	"\t                                    In ivy: abs(A) gives count, A <= 0 inserts zero",
	"\tIndex of              A⍳B   iota    The location (index) of B in A; 1+⌈/⍳⍴A if not found",
	"\t                                    In ivy: origin-1 if not found (i.e. 0 if one-indexed)",
	"\tMatrix divide         A⌹B           Solution to system of linear equations Ax = B",
	"\tRotation              A⌽B   rot     The elements of B are rotated A positions left",
	"\tRotation              A⊖B   flip    The elements of B are rotated A positions along the first axis",
	"\tLogarithm             A⍟B   log     Logarithm of B to base A",
	"\tDyadic format         A⍕B   text    Format B into a character matrix according to A",
	"\t                                    A is the textual format (see format special command);",
	"\t                                    otherwise result depends on length of A:",
	"\t                                    1 gives decimal count, 2 gives width and decimal count,",
	"\t                                    3 gives width, decimal count, and style ('d', 'e', 'f', etc.).",
	"\tGeneral transpose     A⍉B   transp  The axes of B are ordered by A",
	"\tCombinations          A!B   !       Number of combinations of B taken A at a time",
	"\tLess than             A<B   <       Comparison: 1 if true, 0 if false",
	"\tLess than or equal    A≤B   <=      Comparison: 1 if true, 0 if false",
	"\tEqual                 A=B   ==      Comparison: 1 if true, 0 if false",
	"\tGreater than or equal A≥B   >=      Comparison: 1 if true, 0 if false",
	"\tGreater than          A>B   >       Comparison: 1 if true, 0 if false",
	"\tNot equal             A≠B   !=      Comparison: 1 if true, 0 if false",
	"\tOr                    A∨B   or      Logic: 0 if A and B are 0; 1 otherwise",
	"\tAnd                   A∧B   and     Logic: 1 if A and B are 1; 0 otherwise",
	"\tNor                   A⍱B   nor     Logic: 1 if both A and B are 0; otherwise 0",
	"\tNand                  A⍲B   nand    Logic: 0 if both A and B are 1; otherwise 1",
	"\tXor                         xor     Logic: 1 if A != B; otherwise 0",
	"\tBitwise and                 &       Bitwise A and B (integer only)",
	"\tBitwise or                  |       Bitwise A or B (integer only)",
	"\tBitwise xor                 ^       Bitwise A exclusive or B (integer only)",
	"\tLeft shift                  <<      A shifted left B bits (integer only)",
	"\tRight Shift                 >>      A shifted right B bits (integer only)",
	"",
	"Operators and axis indicator",
	"",
	"\tName                APL  Ivy  APL Example  Ivy Example  Meaning (of example)",
	"\tReduce (last axis)  /    /    +/B          +/B          Sum across B",
	"\tReduce (first axis) ⌿         +⌿B                       Sum down B",
	"\tScan (last axis)    \\    \\    +\\B          +\\B          Running sum across B",
	"\tScan (first axis)   ⍀         +⍀B                       Running sum down B",
	"\tInner product       .    .    A+.×B        A +.* B      Matrix product of A and B",
	"\tOuter product       ∘.   o.   A∘.×B        A o.* B      Outer product of A and B",
	"\t                                                    (lower case o; may need preceding space)",
	"\tComplex number      J    j    AJB          AjB          A and B are the real and imaginary parts",
	"",
	"Type-converting operations",
	"",
	"\tName              APL   Ivy     Meaning",
	"\tCode                    code B  The integer Unicode value of char B",
	"\tChar                    char B  The character with integer Unicode value B",
	"\tFloat                   float B The floating-point representation of B",
	"",
	"Pre-defined constants",
	"",
	"The constants e (base of natural logarithms) and pi (π) are pre-defined to high",
	"precision, about 3000 decimal digits truncated according to the floating point",
	"precision setting.",
	"",
	"Character data",
	"",
	"Strings are vectors of \"chars\", which are Unicode code points (not bytes).",
	"Syntactically, string literals are very similar to those in Go, with back-quoted",
	"raw strings and double-quoted interpreted strings. Unlike Go, single-quoted strings",
	"are equivalent to double-quoted, a nod to APL syntax. A string with a single char",
	"is just a singleton char value; all others are vectors. Thus ``, \"\", and '' are",
	"empty vectors, `a`, \"a\", and 'a' are equivalent representations of a single char,",
	"and `ab`, `a` `b`, \"ab\", \"a\" \"b\", 'ab', and 'a' 'b' are equivalent representations",
	"of a two-char vector.",
	"",
	"Unlike in Go, a string in ivy comprises code points, not bytes; as such it can",
	"contain only valid Unicode values. Thus in ivy \"\\x80\" is illegal, although it is",
	"a legal one-byte string in Go.",
	"",
	"Strings can be printed. If a vector contains only chars, it is printed without",
	"spaces between them.",
	"",
	"Chars have restricted operations. Printing, comparison, indexing and so on are",
	"legal but arithmetic is not, and chars cannot be converted automatically into other",
	"singleton values (ints, floats, and so on). The unary operators char and code",
	"enable transcoding between integer and char values.",
	"",
	"User-defined operators",
	"",
	"Users can define unary and binary operators, which then behave just like",
	"built-in operators. Both a unary and a binary operator may be defined for the",
	"same name.",
	"",
	"The syntax of a definition is the 'op' keyword, the operator and formal",
	"arguments, an equals sign, and then the body. The names of the operator and its",
	"arguments must be identifiers.  For unary operators, write \"op name arg\"; for",
	"binary write \"op leftarg name rightarg\". The final expression in the body is the",
	"return value. Operators may have recursive definitions; see the paragraph",
	"about conditional execution for an example.",
	"",
	"The body may be a single line (possibly containing semicolons) on the same line",
	"as the 'op', or it can be multiple lines. For a multiline entry, there is a",
	"newline after the '=' and the definition ends at the first blank line (ignoring",
	"spaces).",
	"",
	"Conditional execution is done with the \":\" binary conditional return operator,",
	"which is valid only within the code for a user-defined operator. The left",
	"operand must be a scalar. If it is non-zero, the right operand is returned as",
	"the value of the function. Otherwise, execution continues normally. The \":\"",
	"operator has a lower precedence than any other operator; in effect it breaks",
	"the line into two separate expressions.",
	"",
	"Example: average of a vector (unary):",
	"\top avg x = (+/x)/rho x",
	"\tavg iota 11",
	"\tresult: 6",
	"",
	"Example: n largest entries in a vector (binary):",
	"\top n largest x = n take x[down x]",
	"\t3 largest 7 1 3 24 1 5 12 5 51",
	"\tresult: 51 24 12",
	"",
	"Example: multiline operator definition (binary):",
	"\top a sum b =",
	"\t\ta = a+b",
	"\t\ta",
	"",
	"\tiota 3 sum 4",
	"\tresult: 1 2 3 4 5 6 7",
	"",
	"Example: primes less than N (unary):",
	"\top primes N = (not T in T o.* T) sel T = 1 drop iota N",
	"\tprimes 50",
	"\tresult: 2 3 5 7 11 13 17 19 23 29 31 37 41 43 47",
	"",
	"Example: greatest common divisor (binary):",
	"\top a gcd b =",
	"\t\ta == b: a",
	"\t\ta > b: b gcd a-b",
	"\t\ta gcd b-a",
	"",
	"\t1562 gcd !11",
	"\tresult: 22",
	"",
	"On mobile platforms only, due to I/O restrictions, user-defined operators",
	"must be presented on a single line. Use semicolons to separate expressions:",
	"",
	"\top a gcd b = a == b: a; a > b: b gcd a-b; a gcd b-a",
	"",
	"To declare an operator but not define it, omit the equals sign and what follows.",
	"\top foo x",
	"\top bar x = foo x",
	"\top foo x = -x",
	"\tbar 3",
	"\tresult: -3",
	"\top foo x = /x",
	"\tbar 3",
	"\tresult: 1/3",
	"",
	"Within a user-defined operator, identifiers are local to the invocation unless",
	"they are undefined in the operator but defined globally, in which case they refer to",
	"the global variable. A mechanism to declare locals may come later.",
	"",
	"Special commands",
	"",
	"Ivy accepts a number of special commands, introduced by a right paren",
	"at the beginning of the line. Most report the current value if a new value",
	"is not specified. For these commands, numbers are always read and printed",
	"base 10 and must be non-negative on input.",
	"",
	"\t) help",
	"\t\tDescribe the special commands. Run )help <topic> to learn more",
	"\t\tabout a topic, )help <op> to learn more about an operator.",
	"\t) base 0",
	"\t\tSet the number base for input and output. The commands ibase and",
	"\t\tobase control setting of the base for input and output alone,",
	"\t\trespectively.  Base 0 allows C-style input: decimal, with 037 being",
	"\t\toctal and 0x10 being hexadecimal. If the base is greater than 10,",
	"\t\tany identifier formed from valid numerals in the base system, such",
	"\t\tas abe for base 16, is taken to be a number. TODO: To output",
	"\t\tlarge integers and rationals, base must be one of 0 2 8 10 16.",
	"\t\tFloats are always printed base 10.",
	"\t) cpu",
	"\t\tPrint the duration of the last interactive calculation.",
	"\t) debug name 0|1",
	"\t\tToggle or set the named debugging flag. With no argument, lists",
	"\t\tthe settings.",
	"\t) demo",
	"\t\tRun a line-by-line interactive demo. On mobile platforms,",
	"\t\tuse the Demo menu option instead.",
	"\t) format \"\"",
	"\t\tSet the format for printing values. If empty, the output is printed",
	"\t\tusing the output base. If non-empty, the format determines the",
	"\t\tbase used in printing. The format is in the style of golang.org/pkg/fmt.",
	"\t\tFor floating-point formats, flags and width are ignored.",
	"\t) get \"save.ivy\"",
	"\t\tRead input from the named file; return to interactive execution",
	"\t\tafterwards. If no file is specified, read from \"save.ivy\".",
	"\t\t(Unimplemented on mobile.)",
	"\t) maxbits 1e6",
	"\t\tTo avoid consuming too much memory, if an integer result would",
	"\t\trequire more than this many bits to store, abort the calculation.",
	"\t\tIf maxbits is 0, there is no limit; the default is 1e6.",
	"\t) maxdigits 1e4",
	"\t\tTo avoid overwhelming amounts of output, if an integer has more",
	"\t\tthan this many digits, print it using the defined floating-point",
	"\t\tformat. If maxdigits is 0, integers are always printed as integers.",
	"\t) maxstack 1e5",
	"\t\tTo avoid using too much stack, the number of nested active calls to",
	"\t\tuser-defined operators is limited to maxstack.",
	"\t) op X",
	"\t\tIf X is absent, list all user-defined operators. Otherwise,",
	"\t\tshow the definition of the user-defined operator X. Inside the",
	"\t\tdefinition, numbers are always shown base 10, ignoring the ibase",
	"\t\tand obase.",
	"\t) origin 1",
	"\t\tSet the origin for indexing a vector or matrix.",
	"\t) prec 256",
	"\t\tSet the precision (mantissa length) for floating-point values.",
	"\t\tThe value is in bits. The exponent always has 32 bits.",
	"\t) prompt \"\"",
	"\t\tSet the interactive prompt.",
	"\t) save \"save.ivy\"",
	"\t\tWrite definitions of user-defined operators and variables to the",
	"\t\tnamed file, as ivy textual source. If no file is specified, save to",
	"\t\t\"save.ivy\".",
	"\t\t(Unimplemented on mobile.)",
	"\t) seed 0",
	"\t\tSet the seed for the ? operator.",
}

type helpIndexPair struct {
	start, end int
}

var helpUnary = map[string]helpIndexPair{
	"?":      {43, 43},
	"ceil":   {44, 44},
	"floor":  {45, 45},
	"rho":    {46, 46},
	"not":    {47, 47},
	"abs":    {48, 48},
	"iota":   {49, 49},
	"**":     {50, 50},
	"-":      {51, 51},
	"+":      {52, 52},
	"sgn":    {53, 53},
	"/":      {54, 54},
	",":      {55, 55},
	"log":    {58, 58},
	"rot":    {59, 59},
	"flip":   {60, 60},
	"up":     {61, 61},
	"down":   {62, 62},
	"ivy":    {63, 63},
	"text":   {64, 64},
	"transp": {65, 65},
	"!":      {66, 66},
	"^":      {67, 67},
	"sqrt":   {68, 68},
	"sin":    {69, 71},
	"cos":    {69, 71},
	"tan":    {69, 71},
	"asin":   {72, 74},
	"acos":   {72, 74},
	"atan":   {72, 74},
	"real":   {75, 77},
	"imag":   {75, 77},
	"phase":  {75, 77},
	"code":   {151, 151},
	"char":   {152, 152},
	"float":  {153, 153},
}

var helpBinary = map[string]helpIndexPair{
	"+":      {82, 82},
	"-":      {83, 83},
	"*":      {84, 84},
	"/":      {85, 87},
	"**":     {88, 88},
	"?":      {89, 89},
	"in":     {90, 90},
	"max":    {91, 91},
	"min":    {92, 92},
	"rho":    {93, 93},
	"take":   {94, 94},
	"drop":   {95, 95},
	"decode": {96, 96},
	"encode": {97, 97},
	"mod":    {99, 100},
	",":      {101, 101},
	"fill":   {102, 103},
	"sel":    {104, 105},
	"iota":   {106, 107},
	"rot":    {109, 109},
	"flip":   {110, 110},
	"log":    {111, 111},
	"text":   {112, 116},
	"transp": {117, 117},
	"!":      {118, 118},
	"<":      {119, 119},
	"<=":     {120, 120},
	"==":     {121, 121},
	">=":     {122, 122},
	">":      {123, 123},
	"!=":     {124, 124},
	"or":     {125, 125},
	"and":    {126, 126},
	"nor":    {127, 127},
	"nand":   {128, 128},
	"xor":    {129, 129},
	"&":      {130, 130},
	"|":      {131, 131},
	"^":      {132, 132},
	"<<":     {133, 133},
	">>":     {134, 134},
}

var helpAxis = map[string]helpIndexPair{
	"/":  {139, 139},
	"\\": {141, 141},
	".":  {143, 143},
	"o.": {144, 144},
	"j":  {146, 146},
}
