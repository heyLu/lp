(load "compiler.scm")

(compile-program '(let ((x 3) (y 4)) (if (integer? (+ x y)) 42 #f)))
