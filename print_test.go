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
	"testing"
)

var (
	testcaseSprintList   = []string{"test string 1", "gdg", "value"}
	testcasesSprintIList = make([]interface{}, len(testcaseSprintList))
)

func init() {
	for i, v := range testcaseSprintList {
		testcasesSprintIList[i] = v
	}
}

// go test -count=1 -v
// go test -count=1 -o fmt.exe -gcflags "-m -m -d=ssa/check_bce/debug=1" -v 2> fmt.log
// go tool objdump -S -s "go-xfmt" fmt.exe > fmt.disasm

// go test -count=1 -v -run "^TestSprint1$"
// go test -count=1 -o fmt.exe -gcflags "-m -m -d=ssa/check_bce/debug=1" -v -run "^TestSprintLn1$" 2> fmt.log
// go tool objdump -S -s "go-xfmt" fmt.exe > fmt.disasm
func TestSprint1(t *testing.T) {

	s := Sprint(testcaseSprintList...)

	t.Log(s)

	rs := fmt.Sprint(testcasesSprintIList...)

	if s != rs {
		t.Fatalf("result string mismatch with ref: want %q, got %q", rs, s)
	}
}

// go test -count=1 -v -run "^TestSprintLn1$"
// go test -count=1 -o fmt.exe -gcflags "-m -m -d=ssa/check_bce/debug=1" -v -run "^TestSprintLn1$" 2> fmt.log
// go tool objdump -S -s "go-xfmt" fmt.exe > fmt.disasm
func TestSprintLn1(t *testing.T) {

	s := Sprintln(testcaseSprintList...)

	t.Log(s)

	rs := fmt.Sprintln(testcasesSprintIList...)

	if s != rs {
		t.Fatalf("result string mismatch with ref: want %q, got %q", rs, s)
	}
}
