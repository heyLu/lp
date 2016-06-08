JSC_CFLAGS = $(shell pkg-config javascriptcoregtk-4.0 --cflags --libs)
LIBZIP_CFLAGS = $(shell pkg-config libzip --cflags --libs)

ton: *.c
	clang $(JSC_CFLAGS) $(LIBZIP_CFLAGS) $^ -o $@

jsc-funcs:
	grep -Rh JS_EXPORT /usr/include/webkitgtk-4.0/JavaScriptCore | sed 's/^JS_EXPORT //' | grep -v '^#' > $@
