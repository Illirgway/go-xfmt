
# Go xfmt package
Almost drop-in replacement of std `fmt` package but only for `string` args

# Overview
This package is an almost drop-in replacement for std `fmt` package fns `*print`, `*printf`, and
`Errorf` with advantages over the `fmt` package, however only exclusively for args of `string` type 
(or cast to this type). Advantages:
* `reflect` package is excluded, so no heavy typecast and conversions with implicit checks
* `format` string value is parsed as tokens list' struct, suitable for caching
* configurable caching of parsed format's tokens list' structs to avoid doing the same parsing job multiple times
* reusable [byte buffer pool](https://github.com/valyala/bytebufferpool) (for simple `*print` functions 
  for results building) + special `buffer` type with all of needed fns for `*printf` processing + with inplace subbuf for   
  internal `fmtXxx` fns' - mostly used by similar to strconv.Append*(buf []byte) fns, e.g. strconv.AppendQuote)   
  drastically reduces heap reallocate-and-copy ops during write of results to buffer byteslice
* unnecessary memory allocations have been drastically reduced - there are no unnecessary allocations at all
* a lot of BCE optimizations performed (and many more needs to be done)
* as a result, it is more faster and cheaper (approximately up to 2 times for both cpu and mem)

## Status: _ALPHA_
* **CAN'T BE USED IN PRODUCTION FOR NOW**
* finish up currently unfinished tests (printf_test)
* adapted tests from the original `fmt` package are required
* more tests are required (up to full code coverage: buffer_test, cache_test)

## API
It has the same subset of functions as the original `fmt` package, but with args of `string` type instead of 
`interface{}` type

```gotemplate

// simple fns without format

func Fprint(w io.Writer, s ...string) (n int, err error)
func Print(s ...string) (n int, err error)
func Sprint(s ...string) string

// ln versions

func Fprintln(w io.Writer, s ...string) (n int, err error)
func Println(s ...string) (n int, err error)
func Sprintln(s ...string) string 

// format fns

func Fprintf(w io.Writer, format string, args ...string) (n int, err error)
func Printf(format string, args ...string) (n int, err error)
func Sprintf(format string, args ...string) string

// especial fns

func Errorf(format string, args ...string) error

// cache control

const (
	CacheAlways // always cache any fmt
	CacheRepetitions // cache only if format occurred more than once
	CacheDisabled // disables caching at all
)

// both thread-safe
func SetCacheThreshold(threshold uint)
func CacheThreshold() uint

```
 
### Caching

By default, package cache tokens' list of any parsed `format` (`CacheAlways`), but this behavior is controlled by 
`cacheThreshold` internal value. Use `SetCacheThreshold` to set how many times format should be parsed (saw)
by the parser before it will be cached. To check the threshold value, use `CacheThreshold` fn. Both fns are 
thread-safe (atomic r/w).

 
### Mismatches with std `fmt`
* simple no-ln fns `Sprint`, `Fprint`, and `Print` don't try to recognize args' initial types at all, so unlike such 
  `fmt` fns they can't determine when to `Spaces are added between operands when neither is a string` and
  don't add spaces between args et all
* some errors mark (especially for errors related to the tail of `format`) of formatting fns may differ from such  
  returned by original `fmt` format fns

## Contributing

Plz don't send pull requests for now, only write issue tickets

## LICENSE

This program is free software: you can redistribute it and/or modify it under the terms of the 
GNU General Public License as published by the Free Software Foundation, either version 3 of the License, 
or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied 
warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program.
If not, see <https://www.gnu.org/licenses/>.

Copyright &copy; 2021 Illirgway
