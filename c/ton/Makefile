JSC_CFLAGS = $(shell pkg-config javascriptcoregtk-4.0 --cflags --libs)

ton: *.c
	clang $(JSC_CFLAGS) main.c -o $@
