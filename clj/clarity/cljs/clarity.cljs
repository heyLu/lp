(ns clarity
  (:require [om.core :as om :include-macros true]
            [om.dom :as dom :include-macros true]))

(enable-console-print!)

(extend-type string
  ICloneable
  (-clone [s] (js/String. s)))

; data for one node: type and optionally data
;  if data is present, fill the node with the data
;  on input, check if input conforms to data, if it
;  does change the data, if it doesn't signal the
;  error and don't change the data

(defn valid? [el]
  (.. el -validity -valid))

(defn read-keyword [str]
  (let [name (.substr str 1)]
    (if (seq name)
      (keyword name)
      nil)))

(defn typed-keyword [kw owner]
  (om/component
    (dom/input #js {:type "text"
                    :className "field"
                    :value (om/value kw)
                    :pattern "^:(\\w+|\\w+(\\.\\w+)*\\/\\w+)$"
                    :onChange (fn [ev]
                                (when (valid? (.-target ev))
                                  (om/update! kw (fn [o n] n) (read-keyword (.. ev -target -value)))))})))

(defn typed-string [string owner]
  (om/component
    (dom/input #js {:type "text"
                    :className "field"
                    :value (om/value string)
                    :onChange #(om/update! string (fn [o n] n) (.. % -target -value))})))

(defn typed-input [data owner]
  (reify
    om/IWillUpdate
    (will-update [_ p s] (prn p s))
    om/IRender
    (render [_]
      (dom/div nil
        (dom/span nil "{")
        (om/build typed-keyword (:kw data))
        (om/build typed-string (:str data))
        (dom/span nil "}")))))

(def app-state
  (atom
    {:kw :hello
     :str "hello"
     :many [1 2 3]}))

(om/root app-state typed-input (.getElementById js/document "typed_input"))