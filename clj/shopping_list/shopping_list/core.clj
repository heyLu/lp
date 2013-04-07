(ns shopping-list.core)

(defn map-map [f m]
  (into (empty m) (map (fn [[key val]]
                  [key (f val)])
                m)))

(defn sort-map
  ([m] (sort-map compare m))
  ([c m] (apply sorted-map-by c (reduce #(into %1 %2) (seq m)))))

(def default
  {:group (constantly :all-articles)
   :sort  :name})

(def category
  {:group :category
   :sort  :category})

(defn categorize-number [categories n]
  (last (filter #(< % n) (concat [0] categories [Double/POSITIVE_INFINITY]))))

(defrecord PriceInEuros [price-in-cents]
  Comparable
  (compareTo [this other-price]
             (compare price-in-cents (:price-in-cents other-price)))
  Object
  (toString [this]
            (if (< price-in-cents 100)
                   (str price-in-cents "ct")
                   (str (float (/ price-in-cents 100)) "€"))))
(defn pretty-price [price-in-cents] (str (PriceInEuros. price-in-cents)))

(def price
  {:group (comp #(PriceInEuros. %) #(categorize-number [50 100 250 500 1000] %) :price)
   :sort  :price})

(defrecord MetricWeight [weight-in-grams]
  Comparable
  (compareTo [this other-weight]
             (compare weight-in-grams (:weight-in-grams other-weight)))
  Object
  (toString [this]
            (if (< weight-in-grams 1000)
              (str weight-in-grams "g")
              (str (float (/ weight-in-grams 1000)) "kg"))))
(defn pretty-weight [price-in-grams] (str (MetricWeight. price-in-grams)))

(def weight
  {:group (comp #(MetricWeight. %) #(categorize-number [100 250 500 1000 2500 5000] %) :weight)
   :sort  :weight})

(def alphabet
  {:group (comp first :name)
   :sort  :name})

(defn ungroup [grouped-articles]
  (flatten (vals grouped-articles)))

(defn regroup [{:keys [group sort reverse-group reverse-sort]} grouped-articles]
  (let [articles   (ungroup grouped-articles)
        compare-fn (fn [reverse?] (if reverse? (comp - compare) compare))]
    (map-map #(sort-by sort (compare-fn reverse-sort) %)
             (sort-map (compare-fn reverse-group) (group-by group articles)))))

(defn combine [{group :group} {sort :sort}]
  {:group group
   :sort  sort})

(def example-articles
  (regroup default
    {:all-articles
       [{:name "spinach", :price 310, :weight 500, :category "vegetables"}
        {:name "apple", :price 230, :weight 1340, :category "vegetables"}
        {:name "onion bread", :price 179, :weight 1100, :category "baked goods"}
        {:name "soy pudding", :price 199, :weight 500, :category "sweeties"}
        {:name "gum hearts", :price 249, :weight 400, :category "sweeties"}]}))

(regroup (assoc alphabet :reverse-group false :reverse-sort true) example-articles)
(map str (keys (regroup price example-articles)))