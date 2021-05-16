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
	"bytes"
	"fmt"
	"runtime"
	"testing"
)

func hideFromVet(s string) string {
	return s
}

// go test -count=1 -v -run "^TestPrinterSprintf$"
// go test -count=1 -o fmt.exe -gcflags "-m -m -d=ssa/check_bce/debug=1" -v -run "^TestPrinterSprintf$" 2> fmt.log
// go tool objdump -S -s "go-xfmt" fmt.exe > fmt.disasm
func TestPrinterSprintf(t *testing.T) {

	// TODO append additional cases to format string

	const (
		format = "head string :: %s mid %% => joined %13.7[6]s%%%.13q%[2]x %#q %#x // %# X\n %#0 -+[2]s ignores flags %#23.15q <==> cdr"
		arg1   = "1HnBCQV$&Y25@WpMNMz*LI42GT6hwNI"
		arg2   = "t@ALoGT+8_-K@ss%RrF!fT5KdG5?O-e"
		arg3   = "Cwf3_tfEMHtK^@_bT1d2#zO-&Buk22@"
		arg4   = "HR6L!DY|cE_Q|mfQ^D62YNTWkA@WCtx"
		arg5   = "53aXrA@ZuX+i7V+8rtaneNkdNw^fFk+"
		arg6   = "j%bQiwg*PbTpQqj&R1E!nUwMpFj31WC"
		arg7   = "N$i0xM$4|94+9ZonGqweYLdnF1Vl=!!"
	)

	// without caching
	r := esprintf(format, arg1, arg2, arg3, arg4, arg5, arg6, arg7)

	//t.Log(r)

	if w := fmt.Sprintf(hideFromVet(format), arg1, arg2, arg3, arg4, arg5, arg6, arg7); r != w {
		t.Fatalf("unexpected result: want <%s> got <%s>", w, r)
	}
}

// benchmarks
// go test -bench "^BenchmarkLinearXfmtOnly$" -run "^$" -benchmem
// go test -bench "^BenchmarkLinearXfmtOnly$" -run "^$" -benchmem -cpuprofile cpu.pprof -memprofile mem.pprof

const (
	linearBenchTypeSprintf = iota
	linearBenchTypeFprintf
)

type linearBenchListItem struct {
	name   string
	typ    uint
	format string
	args   []string
}

func (e *linearBenchListItem) run(b *testing.B) {

	prefix := "Sprintf"

	if e.typ == linearBenchTypeFprintf {
		prefix = "Fprintf"
	}

	name := prefix + e.name

	runtime.GC()

	b.Run(name, func(bb *testing.B) {

		if e.typ == linearBenchTypeSprintf {

			bb.ResetTimer()

			for i := 0; i < bb.N; i++ {
				_ = Sprintf(e.format, e.args...)
			}

			return
		}

		// linearBenchTypeFprintf

		var buf bytes.Buffer

		bb.ResetTimer()

		for i := 0; i < bb.N; i++ {
			_, _ = Fprintf(&buf, e.format, e.args...)
		}

	})
}

var linearBenchList = [...]linearBenchListItem{
	{
		"Padding",
		linearBenchTypeSprintf,
		"%16s",
		[]string{tftc1arg5utf8MS_HnS},
	},
	{
		"String",
		linearBenchTypeSprintf,
		"%s",
		[]string{tftc1arg13str},
	},
	{
		"TruncateString",
		linearBenchTypeSprintf,
		"%.4s",
		[]string{tftc1arg14utf8cjklong},
	},
	{
		"QuoteString",
		linearBenchTypeSprintf,
		"%q",
		[]string{tftc1arg14utf8cjklong},
	},
	{
		"PrefixedString",
		linearBenchTypeSprintf,
		"This is some meaningless prefix text that needs to be scanned %s",
		[]string{tftc1arg5utf8MS_HnS},
	},
	{
		"HexString",
		linearBenchTypeSprintf,
		"% #x",
		[]string{tftc1arg15str},
	},
	{
		"ManyArgs",
		linearBenchTypeSprintf,
		"%2s/%2s/%2s %s:%s:%s %s %s\n",
		[]string{"3", "4", "5", "11", "12", "13", "hello", "world"},
	},
	{
		"ManyArgs",
		linearBenchTypeFprintf,
		"%2s/%2s/%2s %s:%s:%s %s %s\n",
		[]string{"3", "4", "5", "11", "12", "13", "hello", "world"},
	},
}

func BenchmarkLinearXfmtOnly(b *testing.B) {
	for i := 0; i < len(linearBenchList); i++ {
		linearBenchList[i].run(b)
	}
}
