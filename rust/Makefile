diff_code: diff_code.c
	gcc -std=c99 -o diff_code diff_code.c

keylog: keylog.rs libkeylog_xlib.so
	rustc -L rust-xlib -L . keylog.rs

libkeylog_xlib.so: keylog_xlib.c
	clang -fPIC -c keylog_xlib.c
	clang -shared -Wl,-soname,libkeylog_xlib.so -o libkeylog_xlib.so keylog_xlib.o

clean-keylog:
	rm -f keylog libkeylog_xlib.so keylog_xlib.o

clean: clean-keylog
