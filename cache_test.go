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
	"sync"
	"testing"
)

// thread-unsafe
func purgeCaches() {
	countersCache.counters = nil
	xfmtCache.cache = nil
}

//TODO

// check concurrent stability
// go test -count=1 -v -run "^TestCacheConcurrentStability$"
func TestCacheConcurrentStability(t *testing.T) {

	const c = 10

	sources := make([]string, c)

	for i := 0; i < c; i++ {
		sources[i] = fmt.Sprintf("source value: #%d", i)
	}

	startCh := make(chan struct{})

	var (
		cache formatCache
		g     sync.WaitGroup
	)

	setter := func(i int) {

		source := sources[i]

		t := token{
			verb:  verbNone,
			value: source,
			arg:   uint(i),
		}

		v := xfmt{
			[]token{t},
			uint(i),
			0,
		}

		<-startCh

		cache.Set(source, v)

		g.Done()
	}

	g.Add(len(sources))

	for i := range sources {
		go setter(i)
	}

	close(startCh)

	g.Wait()

	if l := cache.Len(); l != len(sources) {
		t.Fatalf("len mismatch: want %d, got %d", len(sources), l)
	}

	for i, s := range sources {
		v, has := cache.Get(s)

		if !has {
			t.Fatalf("%d (%s): absent cache entry", i, s)
		}

		if v.args != uint(i) {
			t.Fatalf("%d (%s): mismatch args value: got %d", i, s, v.args)
		}

		if l := len(v.tokens); l != 1 {
			t.Fatalf("%d (%s): wrong tokens list len: %d", i, s, l)
		}

		tt := v.tokens[0]

		if tt.verb != verbNone {
			t.Fatalf("%d (%s): wrong token verb: %d", i, s, tt.verb)
		}

		if tt.value != s {
			t.Fatalf("%d (%s): mismatch token value: got %s", i, s, tt.value)
		}

		if tt.arg != uint(i) {
			t.Fatalf("%d (%s): mismatch token arg value: got %d", i, s, tt.arg)
		}
	}
}
