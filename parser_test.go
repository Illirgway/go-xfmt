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

type pickNumValueTestCase struct {
	format string
	start  uint
	num    int
	j      uint
}

func (tc *pickNumValueTestCase) run() error {

	num, j := pickNumValue(tc.format, tc.start)

	switch {
	case (num != tc.num) && (j != tc.j):
		return fmt.Errorf("mismatch `num` and `j` values: num want %d got %d, j want %d got %d", tc.num, num, tc.j, j)

	case num != tc.num:
		return fmt.Errorf("mismatch `num`: want %d got %d", tc.num, num)

	case j != tc.j:
		return fmt.Errorf("mismatch `j`: want %d got %d", tc.j, j)
	}

	return nil
}

var pickNumValueTestCases1 = [...]pickNumValueTestCase{
	{"1234", 4, absentValue, 4}, // start == end
	{"1234", 7, absentValue, 4}, // start > end
	{"a123", 0, absentValue, 0}, // starts from char not num - digits is not found

	{"123a", 0, 123, 3},
	{"12a3", 0, 12, 2},
	{"1234", 0, 1234, 4},

	{"1a234", 1, absentValue, 1},
	{"9999999x", 0, absentValue, 8}, // check overflow
}

// go test -count=1 -v -run "^TestPickNumValue1$"
func TestPickNumValue1(t *testing.T) {

	for i := 0; i < len(pickNumValueTestCases1); i++ {

		tcase := &pickNumValueTestCases1[i]

		if err := tcase.run(); err != nil {
			t.Fatalf("%d (%#v => %v): %v", i, tcase.format, tcase.start, err)
		}
	}

}
