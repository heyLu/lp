(ns shame
  (:use [clojure.data.json :as json :only []]))

; roadmap:
;  1. implement add-item, resurrect-item, close-item
;  2. use refs
;  3. persist to disk
;  4. serve as local web-page
;      caveats: client-side needed (add/close/resurrect, index, note)
;  5. publish

;; todo manipulation

(defn add-item [item shaming]
  "Try adding an item to the shaming. If the maximum number of items is
   exceeeded, return the original shaming."
  (let [c (count (:current shaming))]
    (if (< c (get-in shaming [:config :number-of-items]))
      (assoc-in shaming [:current c] item)
      shaming))) ; FIXME: is that the clojure way of doing it?

(defn get-by [map-coll k v]
  (first (drop-while #(not= (k %) v) map-coll)))

(defn transplant [item association from to]
  "Transplants an item in a map of vectors from one key to another."
  (assoc association
         from (vec (filter #(not= item %) (from association)))
         to   (conj (to association) item)))

(defn close-item [item-name status shaming]
  (let [item (get-by (:current shaming) :name item-name)
        item (assoc item :closed-at (java.util.Date.))]
    (transplant item shaming :current :past)))

(defn resurrect-item [item-name shaming]
  (let [item (get-by (:past shaming) :name item-name)
        item (assoc item :started-at (java.util.Date.))]
    (transplant item shaming :past :current)))

;; references (for "mutating" the todo-list)

(def ^:dynamic *shaming*
  (ref nil))

(defn read-shame [filename]
  (json/read-str (slurp filename) :key-fn keyword))

(defn write-shame [filename shaming]
  (spit filename (json/write-str shaming)))
