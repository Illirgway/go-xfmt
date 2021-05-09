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
	"unicode/utf8"
)

// inlined
//go:nosplit
func tailnoop(w io.Writer) (n int, err error) {
	return 0, nil
}

// inlined
//go:nosplit
func tailln(w io.Writer) (n int, err error) {
	return w.Write(rawLF[:])
}

// SEE https://blog.golang.org/strings#TOC_7.
// SEE src/fmt/format.go::*fmt::truncateString()
// NOTE n is amount of utf8 chars (codepoints / runes, NOT bytes !!!) that keep at the beginning
// `go1.13: cannot inline truncateTail: unhandled op FOR`
//go:nosplit
func truncateTail(s string, n int) string {

	i := uint(0) // uint automagically helps BCE optimization without additional conds

	for ; i < uint(len(s)) && (n > 0); n-- {
		_, w := utf8.DecodeRuneInString(s[i:])
		i += uint(w)
	}

	return s[:i]
}
