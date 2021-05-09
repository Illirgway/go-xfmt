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

// Strings for use with buffer.WriteString.
// This is less overhead than using buffer.Write with byte arrays.
const (
	percentString     = "%"
	commaSpaceString  = ", "
	reflectStringType = "string" // avoid import reflect package
	percentBangString = "%!"
	missingString     = "(MISSING)"
	badIndexString    = "(BADINDEX)"

	// format errors
	extraString    = "%!(EXTRA "
	badWidthString = "%!(BADWIDTH)"
	badPrecString  = "%!(BADPREC)"
	noVerbString   = "%!(NOVERB)"
	nilToken       = "%!(NILTOKEN)"
)

// chars aliases
const (
	charEquals      = '='
	charPercent     = '%'
	charOpenArgNum  = '['
	charCloseArgNum = ']'
	charAsterisk    = '*'
	charDot         = '.'
	charLeftParens  = '(' // parenthesis
	charRightParens = ')' // parenthesis
	charSpace       = ' '
	charZero        = '0'
	charBackquote   = '`'
)

// const strings
const (
	percentStr       = string(charPercent)
	doublePercentStr = percentStr + percentStr
	backquoteStr     = string(charBackquote)

	emptyString = ""
)
