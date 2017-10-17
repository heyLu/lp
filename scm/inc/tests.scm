(load "compiler.scm")

(define (do-test name input output)
  (let [[test-file "current-test"]
        [test-src  "current-test.s"]
        [program (with-input-from-string input read)]]
  (delete-file test-src)
  (with-output-to-file test-src
    (lambda ()
      (compile-program program)))
  (unless (zero? (system (format "gcc -O3 ~a driver.c -o ~a" test-src test-file)))
    (error 'gcc "could not compile"))
  (unless (zero? (system (format "./~a > tests-output" test-file)))
    (error 'run "could not run"))
  (let [[out (read-file "tests-output")]]
    (unless (string=? out output)
      (error 'test (format "test[~a]: ~s: Expected ~s, but got ~s" name input output out))))
  (delete-file "tests-output")
  (delete-file test-src)
  (delete-file test-file)
  (display (format "test[~a]: OK\t\t~s\n" name input))))

(define (read-file path)
  (with-output-to-string
    (lambda ()
      (with-input-from-file path
        (lambda ()
          (let recur ()
            (let [[c (read-char)]]
              (unless (eof-object? c)
                (display c)
                (recur)))))))))

(define (do-tests name tests)
  (unless (null? tests)
    (do-test name (caar tests) (cadar tests))
    (do-tests name (cdr tests))))

(do-tests "3.1 - integers"
  '[["42" "42\n"]
    ["10" "10\n"]
    ["-1001" "-1001\n"]])

(do-tests "3.2 - immediate constants"
  '[["#\\d" "#\\d\n"]
    ["#\\y" "#\\y\n"]
    ["#t" "#t\n"]
    ["#f" "#f\n"]
    ["()" "()\n"]])

(do-tests "3.3 - unary primitives"
  '[["(add1 0)" "1\n"]
    ["(add1 41)" "42\n"]
    ["(add1 -131124)" "-131123\n"]

    ["(integer->char 121)" "#\\y\n"]
    ["(char->integer #\\y)" "121\n"]

    ["(zero? 0)" "#t\n"]
    ["(zero? 1)" "#f\n"]
    ["(zero? 21425)" "#f\n"]
    ["(zero? -142)" "#f\n"]
    ["(zero? ())" "#f\n"]
    ["(zero? #f)" "#f\n"]
    ["(zero? #\\x)" "#f\n"]

    ["(null? ())" "#t\n"]
    ["(null? 0)" "#f\n"]
    ["(null? #\\y)" "#f\n"]
    ["(null? 13)" "#f\n"]
    ["(null? #t)" "#f\n"]
    ["(null? #f)" "#f\n"]

    ["(integer? 0)" "#t\n"]
    ["(integer? 13)" "#t\n"]
    ["(integer? -1)" "#t\n"]
    ["(integer? 15325232)" "#t\n"]
    ["(integer? -125252121)" "#t\n"]
    ["(integer? #\\y)" "#f\n"]
    ["(integer? #t)" "#f\n"]
    ["(integer? #f)" "#f\n"]
    ["(integer? ())" "#f\n"]

    ["(boolean? #t)" "#t\n"]
    ["(boolean? #f)" "#t\n"]
    ["(boolean? 0)" "#f\n"]
    ["(boolean? 12421)" "#f\n"]
    ["(boolean? #\\y)" "#f\n"]
    ["(boolean? ())" "#f\n"]
    ])
