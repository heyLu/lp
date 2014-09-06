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

## qst - run things quickly (and easily)

intended to be run in unfamilar environments, might detect the project
type and everything later, for now you pass it a file and it will decide
what to do with it.

- `qst hello_world.go`: compiles and runs `hello_world.go`, rerunning
	after it exits or the file is saved

	quite fun for small things, just throw some code in a file, have `qst`
	watch and restart when appropriate.
- the future: just run `qst` and it will run the first thing it detects