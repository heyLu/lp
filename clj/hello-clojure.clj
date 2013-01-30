(ns hello-clojure
  (:require [clojure.repl :as repl])
  #_(:use zetta.core)
  #_(:require [zetta.combinators :as c])
  #_(:require [zetta.parser.seq :as p]))

(defn human-list
  [l]
  (let [len (count l)]
    (cond (= len 0)
          ""
          (= len 1)
          (str (first l))
          :else
          (apply str
                 (concat (butlast (interleave (map str (butlast l)) (repeat ", ")))
                         [(str " and " (last l))])))))

(human-list [])
(human-list [1])
(human-list [1 2])
(human-list [1 2 3])
(println (human-list (take 11 (repeat "love"))))

(defn same-start [s ss]
  (= (apply str (take (count ss) s)) ss))

(same-start "sorted-map" "foo")
(same-start "sorted-map" "sorted-")

(defn filter-ns
  [filter-fun namespace]
  (let [syms (map str (keys (ns-publics namespace)))]
    (filter filter-fun syms)))

(defn complete
  [to-complete namespace]
  (filter-ns #(.startsWith % to-complete) 'clojure.core))

(defn complete'
  [to-complete namespace]
  (filter-ns #(same-start % to-complete) namespace))

(complete' "st" 'clojure.core)
(print (filter-ns #(.endsWith % "?") 'clojure.core))
(= (complete "take" 'clojure.core) (complete' "take" 'clojure.core))

(def letsample
  '(let [x 10]
     (+ x 3)))

(defn replace-vars
  [[fun bindings body]]
  '(fun bindings
        body))

(replace-vars letsample)

(repl/doc replace)
(replace { 'x 10 } '(+ x 13 [10 x 13]))

(for [x '[1 2 3 4 y]] x) ; i need a for that retains the original seqs

(map inc '(1 2 3))
; already 'good'
(map inc [1 2 3])
(into [] (map inc [1 2 3]))
(map inc { 1 1 2 2 3 3 })
; ???
(map inc #{1 2 3 4 5})
(into #{} (map inc #{1 2 3 4 5}))

(defn map-retain
  [f coll]
  (let [target (cond (vector? coll)
                     []
                     (set? coll)
                     #{}
                     :else
                     nil)]
    (if target
    	(into target (map f coll))
      (map f coll))))

(map-retain inc '(1 2 3))
(map-retain inc [1 2 3])
(map-retain inc #{1 2 3 4 5})

(defn map-retain-rec
  [f coll]
  (map-retain #(if (coll? %)
                 (map-retain-rec f %)
                 (f %))
              coll))

(map-retain-rec inc '(1 2 3 [4 5 6] #{1 2}))
(map-retain-rec str '(let [x 10] (+ x 11)))
(map-retain-rec #(or ({'x 10 'y 42 'z -3} %) %) '(let [x 10] (* y (+ x 1))))
(map-retain-rec #(or ({'coll [1 2 3]} %) %) '(defn map-retain
                                               [f coll]
                                               (let [target (cond (vector? coll)
                                                                  []
                                                                  (set? coll)
                                                                  #{}
                                                                  :else
                                                                  nil)])))

(def all-vars
  (fn all-vars
    [vars body]
    (reduce #(cond (symbol? %2)
                   (conj %1 %2)
                   (coll? %2)
                   (all-vars %1 %2)
                   :else
                   %1)
            vars body)))

(def letsample '(let [[x y z] [1 2 3]
                     {a :a b :b c :c} {:a 4 :b 5 :c 6}]))
(sort (all-vars #{} letsample))

(reduce #(assoc %1 `'~%2 %2) {} (all-vars #{} letsample))
(repl/doc assoc)

(defmacro bind
  [bindings]
  (let [bindings (eval bindings)] ; macro arguments are usually not evaluated?
    															; is there a shorter/more ideomatic way?
    `(let ~bindings
       ~(reduce #(assoc %1 `'~%2 %2)
                {}
                (all-vars #{} bindings)))))

(bind (second letsample))
(bind '[[x y z] [1 2 3] {a :a b :b c :c} {:a 4 :b 5 :c 6}])

(repl/doc binding)
(defn read-file
  [filename]
  (with-open [file (java.io.PushbackReader.
                    (clojure.java.io/reader filename))]
    (binding [*read-eval* false]
      (doall
        (take-while
          #(not= % nil)
          (repeatedly #(read file false nil)))))))

(def hello-forms
  (let [forms (read-file "/home/lu/k/lp/clj/hello-clojure.clj")]
  { :number-of-forms (count forms)
    :forms forms }))
hello-forms

(bind
	(let [human-list-defn (nth (hello-forms :forms) 1)]
    (println (first (drop-while #(not (vector? %)) human-list-defn)))
    (println (first (drop-while #(not (list? %)) human-list-defn)))
  	(into '[l [1 2 3]] (nth (nth human-list-defn 3) 1)))) ; bind in here yields an InstantiationException

(defn apropos+
  "Search a namespace for a symbol whose documentation contains `string`"
  ([string]
   (apropos+ string 'clojure.core))
  ([string namespace]
   (let [symbols (keys (ns-publics namespace))
         docs (map #(-> `#'~% eval meta :doc) symbols)]
     (filter #(and (not= nil %) (.contains % string)) docs))))

(apropos+ "containing" 'clojure.core)
(repl/doc apropos+)

(defn p
  ([x]
   (p x println))
  ([x print-func]
   (print-func x)
   x))

(filter true? (map #(p (= 42 %)) [1 2 3 4 5]))
