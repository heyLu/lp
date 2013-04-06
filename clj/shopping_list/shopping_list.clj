(ns shopping-list
  (:use [clojure.core :as core])
  (:use [clojure.repl :as repl]))

(defn sort-map
  ([m] (sort-map compare m))
  ([c m] (apply sorted-map-by c (reduce #(into %1 %2) (seq m)))))

(def example-article
  (sort-map
    {:all-articles
       [{:name "spinach", :price 3.1, :weight 0.5, :category "vegetables"}
        {:name "apple", :price 2.3, :weight 1.34, :category "vegetables"}
        {:name "onion bread", :price 1.79, :weight 1.1, :category "baked goods"}
        {:name "soy pudding", :price 1.99, :weight 0.5, :category "sweeties"}
        {:name "gum hearts", :price 2.49, :weight 0.4, :category "sweeties"}]}))

(defn create []
  {:known_articles []
   :articles []
   :sort-method :default})

(defn mmap [f m]
  (into (sorted-map) (map (fn [[key val]]
                  [key (f val)])
                m)))

(defn sort-by [key-or-keyfn grouped-articles]
  (mmap #(core/sort-by key-or-keyfn %) grouped-articles))

(def sort-by-category
  (partial sort-by :category))

(sort-by-category {:expensive-things
                     [{:category "a" :n 1}, {:category "b" :n 2} {:category "c" :n 3}]
                   :cheap-things
                     [{:category "d" :n 4}, {:category "a" :n 5}]})
(defn regroup-by [key-or-keyfn grouped-articles]
  (sort-map
    (group-by key-or-keyfn (flatten (vals grouped-articles)))))

(defn first-char-of-name [article]
  (first (:name article)))

(def group-by-name
  (partial mmap #(group-by first-char-of-name)))

(def regroup-by-name
  (partial regroup-by first-char-of-name))

(def group-by-category
  (partial mmap #(group-by :category %)))

(def regroup-by-category
  (partial regroup-by :category))

(defn categorize-by-price [article]
  (let [price (:price article)]
    (condp >= price ; (> x price)
      0.5 "<0.5$"
      1 "<1$"
      2.5 "<2.5$"
      5 "<5$"
      10 "<10$"
      (Integer/MAX_VALUE) ">10$")))

(def group-by-price
  (partial mmap #(group-by categorize-by-price %)))

(def regroup-by-price
  (partial regroup-by categorize-by-price))

(defn categorize-by-weight [article]
  (let [weight (:weight article)]
    (condp >= weight
      0.1 "<100g"
      0.25 "<250g"
      0.5 "<500g"
      1 "<1kg"
      2.5 "<2.5kg"
      5 "<5kg"
      (Integer/MAX_VALUE) ">5kg")))

(def group-by-weight
  (partial mmap #(group-by categorize-by-weight %)))

(def regroup-by-weight
  (partial regroup-by categorize-by-weight))
