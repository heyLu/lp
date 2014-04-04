CC=gcc -std=c99

run: mul
	./mul

mul: main.c mpc.c mpc.h
	$(CC) -ledit -lm mpc.c main.c -o mul

mpc_base_url=https://raw.githubusercontent.com/orangeduck/mpc/a25dc955a

mpc.h:
	curl -O $(mpc_base_url)/mpc.h

mpc.c:
	curl -O $(mpc_base_url)/mpc.c

clean:
	rm -f mul
