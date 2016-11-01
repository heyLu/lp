DESTDIR = $(HOME)/.local/bin

RUNNERS = $(wildcard bin/run-*)

all: deps

deps: lib/codemirror.js lib/codemirror.css

lib:
	mkdir lib

lib/codemirror.js: lib
	curl -o $@ 'http://marijnhaverbeke.nl/uglifyjs' --data 'download=codemirror-compressed.js' \
		--data 'code_url=http://codemirror.net/lib/codemirror.js' \
		--data 'code_url=http://codemirror.net/mode/clike/clike.js'\
		--data 'code_url=http://codemirror.net/mode/clojure/clojure.js'\
		--data 'code_url=http://codemirror.net/mode/go/go.js'\
		--data 'code_url=http://codemirror.net/mode/haskell/haskell.js'\
		--data 'code_url=http://codemirror.net/mode/javascript/javascript.js'\
		--data 'code_url=http://codemirror.net/mode/julia/julia.js'\
		--data 'code_url=http://codemirror.net/mode/lua/lua.js'\
		--data 'code_url=http://codemirror.net/mode/python/python.js'\
		--data 'code_url=http://codemirror.net/mode/ruby/ruby.js'\
		--data 'code_url=http://codemirror.net/mode/rust/rust.js'\
		--data 'code_url=http://codemirror.net/mode/shell/shell.js'\
		--data 'header=%2F*+CodeMirror+-+Minified+%26+Bundled+*%2F'

lib/codemirror.css: lib
	curl -o $@ http://codemirror.net/lib/codemirror.css

linguaevalia: *.go
	go build .

install: linguaevalia
	install -t $(DESTDIR) linguaevalia $(RUNNERS)

clean:
	rm -rf lib
