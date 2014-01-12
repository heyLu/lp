(ns clarity.client
  (:use [cljs.reader :only (read-string)])

  (:require [clarity.types :as t]))

(def query-field (js/document.querySelector "#query"))
(def query-result (js/document.querySelector "#query-type-check"))

(defn main []
  (.addEventListener query-field "input"
    (fn [e]
      (let [content (.-value query-field)
            check   (t/friendly-check (read-string content) t/datomic-attr-type)]
        (set! (.-textContent query-result) check)))))
