; let's have a bit of clojure in scheme

(define nil #f)

(define ex-m '((x . 10) (y . 13) (z . -37)))

(define (clj-assoc m k v)
  (if (null? m)
    `((,k . ,v))
    (let ((e (car m)))
      (if (equal? (car e) k)
        `((,k . ,v) . ,(cdr m))
        `(,e . ,(clj-assoc (cdr m) k v))))))

(clj-assoc ex-m 'a 10)
(clj-assoc ex-m 'y 3)

(define (clj-dissoc m k)
  (if (null? m)
    m
    (let ((e (car m)))
      (if (equal? (car e) k)
        (cdr m)
        `(,e . ,(clj-dissoc (cdr m) k))))))

(clj-dissoc ex-m 'z)
(clj-dissoc ex-m 'a)

(define (clj-get m k)
  (if (null? m)
    nil
    (let ((e (car m)))
      (if (equal? (car e) k)
        (cdr e)
        (clj-get (cdr m) k)))))

(clj-get ex-m 'z)
(clj-get ex-m 'a)

(define (clj-get-in m ks)
  (cond
    ((null? m) nil)
    ((null? ks) m)
    (else (let ((s (clj-get m (car ks))))
            (if s
              (clj-get-in s (cdr ks))
              nil)))))

(define ex-nested-m `((a . ,ex-m) (b . ((c . ((d . ((e . ,ex-m) (f . 42)))))))))

(clj-get-in ex-nested-m '(a z))
(clj-get-in ex-nested-m '(b c d f))
(clj-get-in ex-nested-m '(b x y z))