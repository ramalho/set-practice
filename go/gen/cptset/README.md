[![GoDoc](https://godoc.org/github.com/standupdev/set-practice/go/gen/cptset?status.svg)](https://godoc.org/github.com/standupdev/set-practice/go/gen/cptset)

# `CodepointSet`: a set of runes built with `gen`

Install `gen`:

```
$ go get github.com/clipperhouse/gen
```

Add the set _typewriter_ to `gen`:

```
$ gen add github.com/clipperhouse/set 
```

Verify that the `set` _typewriter_ is installed:

```
$ gen list
Imported typewriters:
  slice
  stringer
  set

```

Use `gen` to generate a `codepoint_set.go` file:

```
$ gen
$ ls
codepoint_set.go  codepoints.go  _gen.go
```

Run tests:

```
$ go test
```

Check out the examples in `codepoint_test.go`.
