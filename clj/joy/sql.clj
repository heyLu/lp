(ns joy-sql
  "A tiny DSL for SQL.
  
  Note: This code is taken verbatim from The Joy of Clojure."
  (:use [clojure.string :as str :only []]))

(defn expand-expr [expr]
  (if (coll? expr)
    (if (= (first expr) `unquote)
      "?"
      (let [[op & args] expr]
        (str "(" (str/join (str " " op " ")
                           (map expand-expr args))
             ")")))
    expr))

(declare expand-clause)

(def clause-map
  {'SELECT (fn [fields & clauses]
             (apply str "SELECT " (str/join ", " fields)
                    (map expand-clause clauses)))
   'FROM (fn [table & joins]
            (apply str " FROM " table
                   (map expand-clause joins)))
   'LEFT-JOIN (fn [table on expr]
                (str " LEFT JOIN " table
                     " ON " (expand-expr expr)))
   'WHERE (fn [expr]
            (str " WHERE " (expand-expr expr)))})

(defn expand-clause [[op & args]]
  (apply (clause-map op) args))

(defmacro SELECT [& args]
  [(expand-clause (cons 'SELECT args))
   (vec (for [n (tree-seq coll? seq args)
              :when (and (coll? n) (= (first n) `unquote))]
          (second n)))])
