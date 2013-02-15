(ns pings
  (:use clojure.repl)
  (:use hello-clojure)
  (:require [datomic.api :as d])
  (:import java.util.Date)
  (:import java.text.SimpleDateFormat))

(.getMethods Date)
(.getMethods java.lang.reflect.Method)

(map #(.getName %) (.getMethods Date))

(def rfc-3339-format (SimpleDateFormat. "yyyy-MM-dd HH:mm:ssXXX"))
(defn parse-date [s]
  "Parse a date (RFC3339 format with seconds). Returns nil if no date can be parsed."
  (try (.parse rfc-3339-format s)
    (catch java.text.ParseException pe nil)))

(parse-date "2013-02-14 15:52:14+01:00")
(doc-ns 'clojure.core)

(def pings
  "Minute-wise 'pings' that are generated when my laptop is running."
  (map #(-> % second parse-date)
       (re-seq #"(.*)\n" (slurp "/home/lu/.pings"))))

(defn diff-seconds [d1 d2]
  (Math/abs (/ (- (.getTime d1) (.getTime d2)) 1000)))

(defn usual-ping? [diff] (< (/ diff 60) 10))

(filter #(-> % usual-ping? not) (map #(->> % (apply diff-seconds)) (partition 2 1 pings)))

(def unusual
  (filter #(not (usual-ping? (nth % 2)))
          (map (fn [[d1 d2]] [d1 d2 (diff-seconds d1 d2)])
               (partition 2 1 pings))))

(doseq [[d1 d2 diff] unusual] (println d1 d2 (float (/ diff 60 60))))

; avg unusual
(let [[c s] (reduce (fn [[c s] [_ _ d]] [(inc c) (+ s d)]) [0 0] unusual)]
  (float (/ s c 60 60)))

(defn avg
  "Calculate the average of the values in coll, choosing values using
   keyfn (defaults to identity)."
  ([coll]
    (avg identity coll))
  ([keyfn coll]
   (let [[c s] (reduce (fn [[c s] v] [(inc c) (+ s (keyfn v))]) [0 0] coll)]
     (/ s c))))

(doseq [[d1 d2 diff] (take 10 (sort-by (fn [[_ _ d]] d) > unusual))]
  (println d1 d2 (float (/ diff 60 60))))

; avg unusual #2
(float (/ (avg #(nth % 2) unusual) 60 60))

;; interesting data (for pings):
;;  * min/max/avg/daily (graph) 'uptime'
;;  * distribution throughout the day
;;  * stretches/pauses/'busy times'

(take 10 (sort-by (fn [[_ _ d]] d) < unusual))
(take 10 (map #(float (/ (nth % 2) 60 60)) (sort-by (fn [[_ _ d]] d) < unusual)))

(def day-format (SimpleDateFormat. "yyyy-MM-dd"))

(def daily-pings (group-by (fn [d] (.format day-format d)) pings))

(defn stretches [pings]
  (first (reduce (fn [[ss start-date last-date] date]
            (if (usual-ping? (diff-seconds last-date date))
              [ss start-date date]
              [(conj ss [start-date last-date]) date date]))
          [[] (first pings) (first pings)]
          (rest pings))))

(defn print-stretches [stretches]
  (doseq [[from to] stretches]
    (println (.format rfc-3339-format from) "to" (.format rfc-3339-format to))))

(print-stretches (stretches pings))
(print-stretches (filter (fn [[from to]] (not= (.getDay from) (.getDay to))) (stretches pings)))

(def daily-stretches (group-by (fn [[d _]] (.format day-format d)) (stretches pings)))
(daily-stretches "2012-12-07")

(defn daily-total [stretches-for-day]
  (reduce (fn [s [f t]] (+ s (diff-seconds f t)))
          0
          stretches-for-day))

(def hours-per-day (map (fn [[_ v]] (float (/ (daily-total v) 60 60))) daily-stretches))

(take 10 (sort-by identity > hours-per-day))
(reduce #(conj %1 %2) #{} (map #(Math/round %) hours-per-day))
(avg hours-per-day)

(def hour-distribution
  (map (fn [[k v]] [k (count v)]) (group-by #(.getHours %) pings)))

(sort-by second > (map (fn [[h c]] [h (float (/ c 60))]) hour-distribution))
