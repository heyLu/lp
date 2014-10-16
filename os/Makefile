help:
	@echo "make help  - display this message"
	@echo "make build - build the os"
	@echo "make run   - run the os in qemu"
	@echo "make clean - clean up"

build: first.flp

first.flp: mikeos.flp first.bin
	cp mikeos.flp first.flp
	dd status=noxfer conv=notrunc if=first.bin of=first.flp

first.bin: first.asm
	nasm -f bin -o first.bin first.asm

mikeos.flp:
	curl -so mikeos.zip http://mikeos.sourceforge.net/mikeos-4.5b1.zip
	unzip -p mikeos.zip mikeos-4.5b1/disk_images/mikeos.flp > mikeos.flp
	rm mikeos.zip

run: first.flp
	qemu-system-i386 -fda first.flp

clean:
	rm -f *.flp *.bin mikeos*.zip
