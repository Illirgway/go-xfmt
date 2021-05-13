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
)

type xfmt struct {
	tokens  []token
	args    uint // needed args count, may differ from len(tokens) due to [n] notation
	minSize int  // minimal size of result string if all of fmt args are empty strings (== "")
	// NOTE small size struct, may be passed by value
}

var fmtprintbufpool bufferpool

func (fmt *xfmt) bprint(args []string) (buf *buffer) {

	// fast-paths
	// - format is empty string and no args
	if (len(fmt.tokens) == 0) && (len(args) == 0) {
		return nil
	}

	buf = fmtprintbufpool.Get()

	// try to minimize memallocs
	buf.Grow(fmt.minSize)

	// - format is a single raw const string value without any verb
	if (len(fmt.tokens) == 1) && (fmt.args == 0) && (len(args) == 0 /* implies `args == nil` */) &&
		(fmt.tokens[0].verb == verbNone) {
		buf.WriteString(fmt.tokens[0].value)
		//buf.SetString(fmt.tokens[0].value)
		return
	}

	for i := 0; i < len(fmt.tokens); i++ {
		fmt.tokens[i].format(buf, args)
	}

	// uint automagically helps BCE optimization without additional conds
	if nArgs := uint(len(args)); nArgs > fmt.args {

		buf.WriteString(extraString)

		for i := fmt.args; i < nArgs; i++ {

			if i > fmt.args {
				buf.WriteString(commaSpaceString)
			}

			buf.WriteString(reflectStringType)
			buf.WriteByte(charEquals)
			buf.WriteString(args[i])
		}

		buf.WriteByte(charRightParens)
	}

	return buf
}

func (fmt *xfmt) Fprint(w io.Writer, args []string) (n int, err error) {

	// fast-paths
	// - format is empty string and no args
	if (len(fmt.tokens) == 0) && (len(args) == 0) {
		return 0, nil
	}

	// common case, buf needed

	// wrap around bprint

	b := fmt.bprint(args)

	// impossible situation
	if b == nil || b.Len() == 0 {
		return 0, nil
	}

	// WARN write buf BEFORE return it to pool
	n, err = w.Write(b.Bytes())

	b.Free()

	return n, err
}

func (fmt *xfmt) Sprint(args []string) (s string) {

	// fast-paths
	// - format is empty string and no args
	if (len(fmt.tokens) == 0) && (len(args) == 0) {
		return ""
	}

	// - format is a single raw const string value without any verb
	if (len(fmt.tokens) == 1) && (fmt.args == 0) && (len(args) == 0 /* implies `args == nil` */) &&
		(fmt.tokens[0].verb == verbNone) {
		return fmt.tokens[0].value
	}

	// common case, buf needed

	// wrap around bprint

	b := fmt.bprint(args)

	// impossible situation
	if b == nil || b.Len() == 0 {
		return ""
	}

	// implicitly copy buf as string to free buf and ret it back to pool
	// WARN make string from buf BEFORE return buf to pool
	s = b.String()

	b.Free()

	return s
}
