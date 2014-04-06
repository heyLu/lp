CFLAGS=-std=gnu99 -Iluajit

bf: bf.c dynasm-driver.c
	$(CC) $(CFLAGS) -o bf dynasm-driver.c -DJIT=\"bf.c\"

mandelbrot.bf:
	curl -L -O https://github.com/pablojorge/brainfuck/raw/master/samples/mandelbrot.bf

hello_dynasm: hello_dynasm.c dynasm-driver.c
	$(CC) $(CFLAGS) -o hello_dynasm dynasm-driver.c -DJIT=\"hello_dynasm.c\"

%.c: %.dasc
	./dynasm -o $@ $<

%.dasc: dynasm

dynasm-driver.c:
	curl -O https://raw.githubusercontent.com/haberman/jitdemo/master/dynasm-driver.c

dynasm: luajit

luajit:
	git clone http://luajit.org/git/luajit-2.0.git luajit
