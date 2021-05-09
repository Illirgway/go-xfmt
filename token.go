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
	"github.com/valyala/bytebufferpool"
	"strconv"
	"unicode/utf8"
)

/* DOC: String and slice of bytes (treated equivalently with these verbs):
 *
 *	%%	a literal percent sign; consumes no value
 *
 *	%s	the uninterpreted bytes of the string or slice
 *	%q	a double-quoted string safely escaped with Go syntax
 *	%x	base 16, lower-case, two characters per byte
 *	%X	base 16, upper-case, two characters per byte
 *
 */

type verb uint

const (
	verbCharString = 's'
	verbCharQuoted = 'q'
	verbCharHex    = 'x'
)

const (
	verbNone   verb = iota // direct plain raw string const value
	verbString             // %s
	verbQuoted             // %q
	verbHex                // %x or %X

	// bad verb
	badVerb

	maxVerbs
)

// TODO? [26]verb instead of map[rune]verb (26 is count of en letters --> F(v - 'a'))?
var verbMap = map[rune]verb{
	verbCharString: verbString,
	verbCharQuoted: verbQuoted,
	verbCharHex:    verbHex,
}

/* DOC: format flags:
 *
 *	+	always print a sign for numeric values;
 *		guarantee ASCII-only output for %q (%+q)
 *	-	pad with spaces on the right rather than the left (left-justify the field)
 *	#	alternate format: add leading 0b for binary (%#b), 0 for octal (%#o),
 *		0x or 0X for hex (%#x or %#X); suppress 0x for %p (%#p);
 *		for %q, print a raw (backquoted) string if strconv.CanBackquote
 *		returns true;
 *		always print a decimal point for %e, %E, %f, %F, %g and %G;
 *		do not remove trailing zeros for %g and %G;
 *		write e.g. U+0078 'x' if the character is printable for %U (%#U).
 *	' '	(space) leave a space for elided sign in numbers (% d);
 *		put spaces between bytes printing strings or slices in hex (% x, % X)
 *	0	pad with leading zeros rather than spaces;
 *		for numbers, this moves the padding after the sign
 *
 * Flags are ignored by verbs that do not expect them. For example there is no alternate decimal format, so %#d
 * and %d behave identically.
 */

type flags uint

const (
	flagCharPlus  = '+'
	flagCharMinus = '-'
	flagCharSharp = '#'
	flagCharSpace = charSpace
	flagCharZero  = charZero
)

const (
	flagPlus flags = 1 << iota
	flagMinus
	flagSharp
	flagSpace
	flagZero

	// internal flags

	// - verb in uppercase
	flagUpperVerb
	// - width is indirect ("[n]*")
	flagIndirectWidth
	// - precision is indirect (".[n]")
	flagIndirectPrec

	flagsMask flags = (1 << iota) - 1

	flagNone flags = 0 // goes after all iota-depended consts to not to affect iota itself

	// aliases with meaningful names
	flagSign      = flagPlus
	flagAsciiOnly = flagPlus
	flagPadRight  = flagMinus
	flagWithSpace = flagSpace
	flagAltFmt    = flagSharp
	flagPadZeros  = flagZero
)

// TODO
var charFlags = map[byte]flags{}

// inlined
//go:nosplit
func (flags flags) has(flag flags) bool {
	return (flags & flag) != 0
}

// eq. !has
// inlined
//go:nosplit
func (flags flags) omit(flag flags) bool {
	return (flags & flag) == 0
}

const (
	maxNumBitSize = 20 // required number of bits to store src/fmt/print.go::tooLarge()::max === 1e6
	maxNum        = 1<<maxNumBitSize - 1
)

// DOC:
//
// Width is specified by an optional decimal number immediately preceding the verb. If absent, the width is
// whatever is necessary to represent the value. Precision is specified after the (optional) width by a period followed
// by a decimal number. If no period is present, a default precision is used. A period with no following number
// specifies a precision of zero
//
// Width and precision are measured in units of Unicode code points, that is, runes. (This differs from C's printf
// where the units are always measured in bytes.) Either or both of the flags may be replaced with the character '*',
// causing their values to be obtained from the next operand (preceding the one to format), which must be of type int.
//
// For most values, width is the minimum number of runes to output, padding the formatted form with spaces if necessary.
//
// For strings, byte slices and byte arrays, however, precision limits the length of the input to be formatted (not the
// size of the output), truncating if necessary. Normally it is measured in runes, but for these types when formatted
// with the %x or %X format it is measured in bytes.
//

// DOC:
//
// In Printf, Sprintf, and Fprintf, the default behavior is for each formatting verb to format successive arguments
// passed in the call. However, the notation [n] immediately before the verb indicates that the nth one-indexed
// argument is to be formatted instead. The same notation before a '*' for a width or precision selects the argument
// index holding the value. After processing a bracketed expression [n], subsequent verbs will use arguments n+1, n+2,
// etc. unless otherwise directed.
//

/* Putting all docs together (lax ebnf):
 * FormatSpec = '%' Flag* WidthPrec? ArgNum? Verb
 * Flag = ASCIICHAR
 * WidthPrec = Width? Precision?
 * Width = DirectValue | IndirectValue
 * Precision = '.' ( DirectValue | IndirectValue )
 * DirectValue = UINT
 * IndirectValue = ArgNum? '*'
 * Verb = ASCIICHAR
 * ArgNum = '[' UINT ']'
 */

type token struct {
	verb  verb
	value string
	flags flags
	width int
	prec  int  // precision
	arg   uint // arg number, handle notation [n] immediately before the verb; uint automagically helps BCE
}

// pseudoflag `no value defined` for width and prec
const absentValue = -1

type fmtFn = func(buf *bytebufferpool.ByteBuffer /* *strings.Builder */, value string, flags flags, width, prec int) (success bool)

// TODO??
/*
type fmtFuncsTable [maxVerbs]fmtFn

//go:nosplit
func (t *fmtFuncsTable) fn(v verb) fmtFn {
	if (int(v) < len(*t)) {
		// may return nil
		return (*t)[v]
	}

	return nil
}
*/

var verbFmtFuncsTable = [...]fmtFn{ // fmtFuncsTable{
	verbString: fmtStr,
	verbQuoted: fmtQuot,
	verbHex:    fmtHex,
}

// inlined
//go:nosplit
func fmtFuncByVerb(v verb) fmtFn {

	if uint(v) < uint(len(verbFmtFuncsTable)) { // with BCE opt.
		// may return nil
		return verbFmtFuncsTable[v]
	}

	return nil
}

// special token cases
var (
	tokenPercent = token{
		verb:  verbNone,
		value: percentString,
	}

	tokenErrNoVerb = token{
		verb:  verbNone,
		value: noVerbString,
	}

	tokenErrBadPrec = token{
		verb:  verbNone,
		value: badPrecString,
	}
)

// buffer pool for temporary strings
var fmtBufPool bytebufferpool.Pool

// TODO (*strings.Builder).WriteString() не имеет внутренней встроенной проверки на empty string noop `s == ""`,
//      когда ничего просто делать не надо кроме return, и для ее внедрения необходим враппер вокруг strings.Builder
//      Также нужна встроенная функция buf.Pad чтобы делать паддинг

// SEE src/fmt/print.go::(*pp).badArgNum()
//go:nosplit
func (token *token) badArgNum(buf *bytebufferpool.ByteBuffer /* *strings.Builder */) {

	// PPSL: use Grow before multiple writes to minimize memallocs

	buf.WriteString(percentBangString)
	buf.WriteString(token.value)
	buf.WriteString(badIndexString)
}

// SEE src/fmt/print.go::(*pp).missingArg()
//go:nosplit
func (token *token) missingArg(buf *bytebufferpool.ByteBuffer /* *strings.Builder */) {

	// PPSL: use Grow before multiple writes to minimize memallocs

	buf.WriteString(percentBangString)
	buf.WriteString(token.value)
	buf.WriteString(missingString)
}

// SEE src/fmt/print.go::(*pp).badVerb()
//go:nosplit
func (token *token) badVerb(buf *bytebufferpool.ByteBuffer /* *strings.Builder */, arg string) {

	// PPSL: use Grow before multiple writes to minimize memallocs

	buf.WriteString(percentBangString)
	buf.WriteString(token.value)

	buf.WriteByte(charLeftParens) // (

	buf.WriteString(reflectStringType)
	buf.WriteByte(charEquals) // =
	buf.WriteString(arg)

	buf.WriteByte(charRightParens) // )
}

// SEE src/fmt/print.go::(*pp).fmtString()
// such code leads to `moved to heap: buf` at least in go1.13
// `go1.13: cannot inline (*token).fmtStringE2H: function too complex: cost 152 exceeds budget 80`
//go:nosplit
func (token *token) fmtStringE2H(buf *bytebufferpool.ByteBuffer /* *strings.Builder */, arg string, flags flags, width, prec int) bool {

	if fn := fmtFuncByVerb(token.verb); fn != nil {
		return fn(buf, arg, flags, width, prec)
	}

	// impossible situation
	token.badVerb(buf, arg)

	return false
}

// SEE src/fmt/print.go::(*pp).fmtString()
// NOTE long version without `buf` e2h
//go:nosplit
func (token *token) fmtString(buf *bytebufferpool.ByteBuffer /* *strings.Builder */, arg string, flags flags, width, prec int) bool {

	switch token.verb {
	case verbString:
		return fmtStr(buf, arg, flags, width, prec)
	case verbQuoted:
		return fmtQuot(buf, arg, flags, width, prec)
	case verbHex:
		return fmtHex(buf, arg, flags, width, prec)
	}

	// impossible situation
	token.badVerb(buf, arg)

	return false
}

// @return bool success status
func (token *token) format(buf *bytebufferpool.ByteBuffer /* *strings.Builder */, args []string) bool {

	// impossible situation
	if token == nil {
		//buf.WriteString(nilToken)
		return false
	}

	// simple case - solely raw string const value
	if token.verb == verbNone {
		buf.WriteString(token.value)
		return true
	}

	// complex cases

	// similar to original fmt algo, should first check the width and prec, and only then the verb

	// here width and prec are proper values except for indirect which may be missingArg

	flags := token.flags // clone flags for possible inplace modifications

	// second part of the original `goodArgNum` check
	badArgNum := false

	width := token.width

	// is width indirect?
	if flags.has(flagIndirectWidth) {
		// check is width arg exists
		if width >= len(args) {
			badArgNum = true
			buf.WriteString(badWidthString)
			// absent width arg, must reset width value ...
		} else {
			// to repeating the algorithm of the original `fmt` package, `string` arg can't be used as width int
			// SEE src/fmt/print.go::intFromArg() ...
			buf.WriteString(badWidthString)
		}

		// ... so in any case for indir width should mark as absent value
		width = absentValue
	}

	prec := token.prec

	// is precision indirect?
	if flags.has(flagIndirectPrec) {
		// check is width arg exists
		if prec >= len(args) {
			badArgNum = true
			buf.WriteString(badPrecString)
			// absent prec arg, must reset prec value ...
		} else {

			// to repeating the algorithm of the original `fmt` package, `string` arg can't be used as prec int
			// SEE src/fmt/print.go::intFromArg() ...
			buf.WriteString(badPrecString)
		}

		// ... so in any case for indir prec should mark as absent prec
		prec = absentValue
	}

	// if arg num of either width or prec out of bounds
	if badArgNum {
		token.badArgNum(buf)
		return false
	}

	// here width and prec are proper values

	// check that appropriate arg exists in args
	if token.arg >= uint(len(args)) {
		token.missingArg(buf)
		return false
	}

	// pick our arg
	arg := args[token.arg]

	// first check cases which mean `error`
	if token.verb == badVerb {
		token.badVerb(buf, arg)
		return false
	}

	// here token is known, arg exists

	// here verb is guaranteed to exist and be known (unknown `badVerb` filtered above)
	return token.fmtString(buf, arg, flags, width, prec)
}

//

// SEE src/fmt/format.go::(*fmt).fmtS()
//go:nosplit
func fmtStr(buf *bytebufferpool.ByteBuffer /* *strings.Builder */, s string, flags flags, width, prec int) bool {
	s = truncateString(s, prec)
	padString(buf, s, flags, width)
	return true
}

// SEE src/fmt/format.go::(*fmt).fmtQ()
//go:nosplit
func fmtQuot(buf *bytebufferpool.ByteBuffer /* *strings.Builder */, s string, flags flags, width, prec int) bool {

	s = truncateString(s, prec)

	if flags.has(flagAltFmt) && strconv.CanBackquote(s) {

		// `go1.13: fmtQuot backquoteStr + s + backquoteStr does not escape`
		padString(buf, backquoteStr+s+backquoteStr /* concatstring3 with tempBuf */, flags, width)

		// NOTE concatstring3 with tempBuf use stack temp buffer with size 32 if enough or memallocated heap tmp buf if not

		return true
	}

	b := fmtBufPool.Get()

	// hate defer
	//defer fmtBufPool.Put(b)

	// check ascii-only flag
	if flags.has(flagAsciiOnly) {
		b.B = strconv.AppendQuoteToASCII(b.B, s)
	} else {
		b.B = strconv.AppendQuote(b.B, s)
	}

	pad(buf, b.B, flags, width)

	// w/o defer
	fmtBufPool.Put(b)

	return true
}

const (
	ldigits = "0123456789abcdefx"
	udigits = "0123456789ABCDEFX"
)

// SEE src/fmt/format.go::(*fmt).fmtSx() -> (*fmt).fmtSbx()
//go:nosplit
func fmtHex(buf *bytebufferpool.ByteBuffer /* *strings.Builder */, s string, flags flags, width, prec int) bool {

	sz := len(s)

	// Set length to not process more bytes than the precision demands.
	if (prec != absentValue) && (prec < sz) {
		sz = prec
	}

	// If string that should be encoded is empty.
	if sz <= 0 {
		// ... should write padding if any
		if width > 0 /* implies `width != absentValue` */ {
			writePadding(buf, width, flags)
		}

		return true
	}

	// here sz > 0

	withSpace, altFmt := flags.has(flagWithSpace), flags.has(flagAltFmt)

	// Compute width of the encoding taking into account the flagAltFmt and flagWithSpace flag.
	w := 2 * sz

	if withSpace {
		// Each element encoded by two hexadecimals will get a leading 0x or 0X.
		if altFmt {
			w *= 2
		}

		w += sz - 1
	} else if altFmt {
		// Only a leading 0x or 0X will be added for the whole string.
		w += 2
	}

	rpad := flags.has(flagPadRight)

	// Handle padding to the left.
	if (width != absentValue) && !rpad && (width > w) {
		writePadding(buf, width-w, flags)
	}

	// select digits set - lowercased or uppercased

	digits := ldigits

	if flags.has(flagUpperVerb) {
		digits = udigits
	}

	// TODO Write the encoding directly into the output buffer.
	//      Additionally *strings.Builder impl has too many redundant checks (every buf.Write has copyCheck and
	//      implicit growslice check in append)

	/* TODO
	// prepare buffer to minimize memallocs
	buf.Grow(w)
	*/

	if altFmt {
		// Add leading 0x or 0X.
		buf.WriteByte(charZero)
		buf.WriteByte(digits[16])
	}

	// helps BCE below
	s = s[:sz]

	for i := 0; i < len(s); i++ {

		// write inter-elements values if any
		if withSpace && (i > 0) {

			// Separate elements with a space.
			buf.WriteByte(charSpace)

			if altFmt {
				// Add leading 0x or 0X for each element.
				buf.WriteByte(charZero)
				buf.WriteByte(digits[16])
			}
		}

		// write byte as hex
		c := s[i]

		buf.WriteByte(digits[c>>4])   // high
		buf.WriteByte(digits[c&0x0F]) // low
	}

	// Handle padding to the right.
	if (width != absentValue) && rpad && (width > w) {
		writePadding(buf, width-w, flags)
	}

	return true
}

// truncateString truncates the string s to the specified precision, if present.
// DOC: For strings, byte slices and byte arrays, however, precision limits the length of the input to be formatted,
//      truncating if necessary. Normally it is measured in runes {utf8 codepoints}
// SEE src/fmt/format.go::(*fmt).truncateString()
// inlined
//go:nosplit
func truncateString(s string, prec int) string {

	if prec != absentValue {
		s = truncateTail(s, prec)
	}

	return s
}

// padString appends s to buf, padded on left (!flagPadRight) or right (flagPadRight).
// DOC: For most values, width is the minimum number of runes to output, padding the formatted form with spaces if necessary.
// SEE src/fmt/format.go::(*fmt).padString()
//go:nosplit
func padString(buf *bytebufferpool.ByteBuffer /* *strings.Builder */, s string, flags flags, width int) {

	if width <= 0 /* implies `width == absentValue` */ {
		buf.WriteString(s)
		return
	}

	// DOC: width := f.wid - utf8.RuneCountInString(s)
	width -= utf8.RuneCountInString(s)

	rpad := flags.has(flagPadRight)

	// either left padding ...
	if !rpad && (width > 0) {
		writePadding(buf, width, flags)
	}

	buf.WriteString(s)

	// ... or right padding
	if rpad && (width > 0) {
		writePadding(buf, width, flags)
	}
}

// pad appends b to f.buf, padded on left (!f.minus) or right (f.minus).
// DOC: see padString above
//go:nosplit
func pad(buf *bytebufferpool.ByteBuffer /* *strings.Builder */, raw []byte, flags flags, width int) {

	if width <= 0 /* implies `width == absentValue` */ {
		buf.Write(raw)
		return
	}

	// DOC: width := f.wid - utf8.RuneCount(s)
	width -= utf8.RuneCount(raw)

	rpad := flags.has(flagPadRight)

	// either left padding ...
	if !rpad && (width > 0) {
		writePadding(buf, width, flags)
	}

	buf.Write(raw)

	// ... or right padding
	if rpad && (width > 0) {
		writePadding(buf, width, flags)
	}
}

//go:nosplit
func writePadding(buf *bytebufferpool.ByteBuffer /* *strings.Builder */, n int, flags flags) {

	if n <= 0 { // No padding bytes needed.
		return
	}

	// Decide which byte the padding should be filled with.
	padByte := byte(charSpace)

	if flags.has(flagPadZeros) {
		padByte = charZero
	}

	/* TODO
	// Make enough room for padding.
	newLen := buf.Len() + n
	buf.Grow(newLen)
	*/

	// WARN this is very BAD code even taking into account the Grow() above,
	//      need buf.Pad for max efficiency
	for ; n > 0; n-- {
		buf.WriteByte(padByte)
	}
}
