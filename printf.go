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
	"io"
	"os"
)

// TODO replace `by value` xfmt by clever `by ptr` *xfmt (should not to do unnecessary heapalloc for xfmt in case of CacheDisabled)
//      hint: use implicit `pxfmt := new(xfmt)` in the right place (and next `*pxfmt = xfmt`) to avoid explicit unwanted one
//      in wrong place

// NOTE xfmt is for now by value
func forgeXfmt(format string) (xfmt xfmt) {

	xfmt, has := xfmtCache.Get(format)

	threshold := CacheThreshold()

	// if already in cache, should not recache (and recounting)
	if has {
		return xfmt
	}

	// not in cache, should parse and cache if needed
	xfmt = parseFormat(format)

	if threshold != CacheDisabled {

		shouldCache := threshold == CacheAlways

		// use counters cache only if need it
		if !shouldCache {
			shouldCache = countersCache.Count(format) > threshold
		}

		if shouldCache {
			// store in cache...
			xfmtCache.Set(format, xfmt)

			// ...and then remove format value from counters cache if needed to reduce counters heapsize and memallocs
			if threshold != CacheAlways {
				countersCache.Delete(format)
			}
		}
	}

	return xfmt
}

// `go1.13: cannot inline Fprintf: function too complex: cost 137 exceeds budget 80`
func Fprintf(w io.Writer, format string, args ...string) (n int, err error) {
	xfmt := forgeXfmt(format)
	return xfmt.Fprint(w, args)
}

// `go1.13: cannot inline Printf: function too complex: cost 138 exceeds budget 80`
func Printf(format string, args ...string) (n int, err error) {
	xfmt := forgeXfmt(format)
	return xfmt.Fprint(os.Stdout, args)
}

// `go1.13: cannot inline Sprintf: function too complex: cost 127 exceeds budget 80`
func Sprintf(format string, args ...string) string {
	xfmt := forgeXfmt(format)
	return xfmt.Sprint(args)
}

// `go1.13: cannot inline Errorf: function too complex: cost 135 exceeds budget 80`
func Errorf(format string, args ...string) error {
	xfmt := forgeXfmt(format)
	return errors.New(xfmt.Sprint(args))
}
