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
	"sync/atomic"
)

// repetitions counting cache
// PPSL: may change to or add delta time cache (between two usages of the same format value)

const (
	CacheAlways      uint = 0
	CacheRepetitions      = 1        // cache only if format occurred more than once
	CacheDisabled         = ^uint(0) // max uint value
)

// use uintptr for atomic store/load
var cacheThreshold uintptr = uintptr(CacheAlways)

// thread-safe
// inlined
//go:nosplit
func SetCacheThreshold(threshold uint) {
	atomic.StoreUintptr(&cacheThreshold, uintptr(threshold))
}

// thread-safe
// inlined
//go:nosplit
func CacheThreshold() uint {
	return uint(atomic.LoadUintptr(&cacheThreshold))
}

//

// NOTE due to the fact that the overwhelming majority of `format` values are string constants (i.e., the backed bytearray
//      of their data is in the non-heap constant section (bss)), the string keys of the following hashes will not
//      actually use heap-allocated strings as values, so these hashes are cheap in terms of memory usage

type thresholdCounters struct {
	lock     sync.Mutex
	counters map[string]uint // lazy
}

// thread-safe
//go:nosplit
func (tc *thresholdCounters) Count(key string) (count uint) {

	tc.lock.Lock()
	// hate defer, but we should unlock in case of any write (== memalloc) error
	defer tc.lock.Unlock()

	if tc.counters != nil {
		count = tc.counters[key] // got 0 if not found
	} else {
		// lazy map init, so no key in counters ==> count = 0
		tc.counters = make(map[string]uint, 1)
	}

	// should inc usage counter ...
	count++

	// ... and store back
	tc.counters[key] = count

	return count
}

var countersCache thresholdCounters

//

type formatCache struct {
	lock  sync.RWMutex
	cache map[string]xfmt // NOTE xfmt by value
}

// inlined
// thread-safe
//go:nosplit
func (c *formatCache) Get(format string) (fmt xfmt, has bool) {

	c.lock.RLock()

	fmt, has = c.cache[format]

	c.lock.RUnlock()

	return fmt, has
}

// thread-safe
//go:nosplit
func (c *formatCache) Set(format string, fmt xfmt) {

	c.lock.Lock()
	// hate defer, but we should unlock in case of any write (== memalloc) error
	defer c.lock.Unlock()

	// lazy map init
	if c.cache == nil {
		c.cache = make(map[string]xfmt, 1)
	}

	c.cache[format] = fmt
}

var xfmtCache formatCache
