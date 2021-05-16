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
	"errors"
	"fmt"
	"testing"
)

type flagTestCase struct {
	format string
	verb   verb
	flags  flags
	width  int
	prec   int
	c      int
}

var (
	errFTCAbsentTokens = errors.New("no tokens")
)

func (tc *flagTestCase) run( /*t *testing.T*/ ) error {

	x := parseFormat(tc.format)

	if len(x.tokens) == 0 {
		return errFTCAbsentTokens
	}

	if len(x.tokens) != tc.c {
		return fmt.Errorf("wrong tokens count: want %d, got %d", tc.c, len(x.tokens))
	}

	tt := &x.tokens[0]

	//t.Log(tt)

	if tt.verb != tc.verb {
		return fmt.Errorf("verb mismatch: want %c, got %c", tc.verb, tt.verb)
	}

	if tt.flags != tc.flags {
		return fmt.Errorf("flags mismatch: want %#x, got %#x", tc.flags, tt.flags)
	}

	if tt.width != tc.width {
		return fmt.Errorf("width mismatch: want %d, got %d", tc.width, tt.width)
	}

	if tt.prec != tc.prec {
		return fmt.Errorf("prec mismatch: want %d, got %d", tc.prec, tt.prec)
	}

	const arg = "some arg test value 123467\t\n"

	r := x.Sprint([]string{arg})

	if w := fmt.Sprintf(tc.format, arg); w != r {
		return fmt.Errorf("unexpected result: want <%s>, got <%s>", w, r)
	}

	return nil
}

var flagTestCases = [...]flagTestCase{
	{
		"%s",
		verbString,
		flagNone,
		absentValue,
		absentValue,
		1,
	},
	{
		"%q",
		verbQuoted,
		flagNone,
		absentValue,
		absentValue,
		1,
	},
	{
		"%x",
		verbHex,
		flagNone,
		absentValue,
		absentValue,
		1,
	},
	{
		"%X",
		verbHex,
		flagUpperVerb,
		absentValue,
		absentValue,
		1,
	},
	// '-' flag
	{
		"%-s",
		verbString,
		flagMinus,
		absentValue,
		absentValue,
		1,
	},
	{
		"%-q",
		verbQuoted,
		flagMinus,
		absentValue,
		absentValue,
		1,
	},
	{
		"%-x",
		verbHex,
		flagMinus,
		absentValue,
		absentValue,
		1,
	},
	{
		"%-X",
		verbHex,
		flagMinus | flagUpperVerb,
		absentValue,
		absentValue,
		1,
	},
	// +
	{
		"%+s",
		verbString,
		flagPlus,
		absentValue,
		absentValue,
		1,
	},
	{
		"%+q",
		verbQuoted,
		flagPlus,
		absentValue,
		absentValue,
		1,
	},
	{
		"%+x",
		verbHex,
		flagPlus,
		absentValue,
		absentValue,
		1,
	},
	{
		"%+X",
		verbHex,
		flagPlus | flagUpperVerb,
		absentValue,
		absentValue,
		1,
	},
	// #
	{
		"%#s",
		verbString,
		flagSharp,
		absentValue,
		absentValue,
		1,
	},
	{
		"%#q",
		verbQuoted,
		flagSharp,
		absentValue,
		absentValue,
		1,
	},
	{
		"%#x",
		verbHex,
		flagSharp,
		absentValue,
		absentValue,
		1,
	},
	{
		"%#X",
		verbHex,
		flagSharp | flagUpperVerb,
		absentValue,
		absentValue,
		1,
	},
	// ' '
	{
		"% s",
		verbString,
		flagSpace,
		absentValue,
		absentValue,
		1,
	},
	{
		"% q",
		verbQuoted,
		flagSpace,
		absentValue,
		absentValue,
		1,
	},
	{
		"% x",
		verbHex,
		flagSpace,
		absentValue,
		absentValue,
		1,
	},
	{
		"% X",
		verbHex,
		flagSpace | flagUpperVerb,
		absentValue,
		absentValue,
		1,
	},
	// 0
	{
		"%0s",
		verbString,
		flagZero,
		absentValue,
		absentValue,
		1,
	},
	{
		"%0q",
		verbQuoted,
		flagZero,
		absentValue,
		absentValue,
		1,
	},
	{
		"%0x",
		verbHex,
		flagZero,
		absentValue,
		absentValue,
		1,
	},
	{
		"%0X",
		verbHex,
		flagZero | flagUpperVerb,
		absentValue,
		absentValue,
		1,
	},
	// width prec
	{
		"%1.2s",
		verbString,
		flagNone,
		1,
		2,
		1,
	},
	{
		"%1.2q",
		verbQuoted,
		flagNone,
		1,
		2,
		1,
	},
	{
		"%1.2x",
		verbHex,
		flagNone,
		1,
		2,
		1,
	},
	{
		"%1.2X",
		verbHex,
		flagUpperVerb,
		1,
		2,
		1,
	},
	// - w p
	{
		"%-1.2s",
		verbString,
		flagMinus,
		1,
		2,
		1,
	},
	{
		"%-1.2q",
		verbQuoted,
		flagMinus,
		1,
		2,
		1,
	},
	{
		"%-1.2x",
		verbHex,
		flagMinus,
		1,
		2,
		1,
	},
	{
		"%-1.2X",
		verbHex,
		flagMinus | flagUpperVerb,
		1,
		2,
		1,
	},
	// + w p
	{
		"%+1.2s",
		verbString,
		flagPlus,
		1,
		2,
		1,
	},
	{
		"%+1.2q",
		verbQuoted,
		flagPlus,
		1,
		2,
		1,
	},
	{
		"%+1.2x",
		verbHex,
		flagPlus,
		1,
		2,
		1,
	},
	{
		"%+1.2X",
		verbHex,
		flagPlus | flagUpperVerb,
		1,
		2,
		1,
	},
	// + - w p
	{
		"%+-1.2s",
		verbString,
		flagPlus | flagMinus,
		1,
		2,
		1,
	},
	{
		"%+-1.2q",
		verbQuoted,
		flagPlus | flagMinus,
		1,
		2,
		1,
	},
	{
		"%+-1.2x",
		verbHex,
		flagPlus | flagMinus,
		1,
		2,
		1,
	},
	{
		"%+-1.2X",
		verbHex,
		flagPlus | flagMinus | flagUpperVerb,
		1,
		2,
		1,
	},
	// + - w p tail
	{
		"%+-1.2sqx",
		verbString,
		flagPlus | flagMinus,
		1,
		2,
		2,
	},
	{
		"%+-1.2qsx",
		verbQuoted,
		flagPlus | flagMinus,
		1,
		2,
		2,
	},
	{
		"%+-1.2xsq",
		verbHex,
		flagPlus | flagMinus,
		1,
		2,
		2,
	},
	{
		"%+-1.2Xsq",
		verbHex,
		flagPlus | flagMinus | flagUpperVerb,
		1,
		2,
		2,
	},
	// - w p tail
	{
		"%-1.2sqx",
		verbString,
		flagMinus,
		1,
		2,
		2,
	},
	{
		"%-1.2qsx",
		verbQuoted,
		flagMinus,
		1,
		2,
		2,
	},
	{
		"%-1.2xsq",
		verbHex,
		flagMinus,
		1,
		2,
		2,
	},
	{
		"%-1.2Xsq",
		verbHex,
		flagMinus | flagUpperVerb,
		1,
		2,
		2,
	},
	// all flags except for minus
	{
		"%# +0[1]*.[2]*sqx",
		verbString,
		flagPlus | flagSpace | flagSharp | flagZero | flagIndirectWidth | flagIndirectPrec,
		1 - 1,
		2 - 1,
		2,
	},
	{
		"%# +0[1]*.[2]*qsx",
		verbQuoted,
		flagPlus | flagSpace | flagSharp | flagZero | flagIndirectWidth | flagIndirectPrec,
		1 - 1,
		2 - 1,
		2,
	},
	{
		"%# +0[1]*.[2]*xsq",
		verbHex,
		flagPlus | flagSpace | flagSharp | flagZero | flagIndirectWidth | flagIndirectPrec,
		1 - 1,
		2 - 1,
		2,
	},
	{
		"%# +0[1]*.[2]*Xsq",
		verbHex,
		flagPlus | flagSpace | flagSharp | flagZero | flagIndirectWidth | flagIndirectPrec | flagUpperVerb,
		1 - 1,
		2 - 1,
		2,
	},
	// all possible flags + tail
	// ATN! zero and minus are mutually exclusive: - disables 0, but 0 doesn't disable -, only doesn't applied
	//      ==> minus has higher priority
	{
		"%# +-0[1]*.[2]*sqx",
		verbString,
		flagMinus | flagPlus | flagSpace | flagSharp | flagIndirectWidth | flagIndirectPrec,
		1 - 1,
		2 - 1,
		2,
	},
	{
		"%# +-0[1]*.[2]*qsx",
		verbQuoted,
		flagMinus | flagPlus | flagSpace | flagSharp | flagIndirectWidth | flagIndirectPrec,
		1 - 1,
		2 - 1,
		2,
	},
	{
		"%# +-0[1]*.[2]*xsq",
		verbHex,
		flagMinus | flagPlus | flagSpace | flagSharp | flagIndirectWidth | flagIndirectPrec,
		1 - 1,
		2 - 1,
		2,
	},
	{
		"%# +-0[1]*.[2]*Xsq",
		verbHex,
		flagMinus | flagPlus | flagSpace | flagSharp | flagIndirectWidth | flagIndirectPrec | flagUpperVerb,
		1 - 1,
		2 - 1,
		2,
	},
}

// go test -count=1 -v -run "^TestFlagParser$"
func TestFlagParser(t *testing.T) {

	for i := 0; i < len(flagTestCases); i++ {

		tcase := &flagTestCases[i]

		if err := tcase.run( /*t*/ ); err != nil {
			t.Fatalf("%d (<%s>): %v", i, tcase.format, err)
		}
	}

}

//

type pickNumValueTestCase struct {
	format string
	start  uint
	num    int
	j      uint
}

func (tc *pickNumValueTestCase) run() error {

	num, j := pickNumValue(tc.format, tc.start)

	switch {
	case (num != tc.num) && (j != tc.j):
		return fmt.Errorf("mismatch `num` and `j` values: num want %d got %d, j want %d got %d", tc.num, num, tc.j, j)

	case num != tc.num:
		return fmt.Errorf("mismatch `num`: want %d got %d", tc.num, num)

	case j != tc.j:
		return fmt.Errorf("mismatch `j`: want %d got %d", tc.j, j)
	}

	return nil
}

var pickNumValueTestCases1 = [...]pickNumValueTestCase{
	{"1234", 4, absentValue, 4}, // start == end
	{"1234", 7, absentValue, 4}, // start > end
	{"a123", 0, absentValue, 0}, // starts from char not num - digits is not found

	{"123a", 0, 123, 3},
	{"12a3", 0, 12, 2},
	{"1234", 0, 1234, 4},

	{"1a234", 1, absentValue, 1},
	{"9999999x", 0, absentValue, 8}, // check overflow
}

// go test -count=1 -v -run "^TestPickNumValue$"
func TestPickNumValue(t *testing.T) {

	for i := 0; i < len(pickNumValueTestCases1); i++ {

		tcase := &pickNumValueTestCases1[i]

		if err := tcase.run(); err != nil {
			t.Fatalf("%d (%#v => %v): %v", i, tcase.format, tcase.start, err)
		}
	}

}
