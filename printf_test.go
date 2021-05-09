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

import "testing"

// TODO

// go test -count=1 -v -run "^TestPrinterSprintf1$"
// go test -count=1 -o fmt.exe -gcflags "-m -m -d=ssa/check_bce/debug=1" -v -run "^TestPrinterSprintf1$" 2> fmt.log
// go tool objdump -S -s "go-xfmt" fmt.exe > fmt.disasm
func TestPrinterSprintf1(t *testing.T) {

	SetCacheThreshold(CacheAlways)

	t.Log(Sprintf("%% ||| %s ||| %q ||| % #x", tpfc1argpct1, tpfc1argbqqt1, tpfc1argutf8c))
	t.Logf("%#v", countersCache.counters)
	t.Logf("%#v", xfmtCache.cache)

	for k, v := range xfmtCache.cache {
		t.Logf("%q => %#v", k, v)
	}
}
