(ns macros
  "Dilettante first steps with macros")

(defmacro my-cond [& clauses]
  (when clauses
    `(if ~(ffirst clauses)
       ~(second (first clauses))
       (my-cond ~@(next clauses)))))
