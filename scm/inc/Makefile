run-scheme: scheme
	./scheme

scheme: scheme.s driver.c
	gcc -O3 scheme.s driver.c -o $@

scheme.s: compiler.scm
	petite --script compiler.scm > $@

integers.s: integers.c
	gcc -O3 --omit-frame-pointer -S $<
