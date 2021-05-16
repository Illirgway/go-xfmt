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
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
	"unsafe"
)

//TODO

// buffer

// go test -count=1 -v -run "^TestBuffer"

// go test -count=1 -v -run "^TestBufferSizeClassAutoAdjustment$"
// NOTE `set GOARCH=386` to test on x32
func TestBufferSizeClassAutoAdjustment(t *testing.T) {

	var b buffer

	if got := unsafe.Sizeof(b); got != bufferSizeClass {
		t.Fatalf("buffer size mismatch with its size-class: want %d, got %d", bufferSizeClass, got)
	}

	if off := unsafe.Offsetof(b.inpbuf); off != bufferHeaderSize {
		t.Fatalf("buffer header size mismatch: want %d, got %d", bufferHeaderSize, off)
	}

	t.Log(unsafe.Alignof(b))
}

// go test -count=1 -v -run "^TestBufferInitRelease$"
func TestBufferInitRelease(t *testing.T) {

	var (
		b buffer
		p bufferpool
	)

	// init

	b.init(nil)

	if b.buf == nil {
		t.Fatal("buffer backary is nil")
	}

	if l := b.Len(); l != 0 {
		t.Fatalf("buffer buf size is not zero: %d", l)
	}

	if c := b.Cap(); c != minBufDefaultSize {
		t.Fatalf("buffer buf wrong initial capacity: got %d, want %d", c, minBufDefaultSize)
	}

	oldBuf := b.buf

	if b.p != nil {
		t.Fatal("buffer p ptr must be nil but isn't")
	}

	b.init(&p)

	if ptr := &p; b.p != ptr {
		t.Fatalf("buf init error: p must be %#v, got %#v", ptr, b.p)
	}

	if !IsIdenticalByteSlice(oldBuf, b.buf) {
		t.Fatalf("b.buf bakary mismatch: want %#v, got %#v", oldBuf, b.buf)
	}

	// release

	var data = []byte("some data")

	b.buf = append(b.buf, data...)

	oldCap := cap(b.buf)

	b.reset()

	if len(b.buf) != 0 {
		t.Fatalf("buf length after reset != 0: %d (%#v)", len(b.buf), b.buf)
	}

	if c := b.Cap(); oldCap != c {
		t.Fatalf("buf cap mismatch (lost) after Free: want %d, got %d", oldCap, c)
	}

	if !IsEqualByteSlicesBakAry(oldBuf, b.buf) {
		t.Fatalf("wrong buf bakary after free: want %#v, got %#v", oldBuf, b.buf)
	}

	if b.p != nil {
		t.Fatalf("b pool has not detached: %#v", b.p)
	}

	b.buf = make([]byte, 0, maxAllowedBufSize+10)

	b.reset()

	if b.buf != nil {
		t.Fatalf("buffer extra large bakary preserved after reset: %#v", b.buf)
	}
}

// go test -count=1 -v -run "^TestBufferStdFlow$"
func TestBufferStdFlow(t *testing.T) {

	const (
		carstr  = "prefix\n"
		midstr  = "string with spaces"
		cdrstr  = "\tsuffix"
		midchar = '\b'

		resultstr = carstr + midstr + string(midchar) + cdrstr

		rawbytestring = "\n\nbyteslice\n"
		finalresult   = resultstr + rawbytestring
	)

	var (
		rawbytes = []byte(rawbytestring)
	)

	var b buffer

	// init has benn already tested above
	b.init(nil)

	sz := len(resultstr)

	b.Grow(sz)

	if cap(b.buf) != minBufDefaultSize {
		t.Fatalf("buffer.Grow fails, mismatch cap: want %d, got %d", sz, cap(b.buf))
	}

	if c := b.Cap(); c != cap(b.buf) {
		t.Fatalf("buffer.Cap mismatch result: want %d, got %d", c, cap(b.buf))
	}

	b.WriteString(carstr)
	b.WriteString(midstr)
	b.WriteByte(midchar)
	b.WriteString(cdrstr)

	if s := b.String(); s != resultstr {
		t.Fatalf("buffer write fns fail, mismatch result: want %s, got %s", resultstr, s)
	}

	b.Write(rawbytes)

	if s := b.String(); s != finalresult {
		t.Fatalf("buffer.Write fails with mismatch result: want %s, got %s", finalresult, s)
	}

	if l := b.Len(); l != len(b.buf) {
		t.Fatalf("buffer.Len unexpected mismatch result: want %d, got %d", len(b.buf), l)
	}

	if c, w := b.Cap(), minBufDefaultSize; c != w {
		t.Fatalf("unecessary buffer cap after Write: want %d, got %d", w, c)
	}

	const (
		advance      = "thisisadbvancestring"
		finadvresult = finalresult + advance
	)

	oldCap := b.Cap()

	tail := b.Advance(len(advance))

	if len(tail) != len(advance) {
		t.Fatalf("buffer.Advance error - mismatch len: want %d, got %d", len(advance), len(tail))
	}

	if c, w := b.Cap(), oldCap+len(finalresult)+len(advance); c != w {
		t.Fatalf("unecessary buffer cap after Advance: want %d, got %d", w, c)
	}

	for i := 0; i < len(advance); i++ {
		tail[i] = advance[i]
	}

	if s := b.String(); s != finadvresult {
		t.Fatalf("buffer.Advance error: write to advance tail has mismatch result: want %s, got %s", finadvresult, s)
	}

	if bb := b.Bytes(); !IsIdenticalByteSlice(bb, b.buf) {
		t.Fatalf("buffer.Bytes returns mismatch bakary: want %#v, got %#v", b.buf, bb)
	}

	oldCap = b.Cap()

	// tempbuf
	tb := b.TempBuf()

	if len(tb) != 0 {
		t.Fatalf("wrong temp buf len: %d (want 0)", len(tb))
	}

	if c := cap(tb); c != inplaceBufSize {
		t.Fatalf("wrong temp buf cap: want %d got %d", inplaceBufSize, c)
	}

	if !IsIdenticalByteSlice(tb, b.inpbuf[:0]) {
		t.Fatalf("unexpected temp buf slice bakary ptr: want %p, got %p", &b.inpbuf, &tb[:1][0])
	}

	b.reset()

	if s := b.String(); s != "" {
		t.Fatalf("buffer.reset fails: buf is not empty: %#v", b.buf)
	}

	if c := b.Cap(); oldCap != c {
		t.Fatalf("buffer.reset unexpectedly change buf cap: want %d, got %d", oldCap, c)
	}

	b.reset()
}

// pool

func testPoolBufRW(p *bufferpool, mark string) error {

	const n = 20

	for i := 0; i < n; i++ {

		b := p.Get()

		want := fmt.Sprintf("value: %s %-+ d", mark, i)

		b.WriteString("value: ")
		b.WriteString(mark)
		b.WriteByte(' ')
		b.WriteString(fmt.Sprintf("%-+ d", i))

		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

		if s := b.String(); s != want {
			return fmt.Errorf("seq r/w unexpected result: want <%s>, got <%s>", want, s)
		}

		b.Free()
	}

	return nil
}

// go test -count=1 -v -run "^TestPoolSerialFlow$"
func TestPoolSerialFlow(t *testing.T) {

	var p bufferpool

	defer func() {
		if err := recover(); err != nil {
			t.Error(err)
		}
	}()

	// test nil checking
	p.Put(nil)

	//

	b := p.Get()

	if b.p != &p {
		t.Fatalf("buffer p wrong value after getting from pool: watn %p, got %p", &p, b.p)
	}

	oldB := b

	const (
		sv1 = "some string\n"
		sv2 = "\tanother line of text\n"

		resultsv = sv1 + sv2
	)

	b.WriteString(sv1)
	b.WriteString(sv2)

	if s := b.String(); s != resultsv {
		t.Fatalf("buffer.WriteString has been failed, mismatch result: want %s, got %s", resultsv, s)
	}

	oldBuf := b.Bytes()
	oldCap := b.Cap()

	b.Free()

	// b postusage ONLY for tests purposes

	if b.p != nil {
		t.Fatalf("wrong buffer p ptr after Free: got %p, must be nil ", b.p)
	}

	if c := b.Cap(); c != oldCap {
		t.Fatalf("wrong buffer cap value after Free: want %d, got %d", oldCap, c)
	}

	if bb := b.Bytes(); !IsEqualByteSlicesBakAry(oldBuf, bb) {
		t.Fatalf("mismatch buffer bakary after free: want %p, got %p", oldBuf, bb)
	}

	b = p.Get()

	if b != oldB {
		t.Fatalf("second pool buffer mismatch with first: want %#v, got %#v", oldB, b)
	}

	curContent := string(b.Bytes()[:len(resultsv)])

	if curContent != resultsv {
		t.Fatalf("wrong bakary content value in same buffer after pooling: want %s, got %s", resultsv, curContent)
	}

	// now test sequential r/ws

	if err := testPoolBufRW(&p, "seq"); err != nil {
		t.Fatal(err)
	}
}

// go test -count=1 -v -run "^TestPoolParallelFlow$"
// go test -count=1 -o fmt.exe -gcflags "-m -m -d=ssa/check_bce/debug=1" -v -run "^TestPoolParallelFlow$" 2> fmt.log
// go tool objdump -S -s "go-xfmt" fmt.exe > fmt.disasm
func TestPoolParallelFlow(t *testing.T) {

	const (
		concurrent = 30
		timeout    = (concurrent / 2) * time.Second
	)

	var (
		wg sync.WaitGroup
		p  bufferpool
	)

	// immediately set the required number of concurrent threads
	wg.Add(concurrent)

	concurrentRoutineClosure := func(p *bufferpool, idx int) {

		mark := fmt.Sprintf("parallel:%d", idx)

		if err := testPoolBufRW(p, mark); err != nil {
			t.Fatalf("%s: %v", mark, err)
		}

		wg.Done()
	}

	for i := 0; i < concurrent; i++ {
		go concurrentRoutineClosure(&p, i)
	}

	stopCh := make(chan struct{})

	go func() {
		wg.Wait()
		close(stopCh)
	}()

	select {
	case <-stopCh:
	case <-time.After(timeout):
		t.Fatalf("timeout after %d", timeout/time.Second)
	}
}

// issues

// go test -count=1 -v -run "^TestUnusedBufferBakAryLeakage"

func TestUnusedBufferBakAryLeakage(t *testing.T) {

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

func TestUnusedBufferBakAryLeakagePool(t *testing.T) {

	var p bufferpool

	mallocs := testing.AllocsPerRun(100, func() {
		b := p.Get()
		b.Free()
	})

	if mallocs > 0 {
		t.Fatalf("unused pool buffer bak ary leakage has been detected: %v", mallocs)
	}
}
