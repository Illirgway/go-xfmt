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
	"bytes"
	"fmt"
	"runtime"
	"testing"
)

// adapted tests and benchmarks from fmt

type fmtTestCase struct {
	format string
	arg    string
	result string
}

//go:nosplit
func (c *fmtTestCase) sprintf() string {
	return Sprintf(c.format, c.arg)
}

const (
	tftc1arg1strrnd             = "arg1_QzR@QK+LIev1@@7rqK$KTFuKe%df0dI"
	tftc1arg2strhex             = "\xAB\xCD\xEF\xF0\xFF\x00\x01\xFF\xF0"
	tftc1arg2strhexVx           = "abcdeff0ff0001fff0"
	tftc1arg2strhexVX           = "ABCDEFF0FF0001FFF0"
	tftc1arg3short              = "xfmt."
	tftc1arg3shortVx            = "78666d742e"
	tftc1arg3shortVX            = "78666D742E"
	tftc1arg3shortVxf_          = "78 66 6d 74 2e"
	tftc1arg3shortVXf_          = "78 66 6D 74 2E"
	tftc1arg3shortVxfs          = "0x78666d742e"
	tftc1arg3shortVXfs          = "0X78666D742E"
	tftc1arg3shortVxfs_         = "0x78 0x66 0x6d 0x74 0x2e"
	tftc1arg3shortVXfs_         = "0X78 0X66 0X6D 0X74 0X2E"
	tftc1arg4utf8cjk            = "攱枬貮"
	tftc1arg4utf8cjkVqfp        = `"\u6531\u67ac\u8cae"` // quoted ascii only
	tftc1arg5utf8MS_HnS         = "☭"
	tftc1arg5utf8MS_HnS_U       = "\u262d"
	tftc1arg5utf8MS_HnS_RCode   = `\u262d`
	tftc1arg5utf8MS_HnSVqfp     = `"` + tftc1arg5utf8MS_HnS_RCode + `"` // quoted ascii only
	tftc1arg6utf8MT_Kbd         = "⌨"
	tftc1arg6utf8MT_Kbd_RCode   = `\u2328`
	tftc1arg7wrongbyte          = "xfmt\xffefg"
	tftc1arg7wrongbyteVq        = `xfmt\xffefg`
	tftc1arg8nonprintrune       = "\U0010fefd"
	tftc1arg8nonprintruneVq     = `\U0010fefd`
	tftc1arg8nonprintruneVq_Raw = `􏻽`
	tftc1arg9wrongrune          = string(rune(0x110000))
	tftc1arg10alphabet          = "abcdefghijklmnopqrstuvwxyz"
	tftc1arg11utf8cjkstr        = "朰朴杲枫未摁"
	tftc1arg11utf8cjkshort      = "朰朴杲"
	tftc1arg12str               = "bNfwVErqVWy"
	tftc1arg13str               = "13str"
	tftc1arg14utf8cjklong       = "亁亃乧殖每箚馘馛馞﨎忰忲忴貎豉豒"
	tftc1arg15str               = "0123456789abcdefgh"
)

var adaptedFmtTestCases1 = [...]fmtTestCase{

	// basic string
	{"%s", tftc1arg1strrnd, tftc1arg1strrnd},
	{"%q", tftc1arg1strrnd, `"` + tftc1arg1strrnd + `"`},
	{"%x", tftc1arg2strhex, tftc1arg2strhexVx},
	{"%X", tftc1arg2strhex, tftc1arg2strhexVX},
	//
	{"%s", emptyString, emptyString},
	{"%x", emptyString, emptyString},
	{"%X", emptyString, emptyString},
	{"% x", emptyString, emptyString},
	{"%#x", emptyString, emptyString},
	{"%# x", emptyString, emptyString},
	//
	{"%x", tftc1arg3short, tftc1arg3shortVx},
	{"%X", tftc1arg3short, tftc1arg3shortVX},
	{"% x", tftc1arg3short, tftc1arg3shortVxf_},
	{"% X", tftc1arg3short, tftc1arg3shortVXf_},
	{"%#x", tftc1arg3short, tftc1arg3shortVxfs},
	{"%#X", tftc1arg3short, tftc1arg3shortVXfs},
	{"%# x", tftc1arg3short, tftc1arg3shortVxfs_},
	{"%# X", tftc1arg3short, tftc1arg3shortVXfs_},

	// escaped strings
	{"%q", emptyString, `""`},
	{"%#q", emptyString, "``"},
	{"%q", "\"", `"\""`},
	{"%#q", "\"", "`\"`"},
	{"%q", "\n", `"\n"`},
	{"%#q", "\n", `"\n"`},
	{"%q", `\n`, `"\\n"`},
	{"%#q", `\n`, "`\\n`"},
	// ASCII
	{"%q", tftc1arg3short, `"` + tftc1arg3short + `"`},
	{"%#q", tftc1arg3short, "`" + tftc1arg3short + "`"},
	// CJK
	{"%q", tftc1arg4utf8cjk, `"` + tftc1arg4utf8cjk + `"`},
	{"%+q", tftc1arg4utf8cjk, tftc1arg4utf8cjkVqfp},
	{"%#q", tftc1arg4utf8cjk, "`" + tftc1arg4utf8cjk + "`"},
	{"%#+q", tftc1arg4utf8cjk, "`" + tftc1arg4utf8cjk + "`"},
	// ctl bqted chars
	{"%q", "\a\b\f\n\r\t\v\"\\", `"\a\b\f\n\r\t\v\"\\"`},
	{"%+q", "\a\b\f\n\r\t\v\"\\", `"\a\b\f\n\r\t\v\"\\"`},
	{"%#q", "\a\b\f\n\r\t\v\"\\", `"\a\b\f\n\r\t\v\"\\"`},
	{"%#+q", "\a\b\f\n\r\t\v\"\\", `"\a\b\f\n\r\t\v\"\\"`},
	// utf8 other symbols (So) Miscellaneous Symbols (MS)
	{"%q", tftc1arg5utf8MS_HnS, `"` + tftc1arg5utf8MS_HnS + `"`},
	{"% q", tftc1arg5utf8MS_HnS, `"` + tftc1arg5utf8MS_HnS + `"`}, // The space modifier should have no effect.
	{"%+q", tftc1arg5utf8MS_HnS, tftc1arg5utf8MS_HnSVqfp},
	{"%#q", tftc1arg5utf8MS_HnS, "`" + tftc1arg5utf8MS_HnS + "`"},
	{"%#+q", tftc1arg5utf8MS_HnS, "`" + tftc1arg5utf8MS_HnS + "`"},
	// utf8 other symbols (So) Miscellaneous Technical (MT)
	{"%10q", tftc1arg6utf8MT_Kbd, `       "` + tftc1arg6utf8MT_Kbd + `"`},
	{"%+10q", tftc1arg6utf8MT_Kbd, `  "` + tftc1arg6utf8MT_Kbd_RCode + `"`},
	{"%-10q", tftc1arg6utf8MT_Kbd, `"` + tftc1arg6utf8MT_Kbd + `"       `},
	{"%+-10q", tftc1arg6utf8MT_Kbd, `"` + tftc1arg6utf8MT_Kbd_RCode + `"  `},
	{"%010q", tftc1arg6utf8MT_Kbd, `0000000"` + tftc1arg6utf8MT_Kbd + `"`},
	{"%+010q", tftc1arg6utf8MT_Kbd, `00"` + tftc1arg6utf8MT_Kbd_RCode + `"`},
	{"%-010q", tftc1arg6utf8MT_Kbd, `"` + tftc1arg6utf8MT_Kbd + `"       `}, // '0' has no effect when '-' is present.
	{"%+-010q", tftc1arg6utf8MT_Kbd, `"` + tftc1arg6utf8MT_Kbd_RCode + `"  `},
	//
	{"%#8q", "\n", `    "\n"`},
	{"%#+8q", "\r", `    "\r"`},
	{"%#-8q", "\t", "`	`     "},
	{"%#+-8q", "\b", `"\b"    `},
	// wrong byte
	{"%q", tftc1arg7wrongbyte, `"` + tftc1arg7wrongbyteVq + `"`},
	{"%+q", tftc1arg7wrongbyte, `"` + tftc1arg7wrongbyteVq + `"`},
	{"%#q", tftc1arg7wrongbyte, `"` + tftc1arg7wrongbyteVq + `"`},
	{"%#+q", tftc1arg7wrongbyte, `"` + tftc1arg7wrongbyteVq + `"`},
	// Runes that are not printable.
	{"%q", tftc1arg8nonprintrune, `"` + tftc1arg8nonprintruneVq + `"`},
	{"%+q", tftc1arg8nonprintrune, `"` + tftc1arg8nonprintruneVq + `"`},
	{"%#q", tftc1arg8nonprintrune, "`" + tftc1arg8nonprintruneVq_Raw + "`"},
	{"%#+q", tftc1arg8nonprintrune, "`" + tftc1arg8nonprintruneVq_Raw + "`"},
	// Runes that are not valid.
	{"%q", tftc1arg9wrongrune, `"�"`},
	{"%+q", tftc1arg9wrongrune, `"\ufffd"`},
	{"%#q", tftc1arg9wrongrune, "`�`"},
	{"%#+q", tftc1arg9wrongrune, "`�`"},

	// width
	{"%7s", tftc1arg3short, "  " + tftc1arg3short},
	{"%2s", tftc1arg5utf8MS_HnS_U, " " + tftc1arg5utf8MS_HnS},
	{"%-7s", tftc1arg3short, tftc1arg3short + "  "},
	{"%07s", tftc1arg3short, "00" + tftc1arg3short},
	{"%5s", tftc1arg10alphabet, tftc1arg10alphabet},
	{"%.5s", tftc1arg10alphabet, tftc1arg10alphabet[:5]},
	{"%.0s", tftc1arg11utf8cjkstr, ""},
	{"%.5s", tftc1arg11utf8cjkstr, tftc1arg11utf8cjkstr[:5*3]},
	{"%.10s", tftc1arg11utf8cjkstr, tftc1arg11utf8cjkstr},
	{"%010q", tftc1arg3short, `000"` + tftc1arg3short + `"`},
	{"%-10q", tftc1arg3short, `"` + tftc1arg3short + `"   `},
	{"%.5q", tftc1arg10alphabet, `"` + tftc1arg10alphabet[:5] + `"`},
	{"%.5x", tftc1arg10alphabet, "6162636465"},
	{"%.3q", tftc1arg11utf8cjkstr, `"` + tftc1arg11utf8cjkstr[:3*3] + `"`},
	{"%.1q", tftc1arg11utf8cjkshort, `"` + tftc1arg11utf8cjkstr[:1*3] + `"`},
	{"%.1x", tftc1arg11utf8cjkshort, "e6"},
	{"%.1X", tftc1arg11utf8cjkshort, "E6"},
	{"%10.1q", tftc1arg11utf8cjkstr, `       "` + tftc1arg11utf8cjkstr[:1*3] + `"`},

	// old tests
	{"%20.5s", tftc1arg12str, "               " + tftc1arg12str[:5]},
	{"%.5s", tftc1arg12str, tftc1arg12str[:5]},
	{"%-20.5s", tftc1arg12str, tftc1arg12str[:5] + "               "},

	// Padding with strings
	{"%2x", "", "  "},
	{"%#2x", "", "  "},
	{"% 02x", "", "00"},
	{"%# 02x", "", "00"},
	{"%-2x", "", "  "},
	{"%-02x", "", "  "},
	{"%8x", "\xab", "      ab"},
	{"% 8x", "\xab", "      ab"},
	{"%#8x", "\xab", "    0xab"},
	{"%# 8x", "\xab", "    0xab"},
	{"%08x", "\xab", "000000ab"},
	{"% 08x", "\xab", "000000ab"},
	{"%#08x", "\xab", "00000xab"},
	{"%# 08x", "\xab", "00000xab"},
	{"%10x", "\xab\xcd", "      abcd"},
	{"% 10x", "\xab\xcd", "     ab cd"},
	{"%#10x", "\xab\xcd", "    0xabcd"},
	{"%# 10x", "\xab\xcd", " 0xab 0xcd"},
	{"%010x", "\xab\xcd", "000000abcd"},
	{"% 010x", "\xab\xcd", "00000ab cd"},
	{"%#010x", "\xab\xcd", "00000xabcd"},
	{"%# 010x", "\xab\xcd", "00xab 0xcd"},
	{"%-10X", "\xab", "AB        "},
	{"% -010X", "\xab", "AB        "},
	{"%#-10X", "\xab\xcd", "0XABCD    "},
	{"%# -010X", "\xab\xcd", "0XAB 0XCD "},

	// erroneous things
	{"", tftc1arg3short, "%!(EXTRA string=" + tftc1arg3short + ")"},
	{"no args", tftc1arg13str, "no args%!(EXTRA string=" + tftc1arg13str + ")"},
	{"%s %", tftc1arg13str, tftc1arg13str + " %!(NOVERB)"},
	{"%s %.2", tftc1arg13str, tftc1arg13str + " %!(NOVERB)"},
	{"%017086757745859969700584773889596955657833456770", tftc1arg3short, "%!(NOVERB)%!(EXTRA string=" + tftc1arg3short + ")"},
	{"%184467440737095516170s", tftc1arg3short, "%!(NOVERB)%!(EXTRA string=" + tftc1arg3short + ")"},
	// Extra argument errors should format without flags set.
	{"%010.2", "12345", "%!(NOVERB)%!(EXTRA string=12345)"},

	// Use spaces instead of zero if padding to the right.
	{"%0-7s", tftc1arg3short, tftc1arg3short + "  "},

	// Tests to check that not supported verbs generate an error string.
	{"%★", tftc1arg13str, "%!★(string=" + tftc1arg13str + ")"},
}

// go test -count=1 -v -run "^TestAdaptedFmtTestCases$"
// go test -count=1 -o fmt.exe -gcflags "-m -m -d=ssa/check_bce/debug=1" -v -run "^TestAdaptedFmtTestCases$" 2> fmt.log
// go tool objdump -S -s "go-xfmt" fmt.exe > fmt.disasm
func TestAdaptedFmtTestCases(t *testing.T) {

	nerrs, sz := 0, len(adaptedFmtTestCases1)

	for i := 0; i < sz; i++ {

		tcase := &adaptedFmtTestCases1[i]

		if want, got := tcase.result, tcase.sprintf(); want != got {
			t.Errorf("adapted test case %d Sprintf(%q, %q) mismatch: want <%s>, got <%s> (xfmt: %#v)",
				i, tcase.format, tcase.arg, want, got, parseFormat(tcase.format))

			nerrs++
		}
	}

	if nerrs > 0 {
		t.Fatalf("some tests finished with errors: %d of %d", nerrs, sz)
	}
}

type SL []string // strings list

type fmtSLTestCase struct {
	format string
	args   SL
	result string
}

func (c *fmtSLTestCase) sprintf() string {
	return Sprintf(c.format, c.args...)
}

type fmtSLTestCases []fmtSLTestCase

func (cases fmtSLTestCases) run(t *testing.T, test string) {

	nerrs, sz := 0, len(cases)

	for i := 0; i < sz; i++ {

		tcase := &cases[i]

		if want, got := tcase.result, tcase.sprintf(); want != got {
			t.Errorf("adapted %s test case %d Sprintf(%q, %#v) mismatch: want <%s>, got <%s> (xfmt: %#v)",
				test, i, tcase.format, tcase.args, want, got, parseFormat(tcase.format))

			nerrs++
		}
	}

	if nerrs > 0 {
		t.Fatalf("some tests finished with errors: %d of %d", nerrs, sz)
	}
}

var adaptedFmtReorderTestCases1 = [...]fmtSLTestCase{
	{"%[1]s", SL{"1"}, "1"},
	{"%[2]s", SL{"2", "1"}, "1"},
	{"%[2]s %[1]s", SL{"1", "2"}, "2 1"},
	// An actual use! Print the same arguments twice.
	{"%s %s %s %#[1]q %#q %#q", SL{"13", "14", "15"}, "13 14 15 `13` `14` `15`"},

	// Erroneous cases.
	{"%[2]*[1]s", SL{"2", "5"}, "%!(BADWIDTH)2"},
	{"%[3]*.[2]*[1]s", SL{"13.0", "2", "6"}, "%!(BADWIDTH)%!(BADPREC)13.0"},
	{"%[1]*.[2]*[3]s", SL{"6", "2", "13.0"}, "%!(BADWIDTH)%!(BADPREC)13.0"},
	{"%[1]*[3]s", SL{"10", "99", "13.0"}, "%!(BADWIDTH)13.0"},
	{"%.[1]*[3]s", SL{"6", "99", "13.0"}, "%!(BADPREC)13.0"},
	{"%[1]*.[3]s", SL{"6", "3", "13.0"}, "%!(BADWIDTH)"},
	//
	{"%[s", SL{"2", "1"}, "%!s(BADINDEX)%!(EXTRA string=1)"}, // WARN mismatch with fmt package
	{"%]s", SL{"2", "1"}, "%!](string=2)s%!(EXTRA string=1)"},
	{"%[]s", SL{"2", "1"}, "%!s(BADINDEX)%!(EXTRA string=1)"},   // WARN mismatch with fmt package
	{"%[-3]s", SL{"2", "1"}, "%!s(BADINDEX)%!(EXTRA string=1)"}, // WARN mismatch with fmt package
	{"%[99]s", SL{"2", "1"}, "%!s(MISSING)"},                    // WARN mismatch with fmt package
	{"%[3]", SL{"2", "1"}, "%!(NOVERB)"},
	{"%[1].2s", SL{"5", "6"}, "%!s(BADINDEX)%!(EXTRA string=6)"}, // WARN mismatch with fmt package
	{"%[1]2s", SL{"2", "1"}, "%!s(BADINDEX)%!(EXTRA string=1)"},  // WARN mismatch with fmt package
	{"%3.[2]s", SL{"7"}, "%!s(MISSING)"},                         // WARN mismatch with fmt package
	{"%.[2]s", SL{"7"}, "%!s(MISSING)"},                          // WARN mismatch with fmt package
	{"%s %s %s %#[1]q %#q %#q %#q", SL{"13", "14", "15"}, "13 14 15 `13` `14` `15` %!q(MISSING)"},
	{"%[5]s %[2]s %s", SL{"1", "2", "3"}, "%!s(MISSING) 2 3"}, // WARN mismatch with fmt package
	{"%s %[3]s %s", SL{"1", "2"}, "1 %!s(BADINDEX) 2"},        // TODO mismatch with fmt package: erroneous index DOES affect sequence (but must doesn't)
	{"%.[]", SL{}, "%!](BADINDEX)"},
	{"%.-3s", SL{"zyx"}, "%!-(string=zyx)3s"},
	{"%2147483648s", SL{"zyx"}, "%!(NOVERB)%!(EXTRA string=zyx)"},
	{"%-2147483648s", SL{"zyx"}, "%!(NOVERB)%!(EXTRA string=zyx)"},
	{"%.2147483648s", SL{"zyx"}, "%!(NOVERB)%!(EXTRA string=zyx)"},
}

// go test -count=1 -v -run "^TestAdaptedFmtReorderTestCases$"
// go test -count=1 -o fmt.exe -gcflags "-m -m -d=ssa/check_bce/debug=1" -v -run "^TestAdaptedFmtReorderTestCases$" 2> fmt.log
// go tool objdump -S -s "go-xfmt" fmt.exe > fmt.disasm
func TestAdaptedFmtReorderTestCases(t *testing.T) {
	fmtSLTestCases(adaptedFmtReorderTestCases1[:]).run(t, "reorder")
}

// benchmarks
// go test -bench=. -run "^$" -benchmem
// go test -bench=. -run "^$" -benchmem -cpuprofile cpu.pprof -memprofile mem.pprof

func doBenchmarkSprintf(b *testing.B, format string, args ...string) {

	var iargs []interface{}

	if len(args) > 0 {
		iargs = make([]interface{}, len(args))
	}

	purgeCaches()

	runtime.GC()

	b.Run("xfmt", func(bb *testing.B) {
		bb.ResetTimer()
		bb.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = Sprintf(format, args...)
			}
		})
	})

	//b.Logf("%#v", bpool)

	runtime.GC()

	b.Run("fmt", func(bb *testing.B) {
		bb.ResetTimer()
		bb.RunParallel(func(pb *testing.PB) {

			for pb.Next() {

				// NOTE in real scenarios `string` always cast to `interface{}` therefore it
				//      must be taken into account in the benchs

				if len(args) > 0 {
					for i, s := range args {
						iargs[i] = s
					}
				}

				_ = fmt.Sprintf(format, iargs...)
			}
		})
	})
}

// go test -bench "^BenchmarkSprintfPadding$" -run "^$"
func BenchmarkSprintfPadding(b *testing.B) {
	doBenchmarkSprintf(b, "%16s", tftc1arg5utf8MS_HnS)
}

// go test -bench "^BenchmarkSprintfEmpty$" -run "^$"
func BenchmarkSprintfEmpty(b *testing.B) {
	doBenchmarkSprintf(b, "")
}

// go test -bench "^BenchmarkSprintfString$" -run "^$"
func BenchmarkSprintfString(b *testing.B) {
	doBenchmarkSprintf(b, "%s", tftc1arg13str)
}

// go test -bench "^BenchmarkSprintfTruncateString$" -run "^$"
func BenchmarkSprintfTruncateString(b *testing.B) {
	doBenchmarkSprintf(b, "%.4s", tftc1arg14utf8cjklong)
}

// go test -bench "^BenchmarkSprintfQuoteString$" -run "^$"
func BenchmarkSprintfQuoteString(b *testing.B) {
	doBenchmarkSprintf(b, "%q", tftc1arg14utf8cjklong)
}

// go test -bench "^BenchmarkSprintfPrefixedString$" -run "^$"
func BenchmarkSprintfPrefixedString(b *testing.B) {
	doBenchmarkSprintf(b, "This is some meaningless prefix text that needs to be scanned %s", tftc1arg5utf8MS_HnS)
}

// go test -bench "^BenchmarkSprintfHexString$" -run "^$"
func BenchmarkSprintfHexString(b *testing.B) {
	doBenchmarkSprintf(b, "% #x", tftc1arg15str)
}

// go test -bench "^BenchmarkManyArgs$" -run "^$"
func BenchmarkManyArgs(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var buf bytes.Buffer
		for pb.Next() {
			buf.Reset()
			_, _ = Fprintf(&buf, "%2s/%2s/%2s %s:%s:%s %s %s\n", "3", "4", "5", "11", "12", "13", "hello", "world")
		}
	})
}

// malloc tests
//

var mallocBuf bytes.Buffer

type mallocTestCase struct {
	desc   string
	allocs uint
	fn     func()
}

var mallocTestCases1 = [...]mallocTestCase{
	{`Sprintf("")`, 0, func() { Sprintf(emptyString) }},
	{`Sprintf("xfmt")`, 0, func() { Sprintf("xfmt") }},
	{`Sprintf("%x")`, 1, func() { Sprintf("%x", tftc1arg11utf8cjkstr) }},
	{`Sprintf("%s")`, 1, func() { Sprintf("%s", tftc1arg13str) }},
	{`Sprintf("%x %x")`, 1, func() { Sprintf("%x %x", tftc1arg11utf8cjkstr, tftc1arg1strrnd) }},
	{`Sprintf("%q")`, 1, func() { Sprintf("%q", tftc1arg12str) }},
	{`Fprintf(buf, "%s")`, 0, func() {
		mallocBuf.Reset()
		Fprintf(&mallocBuf, "%s", tftc1arg13str)
	}},
	{`Fprintf(buf, "%x %x %x")`, 0, func() {
		mallocBuf.Reset()
		Fprintf(&mallocBuf, "%x %x %x", tftc1arg13str, tftc1arg15str, tftc1arg3short)
	}},
}

// go test -count=1 -v -run "^TestAssertMallocs$"
func TestAssertMallocs(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping malloc count in short mode")
	}

	for i := 0; i < len(mallocTestCases1); i++ {

		tcase := &mallocTestCases1[i]

		mallocs := testing.AllocsPerRun(100, tcase.fn)

		if got, max := mallocs, float64(tcase.allocs); got > max {
			t.Errorf("%s: got %v allocs, want <= %v", tcase.desc, got, max)
		}
	}
}

// indir width and prec tests (startests) - should check for errors, no indir values are processing

var adaptedFmtIndirTestCases1 = [...]fmtSLTestCase{
	{"%*s", SL{"7", "13"}, "%!(BADWIDTH)13"},
	{"%-*s", SL{"7", "13"}, "%!(BADWIDTH)13"},
	{"%*s", SL{"-7", "13"}, "%!(BADWIDTH)13"},
	{"%-*s", SL{"-7", "13"}, "%!(BADWIDTH)13"},
	{"%.*s", SL{"7", "13"}, "%!(BADPREC)13"},
	{"%*.*s", SL{"11", "3", "13"}, "%!(BADWIDTH)%!(BADPREC)13"},
	{"%0*s", SL{"7", "13"}, "%!(BADWIDTH)13"},
	{"%0*s", SL{"0x07", "13"}, "%!(BADWIDTH)13"},

	// erroneous
	{"%*s", SL{"", "13"}, "%!(BADWIDTH)13"},
	{"%*s", SL{"1048577", "13"}, "%!(BADWIDTH)13"},
	{"%*s", SL{"-1048577", "13"}, "%!(BADWIDTH)13"},
	{"%.*s", SL{"", "13"}, "%!(BADPREC)13"},
	{"%.*s", SL{"-1", "13"}, "%!(BADPREC)13"},
	{"%.*s", SL{"1048577", "13"}, "%!(BADPREC)13"},
	{"%.*s", SL{"9223372036854775808", "13"}, "%!(BADPREC)13"},  // Huge negative (-inf).
	{"%.*s", SL{"18446744073709551615", "13"}, "%!(BADPREC)13"}, // Small negative (-1).
	{"%*% %s", SL{"17", "7"}, "% 7"},
	{"%*", SL{"7"}, "%!(NOVERB)"},
	{"%*", SL{"7", "xfmt"}, "%!(NOVERB)%!(EXTRA string=xfmt)"},
}

// go test -count=1 -v -run "^TestWidthAndPrecision$"
func TestWidthAndPrecision(t *testing.T) {
	fmtSLTestCases(adaptedFmtIndirTestCases1[:]).run(t, "indir")
}
