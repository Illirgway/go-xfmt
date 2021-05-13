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
	"sync"
)

// INFO: fast pool buffer wrapper around bytebufferpool to avoid insufficiently optimized bytebufferpool.Pool.Put()
// for small buffers

const (
	// INPlace BUFfer
	is64bit                 = (^uint(0)) >> 63                   // 0 for x32, 1 for x64
	cpucacheline            = 1 << 6                             // 2^6 == 64 bytes
	bufferSizeClassId       = 9                                  // SEE src/runtime/sizeclasses.go
	bufferSizeClass         = 128                                // tail waste = 0 and is 2 * cpucacheline
	bufferHeaderSize        = (3 + 1) * (4 << is64bit)           // => {x32: 16, x64: 32}
	inplaceBufSize          = bufferSizeClass - bufferHeaderSize // => {x32: 128 - 16 = 112, x64: 128 - 32 = 96}
	minBufDefaultSize       = cpucacheline
	minBufMaxSize           = minBufDefaultSize << 1
	usesCounterThreshold    = 5000
	usesCounterInitialValue = 40
)

const (
	// NOTE "Go does not manage the large allocations with a local cache. Those allocations, greater than 32kb,
	//       are rounded up to the page size and the pages are allocated directly to the heap."
	// SEE  https://medium.com/a-journey-with-go/go-memory-management-and-allocation-a7396d430f44
	maxAllowedBufSize = 16 << 10 // 16k; 4 times less than the similar value in the `fmt` package
)

type bufferpool struct {
	pool sync.Pool // 32 or 64
	/*
		defaultSize uintptr
		// always aligned on 64 boundary
		uses    uint64  // count of pool uses (=== count of issued buffers)
		sizeSum uint64  // ∑(X; n)
		maxSize uintptr // 32 or 64
	*/
}

/*
// inlined
//go:nosplit
func newbufferpool() bufferpool {
	return bufferpool{
		//defaultSize: minBufDefaultSize,
		//maxSize: minBufMaxSize << 1,
	}
}
*/

func (p *bufferpool) Get() (b *buffer) {

	if bb := p.pool.Get(); bb != nil {
		b, _ = bb.(*buffer)
	} else {
		b = new(buffer)
	}

	//sz := uint(atomic.LoadUintptr(&p.defaultSize))

	b.init(p /*, sz*/)

	return b
}

// MATH: default size is M(X) of all prev default buf sizes except defsizes >= maxAllowedBufSize
//       M(X; n) = ∑xi / n
//       defSize(n) = M(X; n) = (x0 + x1 + x2 + ... + xn) / i = ∑xi / n
//       sumSize(n) = x0 + x1 + x2 + ... + xn = ∑xi = defSize(n) * n
//       defSize(n+1) = M(X; n+1) = (x0 + x1 + x2 + ... + xn + x(n+1)) / (n + 1) = (∑xi + x(n+1)) / (n + 1) =
//            = (sumSize(n) + x(i+1)) / (n + 1) = (defSize(n) * n + newSize) / (n + 1)  ===>
//       defSize(n) = (defSize(n - 1) * (n - 1) + curSize) / n = (defSize(n - 1) * n - defSize(n - 1) + curSize) / n =
//           = defSize(n - 1) + (curSize - defSize(n - 1)) / n
// NOTE lim{n -> +∞} (n / (n + 1)) = 1; lim{n -> +∞}( (curSize - defSize(n - 1)) / n ) = 0 (because curSize < M ==> defSize < N(M))
func (p *bufferpool) Put(b *buffer) {

	/*if b == nil {
		return
	}*/

	/*
		curSize, maxSize := uint(b.Len()), uint(maxAllowedBufSize)

		// don't take into account extra large buffers
		if curSize <= maxSize {

			// guarded from too low size values
			if curSize < minBufDefaultSize {
				curSize = minBufDefaultSize
			}

			// here `minBufDefaultSize <= curSize <= maxAllowedBufSize`

			n := atomic.AddUint64(&p.uses, 1)

			var defSize uintptr

			if n > usesCounterThreshold {
				defSize = atomic.LoadUintptr(&p.defaultSize)

				n -= usesCounterThreshold - usesCounterInitialValue	// n - usesCounterThreshold + usesCounterInitialValue

				atomic.StoreUint64(&p.uses, n)
				atomic.StoreUint64(&p.sizeSum, uint64(defSize) * n)
			}

			sumSizes := atomic.AddUint64(&p.sizeSum, uint64(curSize))

			defSize = uintptr(sumSizes / n)

			atomic.StoreUintptr(&p.defaultSize, defSize)

			maxSize = uint(atomic.LoadUintptr(&p.maxSize))

			// different coef. to avoid jitter on bound
			// 2 * defSize > maxSize
			if v := uint(defSize << 1); v > maxSize && maxSize < maxAllowedBufSize {
				maxSize = maxSize << 1
				atomic.StoreUintptr(&p.maxSize, uintptr(maxSize))	// maxSize = maxSize * 2
			} else if v := uint((defSize - defSize >> 2) >> 1 /* defSize / 2.67 ~= defSize/2 - defSize/8 * /); v < maxSize && maxSize > minBufMaxSize {
				maxSize = maxSize >> 1
				atomic.StoreUintptr(&p.maxSize, uintptr(maxSize))	// maxSize = maxSize / 2
			}
		}
	*/

	b.reset( /*maxSize*/ )

	p.pool.Put(b)
}

// buffer
// TODO adjust size of struct to suitable malloc size-class (using the size of the internal buffer)
type buffer struct {
	p   *bufferpool // 4 or 8 bytes
	buf []byte      // 3 * (4 or 8) bytes
	// NOTE total size of above fields = sizeof(uintptr) * (1 /* Ptr */ + 3 /* SliceHeader */) =
	//      = 4 * sizeof(uintptr) ==> {x32: 4 * 4 = 16, x64: 4 * 8 = 32}; delta = 32 - 16 = 16

	//inpbuf *fastrawbuf
	inpbuf [inplaceBufSize]byte

	// Total size === bufferSizeClass
}

// `Hacker's Delight`, $3.1
// inlined
//go:nosplit
func roundUpToPowOf2(x uint, p uint) uint {
	return (x + (p - 1)) & ((^p) + 1) // (x + (p - 1)) & (-p)
}

// inlined
//go:nosplit
func (b *buffer) init(p *bufferpool /*, sz uint*/) {

	if b.buf == nil {
		/*// round up to `cpucacheline`
		sz = roundUpToPowOf2(sz, cpucacheline)
		b.buf = make([]byte, 0, sz)*/
		b.buf = make([]byte, 0, minBufDefaultSize)
	}

	b.p = p
}

// inlined
//go:nosplit
func (b *buffer) reset( /*maxSz uint*/ ) {

	//if c := uint(cap(b.buf)); c > maxSz || c > maxAllowedBufSize {
	if uint(cap(b.buf)) > maxAllowedBufSize {
		b.buf = nil
	} else {
		b.buf = b.buf[:0]
	}
}

// inlined
//go:nosplit
func (b *buffer) Free() {
	b.p.Put(b)
}

// inlined
//go:nosplit
func (b *buffer) Len() int {
	return len(b.buf)
}

// inlined
//go:nosplit
func (b *buffer) Cap() int {
	return cap(b.buf)
}

// NOTE `uint` is used instead of `int` below to help BCE

// Grow grows internal buffer `buf` without reslicing (just increases the capacity)
// inlined
//go:nosplit
func (b *buffer) Grow(n int) {
	if sz := n + len(b.buf); sz > cap(b.buf) {
		b.grow(sz)
	}
}

// NOTE Grow + Expand, eq. to src/bytes.(*buffer).tryGrowByReslice() + .grow() if not ok
//go:nosplit
func (b *buffer) advance(n int) (m int) {

	m = len(b.buf)

	if n == 0 {
		return m
	}

	sz := n + m /* len(b.buf) */

	if sz > cap(b.buf) {
		b.grow(sz)
	}

	b.buf = b.buf[:sz]

	return m
}

// inlined
//go:nosplit
func (b *buffer) Advance(n int) []byte {
	m := b.advance(n)
	return b.buf[m:]
}

//go:nosplit
func (b *buffer) grow(sz int) {

	// worst case, should grow underlying storage
	// 2 * cap(b) + n, but
	//  / sz = len(b) + n ==> n = sz - len
	//  | n + len(b) > cap(b) ==> n > cap - len >= 0
	//  \ len(b) <= cap(b) ==> cap - len >= 0
	//  => n = sz - len > cap - len >= 0 ==> sz > cap >= 0 ==> sz > 0
	// cap + sz = cap + len + n <= 2 * cap + n
	// len + n > cap ==> cap + cap < cap + (len + n) <= cap + cap + n
	// ==> 2 * cap < cap + sz <= 2 * cap + n
	buf := make([]byte, len(b.buf), cap(b.buf)+sz)

	copy(buf, b.buf)

	b.buf = buf
}

// inlined
//go:nosplit
func (b *buffer) Write(p []byte) {
	// NOTE after advance() buf always has enough cap
	if m := b.advance(len(p)); m < len(b.buf) /* BCE hint, always is `true` except for `len(p) == 0` */ {
		copy(b.buf[m:], p)
	}
}

// inlined
//go:nosplit
func (b *buffer) WriteByte(c byte) {
	// NOTE after advance() buf always has enough cap
	if m := b.advance(1); m < len(b.buf) /* BCE hint, always is `true` */ {
		b.buf[m] = c
	}
}

// inlined
//go:nosplit
func (b *buffer) WriteString(s string) {
	// NOTE after advance() buf always has enough cap
	if m := b.advance(len(s)); m < len(b.buf) /* BCE hint, always is `true` except for `s == ""` */ {
		copy(b.buf[m:], s)
	}
}

// NOTE implicit copy
// inlined
//go:nosplit
func (b *buffer) String() string {
	return string(b.buf)
}

// inlined
//go:nosplit
func (b *buffer) Bytes() []byte {
	return b.buf
}

// inlined
//go:nosplit
func (b *buffer) TempBuf() []byte {
	return b.inpbuf[:0]
}
