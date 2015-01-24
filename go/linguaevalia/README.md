# Lingua evalia

A tiny web service that runs code for you. You open it, start writing code
immediately and then run it. What you don't do is worry about file names,
which command it was you had to use to run the code and so on.

Just write the code.

## Quickstart

* `go get github.com/heyLu/lp/go/linguaevalia`
* `$GOPATH/bin/linguaevalia`
* visit <http://localhost:8000> and start writing code (press `ctrl-enter`
    to run the code)

Essentially just need `go`, but to run code in other languages, you need
to have them installed as well.

## Languages

- `go`
- `python`
- `ruby`
- `javascript`
- `haskell`
- `rust`
- `julia`
- `pixie`

Adding more is relatively simple: If there is a command that runs code in
a language given a file, just add [the appropriate line](./linguaevalia.go#L40-L47)
and a corresponding mapping to `languageMappings`.

If there isn't, you can either write a wrapper to do that (similar to the
[one for rust](./bin/run-rust)) or you can implement the `Language`
interface.

## Contributions and feedback welcome!

Tell me what you do with it, when it helped you, what you're missing.

Have fun!

## License

[MIT](./LICENSE)
