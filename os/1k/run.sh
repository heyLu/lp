#!/bin/bash

set -xue

CFLAGS="-std=c11 -O2 -g3 -Wall -Wextra --target=riscv32 -ffreestanding -nostdlib"

clang $CFLAGS -Wl,-Tkernel.ld -Wl,-Map=kernel.map -o kernel.elf \
  kernel.c common.c

qemu-system-riscv32 -machine virt -bios default -nographic -serial mon:stdio --no-reboot \
  -kernel kernel.elf
