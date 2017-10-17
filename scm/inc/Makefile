run-scheme: scheme
	./scheme

scheme: scheme.s driver.c
	gcc -O3 -g scheme.s driver.c -o $@

scheme-static: scheme.s driver.c
	gcc -O3 -static scheme.s driver.c -o $@

asm: scheme.s
	cat scheme.s

scheme.s: compiler.scm
	petite --script compiler.scm > $@

integers.s: integers.c
	gcc -O3 --omit-frame-pointer -S $<
