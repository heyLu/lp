(ns paths
  "paths - Extracting paths from raw OSM data."
  (:use [clojure.xml  :as xml])
  (:use [clojure.repl :as r]))

(defn extract-nodes [osm]
  (let [nodes (filter #(= (:tag %) :node) osm)]
    (into {} (map (fn [node]
                    [(get-in node [:attrs :id]) node])
                  nodes))))

(defn get+
  "Extracts values of `ks` in `m` in the order the `ks` are given."
  [m ks]
  (reduce #(conj %1 (m %2)) [] ks))

(defn tags-as-map [osm-element]
  (let [tags (filter #(= (:tag %) :tag) (:content osm-element))]
    (into {} (map #(-> (get+ (:attrs %) [:k :v])
                       ((fn [[k v]]
                          [(or (keyword k) k) v])))
                  tags))))

(tags-as-map {:content [{:tag :tag, :attrs {:k "hey", :v 1}}, {:tag :tag, :attrs {:k :there, :v 2}}]})

(defn dev-prepare [& [path]]
  (let [path (or path "data/leipzig_highways.osm")]
    (println "; reading in" path "...")
    (def leipzig
      (xml/parse path))
    (println ";=> availlable in `leipzig`")
    (println "; extracting nodes...")
    (def leipzig-nodes
      (extract-nodes {:content leipzig}))
    (println "; => availlable in `leipzig-nodes`")))