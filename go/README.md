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

intended to be run in unfamilar environments, you pass it a file or a
directory and it tries to detect what it is and how to run it.

- run `qst` or `qst -h` to see options and support project types
- `qst hello_world.go`: compiles and runs `hello_world.go`, rerunning
	after it exits or the file is saved

	quite fun for small things, just throw some code in a file, have `qst`
	watch and restart when appropriate.
- `qst -phase=test ...` runs the tests for projects that support it

### Building it yourself

	# set up $GOPATH as desired
	$ export GOPATH=$PWD/.go         # choose whatever you want
	$ go build qst
	...
	$ ./qst -h
	Usage: qst <file>
	...
	$ ./qst examples/hello_web.rb
	...
	^C
	$ ./qst examples/hello_web.go
	...

Try changing something in the files, it's fun. :)