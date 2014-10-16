        BITS 16

start:
        ;; set up 4k stack space after this boot loader (this has something to do with the evil "segment registers")
        mov ax, 0x07c0
        add ax, 288             ; (4096 + 512) / 16 bytes per paragraph (4096 = 4K, 512b for the boot loader, paragraph?)
        mov ss, ax              ; ss = stack space?
        mov sp, 4096            ; sp = stack pointer?

        ;; set data segment to where we're loaded
        mov ax, 0x07c0
        mov ds, ax

        mov si, text_string     ; put string position into si
        call print_string

        ;; jump here, infinite loop!
        jmp $

        text_string db 'This is my very first OS!', 0

print_string:                   ; output string in `si` to screen
        mov ah, 0x0e            ; int 0x10 'print char' function

.repeat:
        ;; get character from string
        lodsb                   ; lodsb = load string byte (loads a byte from `si` and stores it in `al`, which is the lower byte of `ax`)
        cmp al, 0
        je .done                ; if char is zero, we have reached the end of the string
        int 10h                 ; otherwise, print it (e.g. ask the bios to do this)
        jmp .repeat

.done:
        ret

        ;; make this recognizable as a flobby disk boot sector (510b + 0xaa55 at the end = 510b)
        times 510-($-$$) db 0   ; pad remainder of boot sector with 0s
        dw 0xaa55               ; standard pc boot signature (dw = define word, 2 bytes)
