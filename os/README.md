# Writing an OS?

inspired by [Julia Evan's blog post](http://jvns.ca/blog/2013/11/26/day-34-the-tiniest-operating-system/), using the instructions from [MikeOS](http://mikeos.sourceforge.net/).

it clears the screen and then prints some text. exciting!

## How to run it

1. install `nasm` and `qemu`
2. run `make run`
3. have a look at [first.asm](./first.asm) and change something
4. `jmp 2`

## Resources

- interrupts:
    * [interrupt table](http://en.wikipedia.org/wiki/BIOS_interrupt_call#Interrupt_table)
    * [interrupt descriptions](http://www.ctyme.com/intr/int.htm)
- [MikeOS](http://mikeos.sourceforce.net)

## Ideas

- print typed keys to the screen
- port to C
- port to a real language (scheme?!)
- try something graphical
- read [Julia's OS in Rust](https://github.com/jvns/puddle)
