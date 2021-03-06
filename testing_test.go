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
	"reflect"
	"unsafe"
)

// utils for testing

func IsIdenticalByteSlice(a, b []byte) bool {
	return (len(a) == len(b)) && (cap(a) == cap(b)) && IsEqualByteSlicesBakAry(a, b)
}

// SEE https://stackoverflow.com/a/53010178
func IsEqualByteSlicesBakAry(a, b []byte) bool {
	return (*reflect.SliceHeader)(unsafe.Pointer(&a)).Data == (*reflect.SliceHeader)(unsafe.Pointer(&b)).Data
}
