;; assembly resources:
;;  - https://en.wikipedia.org/wiki/X86_instruction_listings
(define (emit instr . args)
  (display "\t")
  (display (apply format instr args))
  (display "\n"))

(define wordsize 4)

(define fixnum-shift 2)
(define char-shift 8)
(define boolean-shift 7)

(define (immediate-rep x)
  (cond
    ((integer? x) (bitwise-arithmetic-shift x fixnum-shift))
    ((char? x) (bitwise-xor (bitwise-arithmetic-shift (char->integer x) char-shift) #b00001111))
    ((boolean? x) (bitwise-xor (bitwise-arithmetic-shift (cond
                                                           ((boolean=? x #t) 1)
                                                           ((boolean=? x #f) 0))
                                                         boolean-shift)
                               #b0011111))
    ((eq? x '()) #b00101111)))

(define (immediate? x)
  (or (integer? x) (char? x) (boolean? x) (eq? x '())))

(define (primcall? x)
  (list? x))

(define (primcall-op x)
  (car x))

(define (primcall-operand1 x)
  (cadr x))

(define (primcall-operand2 x)
  (caddr x))

(define (emit-compare)
  (emit "movl $0,  %eax")                ; zero %eax to put the result of the comparison into
  (emit "sete %al")                      ; set low byte of %eax to 1 if cmp succeeded
  (emit "sall $~a,  %eax" boolean-shift) ; construct correctly tagged boolean value
  (emit "xorl $~a, %eax" #b0011111))

(define (emit-expr x si)
  (cond
    ((immediate? x)
     (emit "movl $~a, %eax" (immediate-rep x)))
    ((primcall? x) (emit-primitive-call x si))))

(define (emit-primitive-call x si)
  (case (primcall-op x)
    ((add1)
     (emit-expr (primcall-operand1 x) si)
     (emit "addl $~a, %eax" (immediate-rep 1)))
    ((integer->char)
     (emit-expr (primcall-operand1 x) si)
     (emit "shl $6, %eax")
     (emit "xorl $15, %eax"))
    ((char->integer)
     (emit-expr (primcall-operand1 x) si)
     (emit "shrl $6, %eax"))
    ((zero?)
     (emit-expr (primcall-operand1 x) si)
     (emit "cmpl $0,  %eax") ; x == 0
     (emit-compare))
    ((null?)
     (emit-expr (primcall-operand1 x) si)
     (emit "cmpl $~a, %eax" #b00101111)
     (emit-compare))
    ((integer?)
     (emit-expr (primcall-operand1 x) si)
     (emit "andl $~a, %eax" #b11)
     (emit-compare))
    ((char?)
     (emit-expr (primcall-operand1 x) si)
     (emit "andl $~a, %eax" #b11111111)
     (emit "cmpl $~a, %eax" #b00001111)
     (emit-compare))
    ((boolean?)
     (emit-expr (primcall-operand1 x) si)
     (emit "andl $~a, %eax" #b1111111)
     (emit "cmpl $~a, %eax" #b0011111)
     (emit-compare))
    ((+)
     (emit-expr (primcall-operand2 x) si)
     (emit "movl %eax, ~a(%rsp)" si) ; move second arg on the stack
     (emit-expr (primcall-operand1 x) (- si wordsize))
     (emit "addl ~a(%rsp), %eax" si))
    ))

(define (compile-program x)
  (display ".globl scheme_entry\n\n")
  (display "scheme_entry:\n")

  (emit-expr x (- wordsize))
  (emit "ret"))
