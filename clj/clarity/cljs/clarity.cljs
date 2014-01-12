(ns clarity
  (:require [om.core :as om :include-macros true]
            [om.dom :as dom :include-macros true]))

(enable-console-print!)

; data for one node: type and optionally data
;  if data is present, fill the node with the data
;  on input, check if input conforms to data, if it
;  does change the data, if it doesn't signal the
;  error and don't change the data

(def example-state
  (atom
   {:type 'String
    :data "hello"}))

(defmulti make-typed-input
  (fn [{type :type} & _]
    (if (seq? type)
      (first type)
      type)))

(defn change-data
  ([ev d]
   (om/update! d assoc :data (.. ev -target -value)))
  ([ev d pred parse]
   (when (pred (.-target ev))
     (om/update! d assoc :data (parse (.. ev -target -value))))))

(defn valid? [el]
  (.. el -validity -valid))

(defn read-keyword [str]
  (let [name (.substr str 1)]
    (if (seq name)
      (keyword name)
      nil)))

(defmethod make-typed-input 'Keyword [d owner]
  (reify
    om/IRender
    (render [_]
      (dom/input #js {:type text
                      :value (om/read d :data)
                      :placeholder ":my.ns/identifier"
                      :pattern "^:(\\w+|\\w+(\\.\\w+)*\\/\\w+)$"
                      :onChange #(change-data % d valid? read-keyword)}))))

(defmethod make-typed-input 'String [d owner]
  (reify
    om/IRender
    (render [_]
      (dom/input #js {:type "text"
                      :value (om/read d :data)
                      :onChange #(change-data % d)}))))

(om/root (atom {:type 'Keyword}) make-typed-input (.getElementById js/document "typed_input"))