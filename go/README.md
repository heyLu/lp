# go

Playing with [go](http://golang.org). Late to the party, but it's fun,
I think.

## thoughts

- *fast*
- some level of *type-safety* (just scratched the surface so far)
- good *tool support* (fast (!) compilation, the `go` tool itself,
	fetching libraries built-in, though versioning is missing)
- *simple* (mostly, goroutines + no proper sync will bite you,
	thinking helps, as always)
- *pragmatic*, possibly similar to clojure in this respect (which is
    somewhat surprising to me)
- misc
    * [cross compilation][] is surprisingly painless

[cross compilation]: http://dave.cheney.net/2012/09/08/an-introduction-to-cross-compilation-with-go

## qst - run things quickly (and easily)

`qst` has already grown up, it now lives [in it's own place](https://github.com/heyLu/qst).
You can get it using `go get github.com/heyLu/qst`.

intended to be run in unfamilar environments, you pass it a file or a
directory and it tries to detect what it is and how to run it.

run `qst .` to run anything.
