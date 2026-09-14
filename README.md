[![CI](https://github.com/dangogh/naturally/actions/workflows/ci.yml/badge.svg)](https://github.com/dangogh/naturally/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/dangogh/naturally.svg)](https://pkg.go.dev/github.com/dangogh/naturally)

# naturally

Go implementation of a *natural* sort. Numbers and non-numbers are treated as separate sub-parts to
make the sort more natural for human readability. Inspired by the Perl implementation Sort::Naturally.

## Install

```sh
go get github.com/dangogh/naturally
```

## Usage

```go
files := []string{"photo10.jpg", "photo9.jpg", "photo1.jpg"}

naturally.Sort(files)
// [photo1.jpg photo9.jpg photo10.jpg]
// a plain sort.Strings would give [photo1.jpg photo10.jpg photo9.jpg]

naturally.SortCI(files) // same, ignoring case
```

For use with `sort.Sort` / `sort.Stable` directly, the `StringSlice` and
`CIStringSlice` types implement `sort.Interface`:

```go
sort.Stable(naturally.StringSlice(files))
sort.Stable(naturally.CIStringSlice(files))
```

## Ordering rules

- Numeric segments compare by value, with no digit limit, so
  `"99999999999999999999"` sorts before `"100000000000000000000"`.
- Equal values with different widths order by the shorter run first:
  `"A1" < "A01" < "A001"`.
- A string containing a digit sorts before one that does not, so
  `"b2" < "Banana"`. The alphabetic prefixes are not compared when only one
  side has digits.
- Segments of non-ASCII digits (for example Arabic-Indic `٣`) are compared as
  text rather than by value.
- `SortCI` uses simple Unicode case folding. Values that differ only by case
  compare equal, so use `sort.Stable` if you need their relative order
  preserved.
