// Code generated by goyacc -o parser.go parser.y. DO NOT EDIT.

//line parser.y:2
package parser

import __yyfmt__ "fmt"

//line parser.y:2

// import "fmt"

//line parser.y:8
type yySymType struct {
	yys int
	b   bool
	i   int
	f   float64
	s   string
	a   []string
}

const IDENT = 57346
const INT = 57347
const FLOAT = 57348
const STRING = 57349
const TRUE = 57350
const FALSE = 57351
const ADD = 57352
const SUB = 57353
const MUL = 57354
const DIV = 57355
const MOD = 57356
const ASSIGN = 57357
const ADD_ASSIGN = 57358
const SUB_ASSIGN = 57359
const MUL_ASSIGN = 57360
const DIV_ASSIGN = 57361
const MOD_ASSIGN = 57362
const AND = 57363
const OR = 57364
const INC = 57365
const DEC = 57366
const EQ = 57367
const NOT_EQ = 57368
const NOT = 57369
const LT = 57370
const GT = 57371
const LTE = 57372
const GTE = 57373
const ELLIPSIS = 57374
const DOT = 57375
const COMMA = 57376
const COLON = 57377
const LPAREN = 57378
const RPAREN = 57379
const LBRACE = 57380
const RBRACE = 57381
const LBRAKET = 57382
const RBRAKET = 57383
const CONTRACT = 57384
const DATA = 57385
const CONDITION = 57386
const ACTION = 57387
const FUNC = 57388
const VAR = 57389
const EXTEND_VAR = 57390
const IF = 57391
const ELSE = 57392
const WHILE = 57393
const BREAK = 57394
const CONTINUE = 57395
const INFO = 57396
const WARNING = 57397
const ERROR = 57398
const NIL = 57399
const RETURN = 57400
const T_BOOL = 57401
const T_INT = 57402
const T_FLOAT = 57403
const T_MONEY = 57404
const T_STRING = 57405
const T_BYTES = 57406
const T_ARRAY = 57407
const T_MAP = 57408
const T_FILE = 57409

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"IDENT",
	"INT",
	"FLOAT",
	"STRING",
	"TRUE",
	"FALSE",
	"ADD",
	"SUB",
	"MUL",
	"DIV",
	"MOD",
	"ASSIGN",
	"ADD_ASSIGN",
	"SUB_ASSIGN",
	"MUL_ASSIGN",
	"DIV_ASSIGN",
	"MOD_ASSIGN",
	"AND",
	"OR",
	"INC",
	"DEC",
	"EQ",
	"NOT_EQ",
	"NOT",
	"LT",
	"GT",
	"LTE",
	"GTE",
	"ELLIPSIS",
	"DOT",
	"COMMA",
	"COLON",
	"LPAREN",
	"RPAREN",
	"LBRACE",
	"RBRACE",
	"LBRAKET",
	"RBRAKET",
	"CONTRACT",
	"DATA",
	"CONDITION",
	"ACTION",
	"FUNC",
	"VAR",
	"EXTEND_VAR",
	"IF",
	"ELSE",
	"WHILE",
	"BREAK",
	"CONTINUE",
	"INFO",
	"WARNING",
	"ERROR",
	"NIL",
	"RETURN",
	"T_BOOL",
	"T_INT",
	"T_FLOAT",
	"T_MONEY",
	"T_STRING",
	"T_BYTES",
	"T_ARRAY",
	"T_MAP",
	"T_FILE",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 231

var yyAct = [...]int{

	118, 53, 113, 16, 66, 72, 61, 56, 75, 6,
	42, 51, 85, 5, 15, 144, 55, 28, 143, 142,
	71, 17, 18, 19, 20, 21, 22, 23, 24, 25,
	79, 13, 76, 17, 18, 19, 20, 21, 22, 23,
	24, 25, 55, 78, 114, 114, 12, 37, 38, 39,
	5, 31, 83, 55, 85, 37, 38, 39, 5, 93,
	94, 95, 92, 29, 86, 47, 91, 85, 84, 3,
	131, 55, 73, 87, 138, 74, 14, 96, 97, 140,
	112, 90, 102, 36, 32, 45, 55, 116, 115, 55,
	119, 108, 55, 55, 124, 117, 26, 120, 89, 145,
	36, 125, 126, 121, 80, 81, 134, 135, 136, 132,
	133, 45, 137, 110, 139, 77, 103, 104, 141, 127,
	128, 129, 130, 109, 67, 68, 29, 69, 4, 9,
	63, 64, 98, 99, 100, 101, 105, 106, 107, 8,
	2, 30, 7, 58, 59, 55, 146, 65, 35, 34,
	33, 111, 41, 49, 67, 68, 70, 69, 15, 82,
	63, 64, 67, 68, 48, 69, 5, 50, 63, 64,
	44, 43, 10, 58, 59, 11, 27, 65, 52, 46,
	88, 58, 59, 54, 60, 65, 70, 123, 15, 40,
	57, 62, 1, 0, 70, 122, 5, 50, 67, 68,
	0, 69, 0, 0, 63, 64, 0, 0, 52, 0,
	0, 0, 0, 0, 0, 0, 0, 58, 59, 0,
	0, 65, 0, 0, 0, 0, 0, 0, 0, 0,
	70,
}
var yyPact = [...]int{

	-33, -33, -1000, -1000, -1000, 135, 125, -1000, 10, -7,
	-24, -38, 59, 12, -1000, 150, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 38, -26, -1000,
	4, -1000, -1000, -1000, -1000, -1000, -1000, -8, -24, -24,
	-1000, 120, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	122, 20, 194, -1000, 51, 83, 60, 26, 194, 194,
	194, 52, -1000, -1000, -1000, -1000, 104, -1000, -1000, -1000,
	194, 106, 124, 119, -1000, 109, -1000, -1000, -1000, 41,
	-1000, -1000, -1000, -1000, -26, 194, 20, 194, 194, -1000,
	194, 194, 158, -1000, -1000, -1000, 194, 194, 194, 194,
	194, 194, 33, 194, 194, 194, 194, 194, -26, 42,
	-1000, 40, -1000, -1000, -38, -1000, -1000, 60, -1000, -1000,
	52, -22, -1000, -19, -1000, 104, 104, 106, 106, 106,
	106, -1000, 124, 124, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, 92, -1000, -1000, 194, -1000, -1000,
}
var yyPgo = [...]int{

	0, 192, 191, 11, 190, 187, 1, 0, 184, 5,
	20, 4, 6, 7, 183, 180, 3, 17, 179, 176,
	175, 172, 69, 65, 171, 170, 164, 153, 152, 10,
	2, 151, 150, 149, 148, 84, 141, 128, 140,
}
var yyR1 = [...]int{

	0, 2, 2, 2, 2, 4, 4, 4, 4, 5,
	5, 7, 7, 7, 7, 8, 8, 8, 9, 9,
	9, 9, 10, 10, 10, 11, 11, 11, 11, 11,
	12, 12, 12, 13, 13, 14, 14, 6, 6, 15,
	3, 3, 16, 16, 16, 16, 16, 16, 16, 16,
	16, 17, 17, 18, 19, 19, 19, 20, 20, 21,
	21, 22, 24, 24, 25, 25, 25, 23, 23, 28,
	28, 29, 29, 26, 26, 27, 27, 30, 30, 31,
	31, 32, 32, 33, 34, 35, 35, 35, 35, 36,
	36, 37, 37, 38, 38, 1, 1,
}
var yyR2 = [...]int{

	0, 1, 1, 1, 3, 1, 4, 3, 4, 1,
	3, 1, 2, 2, 2, 1, 1, 1, 1, 3,
	3, 3, 1, 3, 3, 1, 3, 3, 3, 3,
	1, 3, 3, 1, 3, 1, 3, 1, 3, 1,
	1, 3, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 3, 3, 2, 4, 4, 2, 3, 2,
	1, 4, 1, 1, 1, 1, 1, 2, 3, 1,
	2, 1, 1, 0, 1, 1, 2, 3, 2, 1,
	2, 4, 3, 2, 2, 1, 1, 1, 1, 1,
	2, 5, 4, 1, 1, 1, 2,
}
var yyChk = [...]int{

	-1000, -1, -38, -22, -37, 46, 42, -38, 4, 4,
	-21, -20, 36, 38, -23, 38, -16, 59, 60, 61,
	62, 63, 64, 65, 66, 67, 37, -19, -17, 4,
	-36, 39, -35, -32, -33, -34, -22, 43, 44, 45,
	39, -28, -29, -24, -25, -22, -18, -23, -26, -27,
	47, -3, 58, -6, -14, -7, -13, -4, 23, 24,
	-8, -12, -2, 10, 11, 27, -11, 4, 5, 7,
	36, -10, -9, 34, 37, 34, -16, -35, 39, 38,
	-23, -23, 39, -29, -17, 34, -3, 22, -15, 15,
	21, 40, 36, -7, -7, -7, 25, 26, 28, 29,
	30, 31, -3, 10, 11, 12, 13, 14, -17, 4,
	4, -31, 39, -30, 4, -16, -6, -13, -7, -6,
	-12, -3, 37, -5, -6, -11, -11, -10, -10, -10,
	-10, 37, -9, -9, -7, -7, -7, -16, 32, -30,
	39, -16, 41, 37, 34, 7, -6,
}
var yyDef = [...]int{

	0, -2, 95, 93, 94, 0, 0, 96, 0, 0,
	0, 60, 0, 0, 61, 0, 59, 42, 43, 44,
	45, 46, 47, 48, 49, 50, 57, 0, 0, 51,
	0, 92, 89, 85, 86, 87, 88, 0, 0, 0,
	67, 0, 69, 71, 72, 62, 63, 64, 65, 66,
	0, 74, 75, 40, 37, 18, 35, 11, 0, 0,
	0, 33, 5, 15, 16, 17, 30, 1, 2, 3,
	0, 25, 22, 0, 58, 0, 54, 90, 91, 0,
	83, 84, 68, 70, 0, 0, 76, 0, 0, 39,
	0, 0, 0, 12, 13, 14, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 51,
	52, 0, 82, 79, 0, 53, 41, 36, 18, 38,
	34, 0, 7, 0, 9, 31, 32, 26, 27, 28,
	29, 4, 23, 24, 19, 20, 21, 55, 56, 80,
	81, 78, 6, 8, 0, 77, 10,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61,
	62, 63, 64, 65, 66, 67,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	}
	goto yystack /* stack new state and value */
}
