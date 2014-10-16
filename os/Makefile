help:
	@echo "make help  - display this message"
	@echo "make build - build the os"
	@echo "make run   - run the os in qemu"
	@echo "make clean - clean up"

build: first.flp

first.flp: first.bin
	dd status=noxfer conv=notrunc if=first.bin of=first.flp

first.bin: first.asm
	nasm -f bin -o first.bin first.asm

run: first.flp
	qemu-system-i386 -fda first.flp

clean:
	rm -f *.flp *.bin mikeos*.zip
