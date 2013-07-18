(ns merge-sort.reference
"Reference implementation of merge sort.

For comparing the different approaches."
  (:refer-clojure :exclude [merge]))

(declare split merge)

(defn merge-sort [xs]
  (if (<= (count xs) 1)
    xs
    (let [[left right] (split xs)]
      (merge (merge-sort left) (merge-sort right)))))

(defn split [xs]
  (let [len (count xs)
        pivot (quot len 2)]
    [(subvec xs 0 pivot) (subvec xs pivot len)]))

(defn merge [left right]
  (loop [xs []
         ls (seq left)
         rs (seq right)]
    (cond
      (and (nil? ls) (nil? rs))
      xs
      (or (nil? ls) (nil? rs))
      (apply conj xs (or ls rs))
      :else
      (let [l (first ls)
            r (first rs)
            [n ls' rs'] (if (<= l r)
                          [l (seq (rest ls)) rs]
                          [r ls (seq (rest rs))])]
        (recur (conj xs n) ls' rs')))))
