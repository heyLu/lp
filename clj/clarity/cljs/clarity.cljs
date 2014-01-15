(ns clarity
  (:require [cljs.reader :as r]

            [om.core :as om :include-macros true]
            [om.dom :as dom :include-macros true]))


(enable-console-print!)

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

(defmethod empty-value 'Option [[_ v]]
  (empty-value v))

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

(defmethod make-typed-input 'Boolean [m owner {:keys [type key val]}]
  (om/component
    (dom/input #js {:type "checkbox"
                    :checked (.valueOf (or val (empty-value type)))
                    :onChange #(om/transact! m key (fn [_ n] n) (js/Boolean. (.. % -target -checked)))})))

(defmethod make-typed-input 'Number [m owner {:keys [key val]}]
  (om/component
    (dom/input #js {:type "number"
                    :value val
                    :onChange (update-on-change! m key js/parseFloat)})))

(def keyword-pattern "^:(\\w+|\\w+(\\.\\w+)*\\/\\w+)$")

(defmethod make-typed-input 'Keyword [m owner {:keys [key val]}]
  (om/component
    (dom/input #js {:type "text"
                    :value val
                    :pattern keyword-pattern
                    :onChange (update-on-change! m key #(or (read-keyword %) (empty-value type)))})))

(defn update-on-change!
  ([m k] (update-on-change! m k identity))
  ([m k transform-fn] (update-on-change! m k transform-fn false))
  ([m k transform-fn optional?]
   (fn [ev]
     (let [new-val (transform-fn (.. ev -target -value))]
       (if (and optional? (empty? new-val))
         (om/update! m dissoc k)
         (om/transact! m k (fn [_ n] n) new-val))))))

(defmethod make-typed-input 'String [m owner {:keys [type key val optional?]}]
  (om/component
    (dom/input #js {:type "text"
                    :value val
                    :onChange (update-on-change! m key identity optional?)})))

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

(defmethod make-typed-input 'Option [m owner {type :type :as opts}]
  (let [[_ type] type]
    (make-typed-input m owner (assoc opts :type type))))

(defmethod make-typed-input 'U [m owner {:keys [type key val]}]
  (om/component
    (dom/select #js {:onChange (update-on-change! m key r/read-string)}
      (into-array
        (map (fn [[_ v]]
               (dom/option nil (str v)))
             (rest type))))))

(defmethod make-typed-input 'HMap [m owner {type :type}]
  (let [hmap (apply hash-map (rest type))
        required (:mandatory hmap)
        optional (:optional hmap)]
    (om/component
      (dom/div nil
        (dom/span nil "{")
        (into-array
          (map (fn [[k t]]
                 (dom/div #js {:className "field"}
                   (dom/label nil (str k))
                   (om/build make-typed-input m {:opts {:type t, :key k, :val (k m)
                                                        :optional? (contains? optional k)}})))
               (merge required optional)))
        (dom/span nil "}")))))

(def app-state
  (let [type '[HMap :mandatory
                    {:name {:type String :default "Paul"},
                     :age {:type Number, :default 10},
                     :language (U (Value :en)
                                  (Value :de)
                                  (Value :fr)
                                  (Value :jp))
                     :fun Boolean
                     :gender Keyword}
                    :optional
                    {:secret-skill String}]]
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