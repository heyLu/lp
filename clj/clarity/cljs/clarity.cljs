(ns clarity
  (:require [cljs.reader :as r]

            [om.core :as om :include-macros true]
            [om.dom :as dom :include-macros true]))

(enable-console-print!)

(extend-type js/Boolean
  ICloneable
  (-clone [b] b))

(extend-type number
  ICloneable
  (-clone [n] (js/Number. n)))

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

(defmulti empty-value
  (fn get-type [type]
    (cond
      (sequential? type) (first type)
      (map? type) (get-type (:type type))
      :else type)))

(defmethod empty-value 'Boolean [{:keys [default]}]
  (js/Boolean.
    (if-not (nil? default)
      default
      false)))

(defmethod empty-value 'Number [{:keys [default]}]
  (or default 0))

(defmethod empty-value 'Keyword [{:keys [default]}]
  (or default :keyword))

(defmethod empty-value 'String [{:keys [default]}]
  (or default ""))

(defmethod empty-value 'Value [[_ v]]
  v)

(defmethod empty-value 'U [[_ & [[_ v]]]]
  v)

(defmethod empty-value 'HMap [spec]
  (let [entries (nth spec 2)]
    (into {}
          (map (fn [[k v]]
                 [k (empty-value v)])
               entries))))

(defmulti make-typed-input
  (fn [_ _ {type :type} & _]
    (cond
      (sequential? type) (first type)
      (map? type) (:type type)
      :else type)))

(defmethod make-typed-input 'Boolean [bool owner]
  (om/component
    (dom/input #js {:type "checkbox"
                    :checked (.valueOf bool)
                    :onChange #(om/update! bool (fn [_ n] n) (js/Boolean. (.. % -target -checked)))})))

(defmethod make-typed-input 'Number [number owner]
  (om/component
    (dom/input #js {:type "number"
                    :value (om/value number)
                    :onChange #(om/update! number (fn [_ n] n) (js/parseFloat (.. % -target -value)))})))

(def keyword-pattern "^:(\\w+|\\w+(\\.\\w+)*\\/\\w+)$")

(defmethod make-typed-input 'Keyword [kw owner]
  (om/component
    (dom/input #js {:type "text"
                    :value (om/value kw)
                    :pattern keyword-pattern
                    :onChange (fn [ev]
                                (when (valid? (.-target ev))
                                  (om/update! kw (fn [o n] n) (or (read-keyword (.. ev -target -value))
                                                                  (empty-value 'Keyword)))))})))

(defmethod make-typed-input 'String [string owner]
  (om/component
    (dom/input #js {:type "text"
                    :value (om/value string)
                    :onChange #(om/update! string (fn [_ n] n) (.. % -target -value))})))

(defmethod make-typed-input 'Value [value owner]
  (om/component
   (dom/input (clj->js
                (into {:value (str value)
                       :readOnly "readOnly"}
                  (cond
                    (instance? js/Boolean value) {:type "checkbox", :checked value}
                    (number? value) {:type "number"}
                    (keyword? value) {:type "text", :pattern keyword-pattern}
                    :else {:type "text"}))))))

(defmethod make-typed-input 'U [value owner {type :type}]
  (om/component
    (dom/select #js {:onChange #(om/update! value (fn [_ n] n) (r/read-string (.. % -target -value)))}
      (into-array
        (map (fn [[_ v]]
               (dom/option nil (str v)))
             (rest type))))))

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
  (let [type '(HMap :mandatory
                    {:name {:type String :default "Paul"},
                     :age {:type Number, :default 10},
                     :language (U (Value :en)
                                  (Value :de)
                                  (Value :fr)
                                  (Value :jp))
                     :fun Boolean
                     :gender Keyword})]
    (atom
     {:type type
      :data (empty-value type)})))

(defn typed-input [{:keys [type data]} owner]
  (reify
    om/IWillUpdate
    (will-update [_ p s] (prn (:data p) s))
    om/IRender
    (render [_]
      (om/build make-typed-input data {:opts {:type type}}))))

(om/root app-state typed-input (.getElementById js/document "typed_input"))