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
	"testing"
	"unsafe"
)

//TODO

// go test -count=1 -v -run "^TestBufferSizeClassAutoAdjustment1$"
// NOTE `set GOARCH=386` to test on x32
func TestBufferSizeClassAutoAdjustment1(t *testing.T) {

	var b buffer

	if got := unsafe.Sizeof(b); got != bufferSizeClass {
		t.Fatalf("buffer size mismatch with its size-class: want %d, got %d", bufferSizeClass, got)
	}

	if off := unsafe.Offsetof(b.inpbuf); off != bufferHeaderSize {
		t.Fatalf("buffer header size mismatch: want %d, got %d", bufferHeaderSize, off)
	}

	t.Log(unsafe.Alignof(b))
}

// go test -count=1 -v -run "^TestUnusedBufferBakAryLeakage1$"
func TestUnusedBufferBakAryLeakage1(t *testing.T) {

	var b buffer

	// DOC To compute the number of allocations, the function will first be run once as a warm-up
	mallocs := testing.AllocsPerRun(100, func() {
		b.init(nil)
		b.reset()
	})

	if mallocs > 0 {
		t.Fatalf("unused buffer bak ary leakage has been detected: %v", mallocs)
	}
}

// TODO add similar test for b.Free -> bufferpool.Put(b)
