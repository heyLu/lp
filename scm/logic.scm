(load "/home/lu/t/TheReasonedSchemer/mk.scm")
(load "/home/lu/t/TheReasonedSchemer/mkextraforms.scm")

(define peano-o
  (lambda (x)
    (conde
     ((== x 'z) succeed)
     (else (fresh (y)
                  (== x (list 's y))
                  (peano-o y))))))

(define plus-o
  (lambda (x y z)
    (conde
     ((== x 'z) (== y z))
     (else (fresh (sx sz)
                  (== x (list 's sx))
                  (== z (list 's sz))
                  (plus-o sx y sz))))))

(define cons-o
  (lambda (x l r)
    (== r (cons x l))))

(define append-o
  (lambda (l s r)
    (conde
     ((== l '()) (== s r))
     ((fresh (a d res)
             (== `(,a . ,d) l)
             (== `(,a . ,res) r)
             (append-o d s res))))))
