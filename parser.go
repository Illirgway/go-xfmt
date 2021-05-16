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
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// NOTE return nil `xfmt.tokens` slice for empty format string
// NOTE retval by value
func parseFormat(format string) xfmt {

	needArgs, curArg, minSize := 0, 0, 0

	var tokens []token

	if format != "" {
		// use `make` with cap adjusted by counting '%', excluding double counting for "%%" (approximate count)
		if tokApxCount := strings.Count(format, percentStr) - strings.Count(format, doublePercentStr); tokApxCount > 0 {
			tokens = make([]token, 0, tokApxCount)
		}
	}

	// fsm begins
parseLoop:
	for format != "" /* implies `len(format) > 0` */ {

		// fast-check
		if format[0] != charPercent {

			// here format[0] !== '%'

			// NOTE arch-dep IndexByte (bytealg.IndexByteString) may be faster than simple loop
			i := strings.IndexByte(format, charPercent)

			// if '%' not found - grab all existing tail
			if i == -1 {
				i = len(format)
			}

			// store raw const direct string part into tokens
			tokens = append(tokens, token{
				verb:  verbNone,
				value: format[:i],
			})

			// count min size by value len of raw const string token
			minSize += i

			// done processing format string?
			if i == len(format) { /* xor we can check `format == ""` / `len(format) == 0` after reslicing below */
				break
			}

			// reslice format to '%' char
			format = format[i:]
		}

		// here format[0] === '%'

		flags := flagNone

		// ATN! starts from next char just after '%'
		i := uint(1) // uint automagically helps BCE optimization without additional conds

		// fast-path simple check
	fastLoop:
		for ; i < uint(len(format)); i++ {

			c := format[i]

			// here i < len(format), so we can use i+1 as format' end slice index (`format[i+1:]`)

			switch c {

			// special case - percent char ("%%" case)
			case charPercent:

				// DOC: Percent does not absorb operands and ignores f.wid and f.prec.

				// emulate percent as no verb raw const string value "%"
				// TODO? if prev token is verbNone const raw string value, it can be simply merged (append) with '%' char
				tokens = append(tokens, tokenPercent)

				// count min size
				minSize += len(percentString)

				// slice format forward after cur '%' char
				format = format[i+1:]

				continue parseLoop

			case flagCharPlus:
				flags |= flagPlus
			case flagCharSharp:
				flags |= flagSharp
			case flagCharMinus:
				flags |= flagMinus
				// DOC: Do not pad with zeros to the right.
				flags &^= flagZero // exclude zero flag
			case flagCharSpace:
				flags |= flagSpace
			case flagCharZero:
				// DOC: Only allow zero padding to the left.
				//p.fmt.zero = !p.fmt.minus
				if flags&flagMinus == 0 {
					flags |= flagZero
				}
			default:
				// here i < len(format)

				// Fast path for common case of ascii simple verbs
				// without precision or width or argument indices.
				// TODO? see src/strconv/atoi.go::lower()
				if ('a' <= c) && (c <= 'z') || ('A' <= c) && (c <= 'Z') {

					verb, isUpper := char2verb(rune(c))

					if isUpper {
						flags |= flagUpperVerb
					}

					// append token to tokens
					tokens = append(tokens, token{
						verb:  verb,
						value: format[i : i+1], // slice only verb without flags as token value
						flags: flags,
						width: absentValue,  // no width
						prec:  absentValue,  // no prec
						arg:   uint(curArg), // NOTE curArg always >= 0
					})

					// verb exists, should move curArg forward and check needArgs regardless of whether the verb is known
					curArg++

					// here curArg is next using arg num

					if curArg > needArgs {
						needArgs = curArg
					}

					// reslice format next to cur verb
					format = format[i+1:]

					continue parseLoop
				}

				// DOC: Format is more complex than simple flags and a verb or is malformed.
				break fastLoop
			}
		}

		// complex case, maybe have width and precision
		// NOTE ahead forward flags already processed above in fastLoop, so here we stay at
		// `WidthPrec` or `ArgNum` or format error (i.e utf-8 char or spec char), but not at allowed known `Verb`
		// because simple direct 'ascii char verb' has been processed in fastLoop

		var hasArgNum, properArgNum, properNextArgNum bool

		// all of Width, Precision and ArgNum cases may starts with ArgNum (see ebnf above). Try to parse ArgNum

		curArg, i, hasArgNum, properArgNum = tryArgNum(curArg, format, i)

		// here curArg is parsed cur arg num
		if hasArgNum && (curArg >= needArgs) {
			needArgs = curArg + 1
		}

		// Is there width?
		width := absentValue

		if (i < uint(len(format))) && (format[i] == charAsterisk) {

			// move to next char after asterisk
			i++

			// indirect width found
			width = curArg
			flags |= flagIndirectWidth

			// should inc used curArg
			curArg++

			// here curArg is next using arg num

			if curArg > needArgs {
				needArgs = curArg
			}

			// width consume argNum, note this for following code
			hasArgNum = false
		} else {
			// try to blind pick width num
			width, i = pickNumValue(format, i)
			// ... and check wrong fmt "%[3]2x"
			if hasArgNum && (width != absentValue) {
				properArgNum = false
			}
		}

		// Is there precision?
		prec := absentValue

		if (i < uint(len(format))) && (format[i] == charDot) {

			// NOTE no prec value (either nothing after '.' or no digit or not num) is legal prec value means 0

			/* DEPR

			// proceed only if have anything after '.'
			if (i + 1) >= uint(len(format)) { // WARN `>=` instead of `==`
				// '.' is last char - bad precision format + no verb
				// so it is the last token
				// emulate badPrec as verbNone with error string
				tokens = append(tokens, tokenErrBadPrec)

				break parseLoop
			}

			// here i+1 < sz

			*/

			// move to the next char after '.'
			i++

			// handle wrong fmt "%[3].2d"
			if hasArgNum {
				properArgNum = false
			}

			curArg, i, hasArgNum, properNextArgNum = tryArgNum(curArg, format, i)

			// here curArg is parsed cur arg num
			if hasArgNum && (curArg >= needArgs) {
				needArgs = curArg + 1
			}

			// result is true (cur token' argNums are proper) only if both are true
			properArgNum = properArgNum && properNextArgNum

			if (i < uint(len(format))) && (format[i] == charAsterisk) {

				// move to next char after asterisk
				i++

				prec = curArg
				flags |= flagIndirectPrec

				// should inc used curArg
				curArg++

				// here curArg is next using arg num

				if curArg > needArgs {
					needArgs = curArg
				}

				// prec consume argNum, note this for following code
				hasArgNum = false
			} else {
				// no asterisk, try to blind pick prec num
				prec, i = pickNumValue(format, i)

				// ".x" (no prec value, but has verb after dot) is legal prec means 0
				if prec == absentValue {
					prec = 0
				}
			}
		}

		// if already consume argNum, try once more
		if !hasArgNum {
			curArg, i, hasArgNum, properNextArgNum = tryArgNum(curArg, format, i)

			// result is true (cur token' argNums are proper) only if both are true
			properArgNum = properArgNum && properNextArgNum

			// needArgs check is below
		}

		// handle unfinished fmt with absent last verb
		if i >= uint(len(format)) {
			// it is the last token
			// emulate noVerb as verbNone with error string
			tokens = append(tokens, tokenErrNoVerb)

			break parseLoop
		}

		// has verb, now should adjust needArgs (but only if has arg num)
		// this emulate the same behavior
		//		if ok && 0 <= index && index < numArgs {
		//			return index, i + wid, true
		//		}
		//		p.goodArgNum = false
		// of
		// 		if !afterIndex {
		//			argNum, i, afterIndex = p.argNumber(argNum, format, i, len(a))
		//		}
		//
		if hasArgNum && (curArg >= needArgs) {
			needArgs = curArg + 1
		}

		verbChar, size := rune(format[i]), 1

		if verbChar >= utf8.RuneSelf {
			verbChar, size = utf8.DecodeRuneInString(format[i:])
		}

		var tok token

		switch {
		case verbChar == charPercent:
			// DOC: Percent does not absorb operands and ignores f.wid and f.prec.
			tok = tokenPercent

			// count min size
			minSize += len(percentString)

		case !properArgNum:

			// NOTE emulate bad arg num error as verbNone const raw string value with crafted error string
			tok = token{
				verb:  verbNone,
				value: badArgNumValue(verbChar),
				flags: flags,
				width: width,
				prec:  prec,
			}

			// count min size
			minSize += len(tok.value)

		default:
			verb, isUpper := char2verb(verbChar)

			if isUpper {
				flags |= flagUpperVerb
			}

			// append token to tokens
			tok = token{
				verb:  verb,
				value: format[i : i+uint(size)], // slice only verb as token value + don't use `string(verbChar)` to avoid memalloc
				flags: flags,
				width: width,
				prec:  prec,
				arg:   uint(curArg), // NOTE curArg always >= 0
			}

			// verb exists, should move curArg forward and check lastArg regardless of whether the verb is known
			curArg++

			// here curArg is next using arg num

			if curArg > needArgs {
				needArgs = curArg
			}
		}

		// append token to tokens
		tokens = append(tokens, tok)

		// this helps BCE below
		if (i + uint(size)) >= uint(len(format)) {
			break parseLoop
		}

		// reslice `format` next to processed token
		format = format[i+uint(size):]
	}
	// fsm ends

	// try to shrink unnecessary tokens slice cap
	const tokensOversizeThreshold = 1

	if len(tokens)+tokensOversizeThreshold < cap(tokens) {
		// compiler optimization CL 146719 make+copy pattern
		// https://go-review.googlesource.com/c/go/+/146719/
		tmp := make([]token, len(tokens))
		copy(tmp, tokens)

		tokens = tmp
	}

	return xfmt{
		tokens:  tokens,
		args:    uint(needArgs), // here needArgs cannot be less than 0
		minSize: minSize,
	}
}

//

const badArgNum = -1

// SEE src/fmt/print.go::(*pp).argNumber()
// NOTE `proper` is always true when `found` is false
//go:nosplit
func tryArgNum(argNum int, format string, i uint) (newArgNum int, j uint, found, proper bool) {

	if (i >= uint(len(format))) || (format[i] != charOpenArgNum) {
		return argNum, i, false, true
	}

	newArgNum, offset := parseArgNum(format[i:])

	if (newArgNum == badArgNum) || (newArgNum > maxNum) {
		return argNum, i + offset, true, false
	}

	return newArgNum, i + offset, true, true
}

// parseArgNum returns the value of the bracketed number, minus 1
// (explicit argument numbers are one-indexed but we want zero-indexed).
// The opening bracket is known to be present at format[0].
// The returned values are the index, the number of bytes to consume
// up to the closing paren, if present, and whether the number parsed
// ok. The bytes to consume will be 1 if no closing paren is present.
// NOTE retval badArgNum (-1) means 'not ok' ('not found' or 'wrong index 0' or 'num parse error')
//go:nosplit
func parseArgNum(format string) (num int, offset uint) {

	// DOC: There must be at least 3 bytes: [n].
	if len(format) < 3 {
		return badArgNum, 1
	}

	// skip open charOpenArgNum for BCE below, so must take this into account in return (1)
	format = format[1:]

	// try found closing bracket
	// NOTE here format start from next char after charOpenArgNum
	// ALGO uint(-1) == ^uint(0) == MaxUint - 1 --> max possible uint value guaranteed to be always greater
	//      than any possible `format` string len (because len is int and `uint(MaxInt - 1) < MaxUint - 1`)
	//      so the following trick with `uint(IndexByte)` can be combined with `i < len(format)` check to help BCE
	if i := uint(strings.IndexByte(format, charCloseArgNum)); i < uint(len(format)) {

		// here i is index of ']' inside sliced in (1) `format`, `i+1` is index of ']' inside initial `format` string
		// and `i+1+1` is index of first char of next token after ArgNum token inside initial `format` string

		v64, err := strconv.ParseUint(format[:i], 10, maxNumBitSize) // BCEd by above tricks

		if err != nil {
			return badArgNum, i + 1 + 1 // skip paren. + (1)
		}

		// fix num (and automagically convert "[0]" into badArgNum)
		return int(v64) - 1, i + 1 + 1 // DOC: arg numbers are one-indexed and skip paren. + (1)
	}

	return badArgNum, 1
}

// try to parse int value sequence in string `s` from `start` pos up to last seq dec char
// SEE src/fmt/print.go::parsenum()
//go:nosplit
func pickNumValue(s string, start uint) (num int, j uint) {

	end := uint(len(s))

	if start >= end {
		return absentValue, end
	}

	// here `num` is default int zero value 0

	for j = start; (j < end) && ('0' <= s[j]) && (s[j] <= '9'); j++ {

		num = num*10 + int(s[j]-'0')

		// check overflow
		if num > maxNum {
			return absentValue, end
		}
	}

	if j == start { // digits is not found
		num = absentValue
	}

	return num, j
}

//

// may return `badVerb` for `verb` that is missing in `verbMap`
//go:nosplit
func char2verb(c rune) (verb verb, isUpper bool) {

	// test for upper
	if isUpper = unicode.IsUpper(c); isUpper {
		c = unicode.ToLower(c)
	}

	verb = verbMap[c]

	// if wrong verb (not found in verbMap), mark it
	if verb == verbNone {
		verb = badVerb
	}

	return verb, isUpper
}

//

// SEE src/fmt/print.go::(*pp).badArgNum()
// inlined
//go:nosplit
func badArgNumValue(verb rune) string {
	// compiler optimization CL 3163 https://go-review.googlesource.com/c/go/+/3163
	return percentBangString + string(verb) + badIndexString // e2h
}
