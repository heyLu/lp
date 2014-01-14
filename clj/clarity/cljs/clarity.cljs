(ns clarity
  (:require [om.core :as om :include-macros true]
            [om.dom :as dom :include-macros true]))

(enable-console-print!)

(extend-type string
  ICloneable
  (-clone [s] (js/String. s)))

(extend-type number
  ICloneable
  (-clone [n] (js/Number. n)))

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

(defmulti make-typed-input
  (fn [_ _ {type :type} & _]
    (cond
      (sequential? type) (first type)
      (map? type) (:type type)
      :else type)))

(defmethod make-typed-input 'Number [number owner]
  (om/component
    (dom/input #js {:type "number"
                    :value (om/value number)
                    :onChange #(om/update! number (fn [_ n] n) (js/parseFloat (.. % -target -value)))})))

(defmethod make-typed-input 'Keyword [kw owner]
  (om/component
    (dom/input #js {:type "text"
                    :value (om/value kw)
                    :pattern "^:(\\w+|\\w+(\\.\\w+)*\\/\\w+)$"
                    :onChange (fn [ev]
                                (when (valid? (.-target ev))
                                  (om/update! kw (fn [o n] n) (read-keyword (.. ev -target -value)))))})))

(defmethod make-typed-input 'String [string owner]
  (om/component
    (dom/input #js {:type "text"
                    :value (om/value string)
                    :onChange #(om/update! string (fn [_ n] n) (.. % -target -value))})))

(defmethod make-typed-input 'HMap [m owner {type :type}]
  (om/component
    (dom/div nil
      (dom/span nil "{")
      (into-array
        (map (fn [[k v]]
               (dom/div #js {:className "field"}
                 (dom/label nil (str k))
                 (om/build make-typed-input v {:opts {:type (k (nth type 2))}})))
             m))
      (dom/span nil "}"))))

(def app-state
  (atom
    '{:type (HMap :mandatory {:name String, :age Number, :gender Keyword})
      :data {:name "Paul", :age 3, :gender :unknown}}))

(defn typed-input [{:keys [type data]} owner]
  (reify
    om/IWillUpdate
    (will-update [_ p s] (prn (:data p) s))
    om/IRender
    (render [_]
      (om/build make-typed-input data {:opts {:type type}}))))

(om/root app-state typed-input (.getElementById js/document "typed_input"))