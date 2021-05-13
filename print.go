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
	"io"
	"os"
)

const (
	Space = ' '
	LF    = '\n'
)

const (
	strNUL = ""
	strLF  = string(LF)
)

var (
	rawSpace = [1]byte{Space}
	rawLF    = [1]byte{LF}

	stdprintbufpool bufferpool
)

// return value may be nil === no bytes written (empty buf)
//go:nosplit
func bprint(ln bool, s ...string) (buf *buffer) {

	// fast-path
	if len(s) == 0 {
		return nil
	}

	buf = stdprintbufpool.Get()

	// TODO? once time before merge precalc size and preallocate if needed?

	for i := range s {
		if ln && (i > 0) {
			buf.WriteByte(Space)
		}

		buf.WriteString(s[i])
	}

	if ln {
		buf.WriteByte(LF)
	}

	return buf
}

// `go1.13: cannot inline sprint: function too complex: cost 148 exceeds budget 80`
//go:nosplit
func sprint(ln bool, s ...string) (r string) {

	buf := bprint(ln, s...)

	if buf == nil {
		if ln {
			return strLF
		}

		return strNUL
	}

	// implicit copy []byte->string (runtime.slicebytetostring)
	// WARN make string from buf BEFORE free buf
	r = buf.String()

	// should return buf to its pool
	buf.Free()

	return r
}

// `go1.13: cannot inline fprint: function too complex: cost 281 exceeds budget 80`
//go:nosplit
func fprint(w io.Writer, ln bool, s ...string) (n int, err error) {

	buf := bprint(ln, s...)

	if buf == nil {
		if ln {
			return w.Write(rawLF[:])
		}

		return 0, nil
	}

	// should write buf before free it
	n, err = w.Write(buf.Bytes())

	// should return buf to its pool
	buf.Free()

	return n, err
}

/*
// TODO? on-stack []byte buf with auto-flush on fill for reducing w.Write calls
//       также можно использовать github.com/valyala/bytebufferpool.Pool
func fprintOld(w io.Writer, ln bool, s ...string) (n int, err error) {

	tail := tailnoop

	if ln {
		tail = tailln
	}

	// fast-path
	if len(s) == 0 {
		return tail(w)
	}

	var nn int

	// init bytes written counter
	n = 0

	// fast-path
	if len(s) == 1 {

		// skip empty string
		if s[0] != "" {
			// avoid string->[]byte memalloc
			if n, err = w.Write(xruntime.AssignString2SliceUnsafe(&s[0])); err != nil {
				return n, err
			}
		}

		nn, err = tail(w)

		return n + nn, err
	}

	// optimization - skip empty string for both ends (trim s)
	i, j := uint(0), uint(len(s) - 1)	// uint helps BCE

	for found := true; found && (i <= j) && (j < uint(len(s)) /* helps BCE, always true * /); {

		found = false

		if s[i] == "" {
			found = true
			i++
		}

		if s[j] == "" {
			found = true
			j--
		}
	}

	// stop if only empty strings found
	if (i > j) /* omit bounds checking below * / || (j >= uint(len(s))) /** / {
		// but should print ln if any
		return tail(w)
	}

	// [a b c]  ["" a b]  ["" a ""]
	//  i   j       i j       ij

	// iterate [i; j] inclusively
	for ; (i <= j) && (err == nil); i++ {

		nn, err = w.Write(
			// avoid string->[]byte memalloc
			xruntime.AssignString2SliceUnsafe(&s[i]))

		// regardless of err count the number of bytes written
		n += nn

		// for ln version should insert space between strings, but only between and if no error
		if ln && (i < j) && (err == nil) {
			nn, err = w.Write(rawSpace[:])

			// regardless of err count the number of bytes written
			n += nn
		}
	}

	if err != nil {
		return n, err
	}

	nn, err = tail(w)

	// regardless of err count the number of bytes written
	return n + nn, err
}
*/

// inlined
//go:nosplit
func Fprint(w io.Writer, s ...string) (n int, err error) {
	return fprint(w, false, s...)
}

// `go1.13: cannot inline Print: function too complex: cost 87 exceeds budget 80`
//go:nosplit
func Print(s ...string) (n int, err error) {
	return Fprint(os.Stdout, s...)
}

//go:nosplit
func Sprint(s ...string) string {
	/*var b strings.Builder
	// avoid e2h for b
	_, _ = Fprint((*strings.Builder)(xruntime.NoEscape(unsafe.Pointer(&b))), s...)
	return b.String()*/

	return sprint(false, s...)
}

// `Spaces are always added between operands and a newline is appended.`
// inlined
//go:nosplit
func Fprintln(w io.Writer, s ...string) (n int, err error) {
	return fprint(w, true, s...)
}

// `go1.13: cannot inline Println: function too complex: cost 87 exceeds budget 80`
//go:nosplit
func Println(s ...string) (n int, err error) {
	return Fprintln(os.Stdout, s...)
}

//go:nosplit
func Sprintln(s ...string) string {
	/*var b strings.Builder
	// avoid e2h for b
	_, _ = Fprint((*strings.Builder)(xruntime.NoEscape(unsafe.Pointer(&b))), s...)
	return b.String()*/

	return sprint(true, s...)
}
