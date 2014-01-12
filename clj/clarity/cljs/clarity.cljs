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
  (fn [type & _]
    (if (seq? type)
      (first type)
      type)))

(defn change-data
  ([ev owner]
   (om/set-state! owner :data (.. ev -target -value)))
  ([ev owner pred parse]
   (when (pred (.-target ev))
     (om/set-state! owner :data (parse (.. ev -target -value))))))

(defn valid? [el]
  (.. el -validity -valid))

(defn read-keyword [str]
  (let [name (.substr str 1)]
    (if (seq name)
      (keyword name)
      nil)))

(defmethod make-typed-input 'Keyword [type data owner]
  (reify
    om/IInitState
    (init-state [_] {:data data})
    om/IRender
    (render [_]
      (dom/input #js {:type text
                      :placeholder ":my.ns/identifier"
                      :pattern "^:(\\w+|\\w+(\\.\\w+)*\\/\\w+)$"
                      :onChange #(change-data % owner valid? read-keyword)}))))

(defmethod make-typed-input 'String [type data owner]
  (reify
    om/IInitState
    (init-state [_] {:data data})
    om/IRender
    (render [_]
      (dom/input #js {:type "text"
                      :value (om/get-state owner :data)
                      :onChange #(change-data % owner)}))))

(defn typed-input [typed-data owner]
  (let [{:keys [type data]} typed-data]
    (make-typed-input type data owner)))

(om/root (atom {:type 'Keyword}) typed-input (.getElementById js/document "typed_input"))