JSC_CFLAGS = $(shell pkg-config javascriptcoregtk-4.0 --cflags --libs)

ton: *.c
	clang $(JSC_CFLAGS) $^ -o $@

jsc-funcs:
	grep -Rh JS_EXPORT /usr/include/webkitgtk-4.0/JavaScriptCore | sed 's/^JS_EXPORT //' | grep -v '^#' > $@
