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

type tokenTestCase struct {
	token  token
	args   []string
	result string
}

const errPtnTestCaseResultMismatch = "token result value mismatch: want %q, got %q"

var errTestCaseFormatFail = errors.New("format fail (not ok)")

func (tcase *tokenTestCase) run( /*t *testing.T*/ ) error {

	//var buf strings.Builder
	//var buf bytebufferpool.ByteBuffer
	var buf buffer

	/*if ok := tcase.token.format(&buf, tcase.args); !ok {
		return errTestCaseFormatFail
	}*/

	tcase.token.format(&buf, tcase.args)

	if result := buf.String(); result != tcase.result {
		return fmt.Errorf(errPtnTestCaseResultMismatch, tcase.result, result)
	}

	//t.Logf("%v (%#v) ==> %s <<<EOL", tcase.token, tcase.args, buf.String())

	return nil
}

const (
	ttc1rawsimple     = "simple string const raw value"
	ttc1rawrndstr     = "AaHixWLu8RrpJ RxmbnuLrSC2pL GwN7mYpVSARWG tdjDCBwaHzDby"
	ttc1arg1str       = "str test arg value 1"
	ttc1arg2str       = "other test arg value 2"
	ttc1arg3str       = "another test arg value #3"
	ttc1arg4str       = "let once more test arg v 4"
	ttc1arg5str       = "so many test args v 5"
	ttc1arg6str       = "short6"
	ttc1arg7strquoted = "qstring with \" and bq ` and \t and \b so many 7"
	ttc1argutf8str    = "еще строка для теста %q :: utf-8 1"
	ttc1argutf8bqstr  = "дополнительная non-ascii для %q, но теперь с ` (backquote)"
)

var (
	vttc1arg6str               interface{} = ttc1arg6str
	vttc1badverb1fmt                       = "%z"
	vttc1missingargbadverb1fmt             = "%[2]*.[3]*[1]z"
	vttc1missingargbadverb2fmt             = "%[2]*.[3]*z"
	vttc1missingarg1fmt                    = "%[2]*.[3]*[1]s"
	vttc1missingarg2fmt                    = "%[2]*.[3]*s"
	vttc1missingarg3fmt                    = "%[1]*s"
)

var tokenTestCases1 = [...]tokenTestCase{
	{
		// simple raw const string value
		token: token{
			verb:  verbNone,
			value: ttc1rawsimple,
		},
		args:   nil,
		result: ttc1rawsimple,
	},
	{
		// verbNone with wrong values which should be ignored
		token: token{
			verb:  verbNone,
			value: ttc1rawrndstr,
			flags: flagMinus | flagPlus | flagSpace | flagSharp | flagZero,
			width: 50000,
			prec:  1,
			arg:   10000,
		},
		args:   []string{"unusedArg1", "unusedArg2"},
		result: ttc1rawrndstr,
	},
	{
		// "%z" badVerb
		token: token{
			verb:  badVerb,
			value: "z",
			flags: flagNone,
			width: absentValue,
			prec:  absentValue,
			arg:   0,
		},
		args:   []string{ttc1arg5str},
		result: fmt.Sprintf(vttc1badverb1fmt, ttc1arg5str),
	},
	{
		// "%[2]*.[3]*[1]z" badVerb with wrong indir width and prec and existing arg
		token: token{
			verb:  badVerb,
			value: "z",
			flags: flagIndirectWidth | flagIndirectPrec,
			width: 1,
			prec:  2,
			arg:   0,
		},
		args:   []string{ttc1arg5str, ttc1arg3str, ttc1arg6str},
		result: fmt.Sprintf(vttc1missingargbadverb1fmt, ttc1arg5str, ttc1arg3str, ttc1arg6str),
	},
	{
		// "%[2]*.[3]*z" badVerb with wrong indir width and prec and missing arg
		token: token{
			verb:  badVerb,
			value: "z",
			flags: flagIndirectWidth | flagIndirectPrec,
			width: 1,
			prec:  2,
			arg:   3,
		},
		args:   []string{ttc1arg5str, ttc1arg3str, ttc1arg6str},
		result: fmt.Sprintf(vttc1missingargbadverb2fmt, ttc1arg5str, ttc1arg3str, ttc1arg6str),
	},
	{
		// "%[2]*.[3]*[1]z" badVerb with missing indir width and prec and existing arg
		token: token{
			verb:  badVerb,
			value: "z",
			flags: flagIndirectWidth | flagIndirectPrec,
			width: 1,
			prec:  2,
			arg:   0,
		},
		args:   []string{ttc1arg5str},
		result: fmt.Sprintf(vttc1missingargbadverb1fmt, ttc1arg5str),
	},
	{
		// "%[2]*.[3]*z" badVerb with missing indir width, prec and arg
		token: token{
			verb:  badVerb,
			value: "z",
			flags: flagIndirectWidth | flagIndirectPrec,
			width: 1,
			prec:  2,
			arg:   3,
		},
		args:   []string{ttc1arg5str},
		result: fmt.Sprintf(vttc1missingargbadverb2fmt, ttc1arg5str),
	},
	{
		// "%%" token
		token: token{
			verb:  verbNone,
			value: percentString,
		},
		args:   []string{"argx1", "Argv2"},
		result: percentString,
	},
	{
		// simple "%s" token
		token: token{
			verb:  verbString,
			value: string(verbCharString),
			width: absentValue,
			prec:  absentValue,
			arg:   0,
		},
		args:   []string{ttc1arg1str},
		result: ttc1arg1str,
	},
	{
		// simple "%s" token with arg num
		token: token{
			verb:  verbString,
			value: string(verbCharString),
			width: absentValue,
			prec:  absentValue,
			arg:   3 - 1, // args are indexing from 1
		},
		args:   []string{ttc1arg1str, ttc1arg2str, ttc1arg3str, ttc1arg4str, ttc1arg5str},
		result: ttc1arg3str,
	},
	{
		// simple "%[1]*s" token with missing arg
		token: token{
			verb:  verbString,
			value: string(verbCharString),
			width: absentValue,
			prec:  absentValue,
			arg:   3 - 1, // args are indexing from 1
		},
		args: []string{"10"},
		// emulate `missing arg` as width arg + 1 absent (missing) arg for verb, so must use int for width to avoid unnecessary test for bad width arg value
		result: fmt.Sprintf(vttc1missingarg3fmt, 10),
	},
	{
		// "%10.3s" - prec and width left pad
		token: token{
			verb:  verbString,
			value: string(verbCharString),
			width: 10,
			prec:  3,
			arg:   0,
		},
		args:   []string{ttc1arg3str},
		result: fmt.Sprintf("%10.3s", ttc1arg3str),
	},
	{
		// "%-10s" - width right pad
		token: token{
			verb:  verbString,
			value: string(verbCharString),
			flags: flagMinus,
			width: 10,
			prec:  absentValue,
			arg:   0,
		},
		args:   []string{ttc1arg6str},
		result: fmt.Sprintf("%-10s", ttc1arg6str),
	},
	{
		// "%[1]*s" - width indir arg exists
		token: token{
			verb:  verbString,
			value: string(verbCharString),
			flags: flagIndirectWidth,
			width: 1 - 1,
			prec:  absentValue,
			arg:   2 - 1,
		},
		args:   []string{ttc1arg6str, ttc1arg5str},
		result: fmt.Sprintf(vttc1missingarg3fmt, vttc1arg6str, ttc1arg5str),
	},
	{
		// "%[1]*s" - width missing indir arg with missing verb arg
		token: token{
			verb:  verbString,
			value: string(verbCharString),
			flags: flagIndirectWidth,
			width: 1 - 1,
			prec:  absentValue,
			arg:   2 - 1,
		},
		args:   []string{ttc1arg6str},
		result: fmt.Sprintf(vttc1missingarg3fmt, vttc1arg6str),
	},

	{
		// "%-.7s" - prec right pad
		token: token{
			verb:  verbString,
			value: string(verbCharString),
			flags: flagMinus,
			width: absentValue,
			prec:  7,
			arg:   0,
		},
		args:   []string{ttc1arg2str},
		result: fmt.Sprintf("%-.7s", ttc1arg2str),
	},
	{
		// "%.[1]*s" - prec indir arg exists
		token: token{
			verb:  verbString,
			value: string(verbCharString),
			flags: flagIndirectPrec,
			width: absentValue,
			prec:  1,
			arg:   2 - 1,
		},
		args:   []string{ttc1arg6str, ttc1arg5str},
		result: fmt.Sprintf("%.[1]*s", vttc1arg6str, ttc1arg5str),
	},
	{
		// "%[2]*.[3]*[1]s" wrong indir width and prec and existing arg
		token: token{
			verb:  verbString,
			value: string(verbCharString),
			flags: flagIndirectWidth | flagIndirectPrec,
			width: 1,
			prec:  2,
			arg:   0,
		},
		args:   []string{ttc1arg5str, ttc1arg3str, ttc1arg6str},
		result: fmt.Sprintf(vttc1missingarg1fmt, ttc1arg5str, ttc1arg3str, ttc1arg6str),
	},
	{
		// "%[2]*.[3]*[1]s" wrong indir width and prec and missing arg
		token: token{
			verb:  verbString,
			value: string(verbCharString),
			flags: flagIndirectWidth | flagIndirectPrec,
			width: 1,
			prec:  2,
			arg:   3,
		},
		args:   []string{ttc1arg5str, ttc1arg3str, ttc1arg6str},
		result: fmt.Sprintf(vttc1missingarg2fmt, ttc1arg5str, ttc1arg3str, ttc1arg6str),
	},
	{
		// "%[2]*.[3]*[1]s" missing indir width and prec and existing arg
		token: token{
			verb:  verbString,
			value: string(verbCharString),
			flags: flagIndirectWidth | flagIndirectPrec,
			width: 1,
			prec:  2,
			arg:   0,
		},
		args:   []string{ttc1arg5str},
		result: fmt.Sprintf(vttc1missingarg1fmt, ttc1arg5str),
	},
	{
		// "%[2]*.[3]*[1]s" missing indir width, prec and arg
		token: token{
			verb:  verbString,
			value: string(verbCharString),
			flags: flagIndirectWidth | flagIndirectPrec,
			width: 1,
			prec:  2,
			arg:   3,
		},
		args:   []string{ttc1arg5str},
		result: fmt.Sprintf(vttc1missingarg2fmt, ttc1arg5str),
	},
	{
		// "%q" - quoted string
		token: token{
			verb:  verbQuoted,
			value: string(verbCharQuoted),
			flags: flagNone,
			width: absentValue,
			prec:  absentValue,
			arg:   0,
		},
		args:   []string{ttc1arg4str},
		result: fmt.Sprintf("%q", ttc1arg4str),
	},
	{
		// "%20.13q" - quoted string + width prec
		token: token{
			verb:  verbQuoted,
			value: string(verbCharQuoted),
			flags: flagNone,
			width: 20,
			prec:  13,
			arg:   0,
		},
		args:   []string{ttc1arg3str},
		result: fmt.Sprintf("%20.13q", ttc1arg3str),
	},
	{
		// "%q" - quoted string + arg with non-ascii
		token: token{
			verb:  verbQuoted,
			value: string(verbCharQuoted),
			flags: flagNone,
			width: absentValue,
			prec:  absentValue,
			arg:   0,
		},
		args:   []string{ttc1argutf8str},
		result: fmt.Sprintf("%q", ttc1argutf8str),
	},
	{
		// "%q" - quoted string + arg with backquote
		token: token{
			verb:  verbQuoted,
			value: string(verbCharQuoted),
			flags: flagNone,
			width: absentValue,
			prec:  absentValue,
			arg:   0,
		},
		args:   []string{ttc1arg7strquoted},
		result: fmt.Sprintf("%q", ttc1arg7strquoted),
	},
	{
		// "%#q" - alt quoted string
		token: token{
			verb:  verbQuoted,
			value: string(verbCharQuoted),
			flags: flagSharp,
			width: absentValue,
			prec:  absentValue,
			arg:   0,
		},
		args:   []string{ttc1arg3str},
		result: fmt.Sprintf("%#q", ttc1arg3str),
	},
	{
		// "%#q" - alt quoted string + arg with backquote
		token: token{
			verb:  verbQuoted,
			value: string(verbCharQuoted),
			flags: flagSharp,
			width: absentValue,
			prec:  absentValue,
			arg:   0,
		},
		args:   []string{ttc1arg7strquoted},
		result: fmt.Sprintf("%#q", ttc1arg7strquoted),
	},
	{
		// "%+q" - quoted string + ascii only
		token: token{
			verb:  verbQuoted,
			value: string(verbCharQuoted),
			flags: flagPlus,
			width: absentValue,
			prec:  absentValue,
			arg:   0,
		},
		args:   []string{ttc1arg5str},
		result: fmt.Sprintf("%+q", ttc1arg5str),
	},
	{
		// "%+q" - quoted string + ascii only + arg with backquote
		token: token{
			verb:  verbQuoted,
			value: string(verbCharQuoted),
			flags: flagPlus,
			width: absentValue,
			prec:  absentValue,
			arg:   0,
		},
		args:   []string{ttc1arg7strquoted},
		result: fmt.Sprintf("%+q", ttc1arg7strquoted),
	},
	{
		// "%+q" - quoted string + ascii only + arg with non-ascii
		token: token{
			verb:  verbQuoted,
			value: string(verbCharQuoted),
			flags: flagPlus,
			width: absentValue,
			prec:  absentValue,
			arg:   0,
		},
		args:   []string{ttc1argutf8str},
		result: fmt.Sprintf("%+q", ttc1argutf8str),
	},
	{
		// "%+q" - quoted string + ascii only + alt fmt + arg with non-ascii
		token: token{
			verb:  verbQuoted,
			value: string(verbCharQuoted),
			flags: flagPlus | flagSharp,
			width: absentValue,
			prec:  absentValue,
			arg:   0,
		},
		args:   []string{ttc1argutf8str},
		result: fmt.Sprintf("%+#q", ttc1argutf8str),
	},
	{
		// "%+#q" - quoted string + ascii only + alt fmt + arg with non-ascii
		token: token{
			verb:  verbQuoted,
			value: string(verbCharQuoted),
			flags: flagPlus | flagSharp,
			width: absentValue,
			prec:  absentValue,
			arg:   0,
		},
		args:   []string{ttc1argutf8bqstr},
		result: fmt.Sprintf("%+#q", ttc1argutf8bqstr),
	},
	{
		// "%x" - hex string
		token: token{
			verb:  verbHex,
			value: string(verbCharHex),
			flags: flagNone,
			width: absentValue,
			prec:  absentValue,
			arg:   0,
		},
		args:   []string{ttc1arg7strquoted},
		result: fmt.Sprintf("%x", ttc1arg7strquoted),
	},
	{
		// "%x" - hex string upper
		token: token{
			verb:  verbHex,
			value: string(verbCharHex),
			flags: flagUpperVerb,
			width: absentValue,
			prec:  absentValue,
			arg:   0,
		},
		args:   []string{ttc1argutf8str},
		result: fmt.Sprintf("%X", ttc1argutf8str),
	},
	{
		// "%19.13x" - hex string width prec
		token: token{
			verb:  verbHex,
			value: string(verbCharHex),
			flags: flagNone,
			width: 19,
			prec:  13,
			arg:   0,
		},
		args:   []string{ttc1arg3str},
		result: fmt.Sprintf("%19.13x", ttc1arg3str),
	},
	{
		// "%23.7x" - hex string width prec pad left
		token: token{
			verb:  verbHex,
			value: string(verbCharHex),
			flags: flagNone,
			width: 23,
			prec:  7,
			arg:   0,
		},
		args:   []string{ttc1arg1str},
		result: fmt.Sprintf("%23.7x", ttc1arg1str),
	},
	{
		// "%-23.7x" - hex string width prec pad right
		token: token{
			verb:  verbHex,
			value: string(verbCharHex),
			flags: flagMinus,
			width: 23,
			prec:  7,
			arg:   0,
		},
		args:   []string{ttc1arg4str},
		result: fmt.Sprintf("%-23.7x", ttc1arg4str),
	},
	{
		// "% x" - hex string + space
		token: token{
			verb:  verbHex,
			value: string(verbCharHex),
			flags: flagSpace,
			width: absentValue,
			prec:  absentValue,
			arg:   0,
		},
		args:   []string{ttc1argutf8str},
		result: fmt.Sprintf("% x", ttc1argutf8str),
	},
	{
		// "%#x" - hex string + alt fmt
		token: token{
			verb:  verbHex,
			value: string(verbCharHex),
			flags: flagSharp,
			width: absentValue,
			prec:  absentValue,
			arg:   0,
		},
		args:   []string{ttc1arg7strquoted},
		result: fmt.Sprintf("%#x", ttc1arg7strquoted),
	},
	{
		// "%# x" - hex string + alt fmt + space
		token: token{
			verb:  verbHex,
			value: string(verbCharHex),
			flags: flagSharp | flagSpace,
			width: absentValue,
			prec:  absentValue,
			arg:   0,
		},
		args:   []string{ttc1arg5str},
		result: fmt.Sprintf("%# x", ttc1arg5str),
	},
	{
		// "%# 37.5x" - hex string + alt fmt + space + width + prec
		token: token{
			verb:  verbHex,
			value: string(verbCharHex),
			flags: flagSharp | flagSpace,
			width: 37,
			prec:  5,
			arg:   0,
		},
		args:   []string{ttc1argutf8bqstr},
		result: fmt.Sprintf("%# 37.5x", ttc1argutf8bqstr),
	},
	{
		// "%# -37.5x" - hex string + alt fmt + space + width + prec + pad right
		token: token{
			verb:  verbHex,
			value: string(verbCharHex),
			flags: flagSharp | flagSpace | flagMinus,
			width: 37,
			prec:  5,
			arg:   0,
		},
		args:   []string{ttc1arg2str},
		result: fmt.Sprintf("%# -37.5x", ttc1arg2str),
	},
}

// go test -count=1 -v -run "^TestTokenCases1$"
// go test -count=1 -o fmt.exe -gcflags "-m -m -d=ssa/check_bce/debug=1" -v -run "^TestTokenCases1$" 2> fmt.log
// go tool objdump -S -s "go-xfmt" fmt.exe > fmt.disasm
func TestTokenCases1(t *testing.T) {

	nerrs, sz := 0, len(tokenTestCases1)

	for i := 0; i < sz; i++ {

		tcase := &tokenTestCases1[i]

		if err := tcase.run( /*t*/ ); err != nil {
			t.Errorf("%d (%#v => %v): %v", i, tcase.token, tcase.args, err)
			nerrs++
		}
	}

	if nerrs > 0 {
		t.Fatalf("some tests finished with errors: %d of %d", nerrs, sz)
	}
}
