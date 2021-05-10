/**
 * This file is part of the go-xfmt package (https://github.com/Illirgway/go-xfmt)
 *
 * Copyright (c) 2021 Illirgway
 *
 * This program is free software: you can redistribute it and/or modify it under the terms of the GNU
 * General Public License as published by the Free Software Foundation, either version 3 of the License,
 * or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY;
 * without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 * See the GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License along with this program.
 * If not, see <https://www.gnu.org/licenses/>.
 *
 */

package xfmt

import (
	"fmt"
	"testing"
	"unicode"
)

// emulated sprintf
func esprintf(format string, args ...string) string {
	xfmt := parseFormat(format)
	return xfmt.Sprint(args)
}

type parseFormatTestCase struct {
	format string
	xfmt   xfmt
}

const (
	errParseFormatMismatchArgsPtn     = "args count mismatch: want %d, got %d"
	errParseFormatMismatchMinSize     = "minSize mismatch: want %d, got %d"
	errParseFormatMismatchTokensCount = "tokens count mismatch: want %d, got %d"
	errParseFormatMismatchToken       = "mismatch %dth token: want %#v, got %#v"
)

// inlined
//go:nosplit
func tokenCompare(t1, t2 *token) bool {
	// ATN! compare by value, not ptr
	return *t1 == *t2
}

func (c *parseFormatTestCase) run(t *testing.T) error {

	xfmt := parseFormat(c.format)

	t.Logf("%q ==> %#v", c.format, xfmt)

	if want, got := c.xfmt.args, xfmt.args; want != got {
		return fmt.Errorf(errParseFormatMismatchArgsPtn, want, got)
	}

	if want, got := c.xfmt.minSize, xfmt.minSize; want != got {
		return fmt.Errorf(errParseFormatMismatchMinSize, want, got)
	}

	if want, got := len(c.xfmt.tokens), len(xfmt.tokens); want != got {
		return fmt.Errorf(errParseFormatMismatchTokensCount, want, got)
	}

	sz := uint(len(c.xfmt.tokens))

	if sz == 0 {
		return nil
	}

	_ = xfmt.tokens[sz-1]

	for i := uint(0); i < sz; i++ {
		tw, tg := &c.xfmt.tokens[i], &xfmt.tokens[i]

		if !tokenCompare(tw, tg) {
			return fmt.Errorf(errParseFormatMismatchToken, i, tw, tg)
		}
	}

	return nil
}

const (
	tpfc1rawstring = "simple raw string without format `verbs`"
)

const (
	// fmt

	// %%
	tpfc1fmt_pp = "%%"

	// no verb
	tpfc1fmtnv = "%[3]*.*"

	// badVerb
	tpfc1fmtvb           = "z"
	tpfc1fmt_b           = "%z"
	tpfc1fmt_bi1_wi2_pi3 = "%[2]*.[3]*[1]z"
	tpfc1fmt_b_wi2_pi3   = "%[2]*.[3]*z"

	// %s
	tpfc1fmt_s             = "%s"
	tpfc1fmt_si3           = "%[3]s"
	tpfc1fmt_s_w20         = "%20s"
	tpfc1fmt_s_w20_p7      = "%20.7s"
	tpfc1fmt_s_w20_p7_fm   = "%-20.7s"
	tpfc1fmt_si1_w20_p7_fm = "%-20.7[3]s"
	tpfc1fmt_s_wi1_p_fs    = "%#[1]*.s"
	tpfc1fmt_s_p7_fm       = "%-.7s"
	tpfc1fmt_si1_pi2       = "%.[2]*[1]s"
	tpfc1fmt_si1_wi2_pi3   = "%[2]*.[3]*[1]s"

	// %q
	tpfc1fmt_q         = "%q"
	tpfc1fmt_q_w20_p13 = "%20.13q"
	tpfc1fmt_q_fs      = "%#q"
	tpfc1fmt_q_fp      = "%+q"
	tpfc1fmt_q_fps     = "%+#q"

	// %x
	tpfc1fmt_x             = "%x"
	tpfc1fmt_x_w19_p13     = "%19.13x"
	tpfc1fmt_x_w23_p7_fm   = "%-23.7x"
	tpfc1fmt_x_f_          = "% x"
	tpfc1fmt_x_fs          = "%#x"
	tpfc1fmt_x_fs_         = "%# x"
	tpfc1fmt_x_w37_p5_fs_  = "%# 37.5x"
	tpfc1fmt_x_w37_p5_fsm_ = "%# -37.5x"

	// %X
	tpfc1fmt_X             = "%X"
	tpfc1fmt_X_w19_p13     = "%19.13X"
	tpfc1fmt_X_w23_p7_fm   = "%-23.7X"
	tpfc1fmt_X_f_          = "% X"
	tpfc1fmt_X_fs          = "%#X"
	tpfc1fmt_X_fs_         = "%# X"
	tpfc1fmt_X_w37_p5_fs_  = "%# 37.5X"
	tpfc1fmt_X_w37_p5_fsm_ = "%# -37.5X"
)

const (
	// args
	tpfc1argword1 = "singleword"
	tpfc1argword2 = "short"
	tpfc1argstr1  = "some varied string - 123_456;"
	tpfc1argstr2  = "another & different string ^."
	tpfc1argpct1  = "percent % inside arg str"
	tpfc1argbq1   = "backquoted `string` for testing %q"
	tpfc1argbqqt1 = "arg with `quotes inside \"backquotes\"`"
	tpfc1argutf8  = "значение аргумента в иной кодировке"
	tpfc1argutf8c = "комбинированный utf8 string с внедренныеми символами № # "
	//
)

var parseFormatTestCases1 = [...]parseFormatTestCase{
	// empty string
	{
		emptyString,
		xfmt{
			nil, // nil for empty format string
			0,
			len(emptyString),
		},
	},
	// raw const string value
	{
		tpfc1rawstring,
		xfmt{
			[]token{
				{
					verb:  verbNone,
					value: tpfc1rawstring,
				},
			},
			0,
			len(tpfc1rawstring),
		},
	},
	// no verb
	{
		tpfc1fmtnv,
		xfmt{
			[]token{
				tokenErrNoVerb,
			},
			4,
			0,
		},
	},
	// badVerb
	// simple "%z"
	{
		tpfc1fmt_b,
		xfmt{
			[]token{
				{
					verb:  badVerb,
					value: tpfc1fmtvb,
					flags: flagNone,
					width: absentValue,
					prec:  absentValue,
				},
			},
			1,
			0,
		},
	},
	// "%[2]*.[3]*[1]z" badVerb with indir width and prec, and explicit indexing arg
	{
		tpfc1fmt_bi1_wi2_pi3,
		xfmt{
			[]token{
				{
					verb:  badVerb,
					value: tpfc1fmtvb,
					flags: flagIndirectWidth | flagIndirectPrec,
					width: 2 - 1,
					prec:  3 - 1,
					arg:   1 - 1,
				},
			},
			3,
			0,
		},
	},
	// "%[2]*.[3]*z" badVerb with indir width and prec, and implicit indexing arg
	{
		tpfc1fmt_b_wi2_pi3,
		xfmt{
			[]token{
				{
					verb:  badVerb,
					value: tpfc1fmtvb,
					flags: flagIndirectWidth | flagIndirectPrec,
					width: 2 - 1,
					prec:  3 - 1,
					arg:   4 - 1,
				},
			},
			4,
			0,
		},
	},
	// "%%" token
	{
		tpfc1fmt_pp,
		xfmt{
			[]token{
				tokenPercent,
			},
			0,
			1,
		},
	},
	// simple "%s" token
	{
		tpfc1fmt_s,
		xfmt{
			[]token{
				{
					verb:  verbString,
					value: string(verbCharString),
					width: absentValue,
					prec:  absentValue,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	// simple "%s" token with arg num
	{
		tpfc1fmt_si3,
		xfmt{
			[]token{
				{
					verb:  verbString,
					value: string(verbCharString),
					width: absentValue,
					prec:  absentValue,
					arg:   3 - 1, // args are indexing from 1
				},
			},
			3,
			0,
		},
	},
	// "%s" token with width
	{
		tpfc1fmt_s_w20,
		xfmt{
			[]token{
				{
					verb:  verbString,
					value: string(verbCharString),
					width: 20,
					prec:  absentValue,
					arg:   1 - 1, // args are indexing from 1
				},
			},
			1,
			0,
		},
	},
	// "%s" token with width and prec
	{
		tpfc1fmt_s_w20_p7,
		xfmt{
			[]token{
				{
					verb:  verbString,
					value: string(verbCharString),
					width: 20,
					prec:  7,
					arg:   1 - 1, // args are indexing from 1
				},
			},
			1,
			0,
		},
	},
	// "%s" token with width, prec and flag
	{
		tpfc1fmt_s_w20_p7_fm,
		xfmt{
			[]token{
				{
					verb:  verbString,
					value: string(verbCharString),
					flags: flagMinus,
					width: 20,
					prec:  7,
					arg:   1 - 1, // args are indexing from 1
				},
			},
			1,
			0,
		},
	},
	// "%s" token with width, prec, flag and arg num verb
	{
		tpfc1fmt_si1_w20_p7_fm,
		xfmt{
			[]token{
				{
					verb:  verbString,
					value: string(verbCharString),
					flags: flagMinus,
					width: 20,
					prec:  7,
					arg:   3 - 1, // args are indexing from 1
				},
			},
			3,
			0,
		},
	},
	// "%s" token with indir width, prec with no value and  flag sharp
	{
		tpfc1fmt_s_wi1_p_fs,
		xfmt{
			[]token{
				{
					verb:  verbString,
					value: string(verbCharString),
					flags: flagIndirectWidth | flagSharp,
					width: 1 - 1,
					prec:  0,
					arg:   2 - 1,
				},
			},
			2,
			0,
		},
	},
	// "%s" token without width, with prec and right pad
	{
		tpfc1fmt_s_p7_fm,
		xfmt{
			[]token{
				{
					verb:  verbString,
					value: string(verbCharString),
					flags: flagMinus,
					width: absentValue,
					prec:  7,
					arg:   0,
				},
			},
			1,
			0,
		},
	},
	// "%s" token with prec indir and arg num
	{
		tpfc1fmt_si1_pi2,
		xfmt{
			[]token{
				{
					verb:  verbString,
					value: string(verbCharString),
					flags: flagIndirectPrec,
					width: absentValue,
					prec:  2 - 1,
					arg:   1 - 1,
				},
			},
			2,
			0,
		},
	},
	// "%s" token with indir width, prec and arg num
	{
		tpfc1fmt_si1_wi2_pi3,
		xfmt{
			[]token{
				{
					verb:  verbString,
					value: string(verbCharString),
					flags: flagIndirectWidth | flagIndirectPrec,
					width: 2 - 1,
					prec:  3 - 1,
					arg:   1 - 1,
				},
			},
			3,
			0,
		},
	},
	// simple "%q" token
	{
		tpfc1fmt_q,
		xfmt{
			[]token{
				{
					verb:  verbQuoted,
					value: string(verbCharQuoted),
					flags: flagNone,
					width: absentValue,
					prec:  absentValue,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	// "%q" token with width and prec
	{
		tpfc1fmt_q_w20_p13,
		xfmt{
			[]token{
				{
					verb:  verbQuoted,
					value: string(verbCharQuoted),
					flags: flagNone,
					width: 20,
					prec:  13,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	// "%q" token with flag plus
	{
		tpfc1fmt_q_fp,
		xfmt{
			[]token{
				{
					verb:  verbQuoted,
					value: string(verbCharQuoted),
					flags: flagPlus,
					width: absentValue,
					prec:  absentValue,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	// "%q" token with flag sharp
	{
		tpfc1fmt_q_fs,
		xfmt{
			[]token{
				{
					verb:  verbQuoted,
					value: string(verbCharQuoted),
					flags: flagSharp,
					width: absentValue,
					prec:  absentValue,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	// "%q" token with plus and sharp flags
	{
		tpfc1fmt_q_fps,
		xfmt{
			[]token{
				{
					verb:  verbQuoted,
					value: string(verbCharQuoted),
					flags: flagSharp | flagPlus,
					width: absentValue,
					prec:  absentValue,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	// "%x" token
	{
		tpfc1fmt_x,
		xfmt{
			[]token{
				{
					verb:  verbHex,
					value: string(verbCharHex),
					flags: flagNone,
					width: absentValue,
					prec:  absentValue,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	// "%x" token with width and prec
	{
		tpfc1fmt_x_w19_p13,
		xfmt{
			[]token{
				{
					verb:  verbHex,
					value: string(verbCharHex),
					flags: flagNone,
					width: 19,
					prec:  13,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	// "%x" token with width, prec and minus flag
	{
		tpfc1fmt_x_w23_p7_fm,
		xfmt{
			[]token{
				{
					verb:  verbHex,
					value: string(verbCharHex),
					flags: flagMinus,
					width: 23,
					prec:  7,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	// "%x" token with space flag
	{
		tpfc1fmt_x_f_,
		xfmt{
			[]token{
				{
					verb:  verbHex,
					value: string(verbCharHex),
					flags: flagSpace,
					width: absentValue,
					prec:  absentValue,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	// "%x" token with sharp flag
	{
		tpfc1fmt_x_fs,
		xfmt{
			[]token{
				{
					verb:  verbHex,
					value: string(verbCharHex),
					flags: flagSharp,
					width: absentValue,
					prec:  absentValue,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	// "%x" token with space and sharp flags
	{
		tpfc1fmt_x_fs_,
		xfmt{
			[]token{
				{
					verb:  verbHex,
					value: string(verbCharHex),
					flags: flagSharp | flagSpace,
					width: absentValue,
					prec:  absentValue,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	// "%x" token with width, prec, space and sharp flags
	{
		tpfc1fmt_x_w37_p5_fs_,
		xfmt{
			[]token{
				{
					verb:  verbHex,
					value: string(verbCharHex),
					flags: flagSharp | flagSpace,
					width: 37,
					prec:  5,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	// "%x" token with width, prec, space, minus and sharp flags
	{
		tpfc1fmt_x_w37_p5_fsm_,
		xfmt{
			[]token{
				{
					verb:  verbHex,
					value: string(verbCharHex),
					flags: flagSharp | flagSpace | flagMinus,
					width: 37,
					prec:  5,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	// "%X" token
	{
		tpfc1fmt_X,
		xfmt{
			[]token{
				{
					verb:  verbHex,
					value: string(unicode.ToUpper(verbCharHex)),
					flags: flagUpperVerb,
					width: absentValue,
					prec:  absentValue,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	// "%X" token with width and prec
	{
		tpfc1fmt_X_w19_p13,
		xfmt{
			[]token{
				{
					verb:  verbHex,
					value: string(unicode.ToUpper(verbCharHex)),
					flags: flagUpperVerb,
					width: 19,
					prec:  13,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	// "%X" token with width, prec and minus flag
	{
		tpfc1fmt_X_w23_p7_fm,
		xfmt{
			[]token{
				{
					verb:  verbHex,
					value: string(unicode.ToUpper(verbCharHex)),
					flags: flagUpperVerb | flagMinus,
					width: 23,
					prec:  7,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	// "%X" token with space flag
	{
		tpfc1fmt_X_f_,
		xfmt{
			[]token{
				{
					verb:  verbHex,
					value: string(unicode.ToUpper(verbCharHex)),
					flags: flagUpperVerb | flagSpace,
					width: absentValue,
					prec:  absentValue,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	// "%X" token with sharp flag
	{
		tpfc1fmt_X_fs,
		xfmt{
			[]token{
				{
					verb:  verbHex,
					value: string(unicode.ToUpper(verbCharHex)),
					flags: flagUpperVerb | flagSharp,
					width: absentValue,
					prec:  absentValue,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	// "%X" token with space and sharp flags
	{
		tpfc1fmt_X_fs_,
		xfmt{
			[]token{
				{
					verb:  verbHex,
					value: string(unicode.ToUpper(verbCharHex)),
					flags: flagUpperVerb | flagSharp | flagSpace,
					width: absentValue,
					prec:  absentValue,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	// "%X" token with width, prec, space and sharp flags
	{
		tpfc1fmt_X_w37_p5_fs_,
		xfmt{
			[]token{
				{
					verb:  verbHex,
					value: string(unicode.ToUpper(verbCharHex)),
					flags: flagUpperVerb | flagSharp | flagSpace,
					width: 37,
					prec:  5,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	// "%X" token with width, prec, space, minus and sharp flags
	{
		tpfc1fmt_X_w37_p5_fsm_,
		xfmt{
			[]token{
				{
					verb:  verbHex,
					value: string(unicode.ToUpper(verbCharHex)),
					flags: flagUpperVerb | flagSharp | flagSpace | flagMinus,
					width: 37,
					prec:  5,
					arg:   1 - 1,
				},
			},
			1,
			0,
		},
	},
	/*
		// "%s" token with
		{
			tpfc1fmt_si3,
			xfmt{
				[]token{
					{

					},
				},
				0,
				0,
			},
		},
	*/
}

// TODO xfmt <-> fmt test
// TODO tests for error cases and eq.
// TODO adapted original fmt tests

// go test -count=1 -v -run "^TestParseFormatCases1$"
// go test -count=1 -o fmt.exe -gcflags "-m -m -d=ssa/check_bce/debug=1" -v -run "^TestParseFormatCases1$" 2> fmt.log
// go tool objdump -S -s "go-xfmt" fmt.exe > fmt.disasm
func TestParseFormatCases1(t *testing.T) {

	nerrs, sz := 0, len(parseFormatTestCases1)

	for i := 0; i < sz; i++ {

		testCase := &parseFormatTestCases1[i]

		if err := testCase.run(t); err != nil {
			t.Errorf("testCase %d %q error: %v", i, testCase.format, err)
			nerrs++
		}
	}

	if nerrs > 0 {
		t.Fatalf("some tests finished with errors: %d of %d", nerrs, sz)
	}
}

// TODO

// go test -count=1 -v -run "^TestAssertWithFmtPkg1$"
// go test -count=1 -o fmt.exe -gcflags "-m -m -d=ssa/check_bce/debug=1" -v -run "^TestAssertWithFmtPkg1$" 2> fmt.log
// go tool objdump -S -s "go-xfmt" fmt.exe > fmt.disasm
func TestAssertWithFmtPkg1(t *testing.T) {
	t.Log(esprintf("%% ||| %s ||| %q ||| % #x", tpfc1argpct1, tpfc1argbqqt1, tpfc1argutf8c))
}
