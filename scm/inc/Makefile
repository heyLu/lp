CC=gcc
CFLAGS=-Wall -g -O3

run-scheme: scheme
	./scheme

scheme: scheme.s driver.c
	$(CC) $(CFLAGS) scheme.s driver.c -o $@

scheme-static: scheme.s driver.c
	$(CC) $(CFLAGS) -static scheme.s driver.c -o $@

test: scheme.s driver.c
	petite --script tests.scm

asm: scheme.s
	cat scheme.s

scheme.s: compiler.scm program.scm
	petite --script program.scm > $@

integers.s: integers.c
	gcc -O3 --omit-frame-pointer -S $<
