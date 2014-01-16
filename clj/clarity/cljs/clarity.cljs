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
  (if-not (nil? default)
    default
    false))

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

(defn update-on-change! [m k {:keys [extract-fn transform-fn optional? type]
                              :or {extract-fn #(.-value %)
                                   transform-fn identity
                                   type 'String}}]
  (fn [ev]
    (let [new-val (transform-fn (extract-fn (.-target ev)))]
      (if (and optional? (= new-val (empty-value type)))
        (om/update! m dissoc k)
        (om/transact! m k (fn [_ n] n) new-val)))))

(defmulti make-typed-input
  (fn [_ _ {type :type} & _]
    (cond
      (sequential? type) (first type)
      (map? type) (:type type)
      :else type)))

(defmethod make-typed-input 'Boolean [m owner {:keys [type key val] :as opts}]
  (om/component
    (dom/input #js {:type "checkbox"
                    :checked (or val (empty-value type))
                    :onChange (update-on-change! m key (assoc opts :extract-fn #(.-checked %)))})))

(defmethod make-typed-input 'Number [m owner {:keys [key val] :as opts}]
  (om/component
    (dom/input #js {:type "number"
                    :value val
                    :onChange (update-on-change! m key (assoc opts :transform-fn js/parseFloat))})))

(def keyword-pattern "^:(\\w+|\\w+(\\.\\w+)*\\/\\w+)$")

(defmethod make-typed-input 'Keyword [m owner {:keys [type key val] :as opts}]
  (om/component
    (dom/input #js {:type "text"
                    :value val
                    :pattern keyword-pattern
                    :onChange (update-on-change! m key (assoc opts :transform-fn #(or (read-keyword %) (empty-value type))))})))

(defmethod make-typed-input 'String [m owner {:keys [type key val] :as opts}]
  (om/component
    (dom/input #js {:type "text"
                    :value val
                    :onChange (update-on-change! m key opts)})))

(defmethod make-typed-input 'Value [m owner {type :type}]
  (let [[_ value] type]
    (om/component
      (dom/input (clj->js
                   (into {:value (str value)
                          :readOnly "readOnly"}
                         (cond
                           (or (true? value) (false? value)) {:type "checkbox", :checked value}
                           (number? value) {:type "number"}
                           (keyword? value) {:type "text", :pattern keyword-pattern}
                           :else {:type "text"})))))))

(defmethod make-typed-input 'Option [m owner {type :type :as opts}]
  (let [[_ type] type]
    (make-typed-input m owner (assoc opts :type type))))

(defmethod make-typed-input 'U [m owner {:keys [type key val] :as opts}]
  (om/component
    (dom/select #js {:onChange (update-on-change! m key (assoc opts :transform-fn r/read-string))}
      (into-array
        (map (fn [[_ v]]
               (dom/option nil (str v)))
             (rest type))))))

(defn make-typed-field [m optional?]
  (fn [[k t]]
    (dom/div #js {:className "field"}
      (dom/label nil (str k (when optional? " (optional)")))
      (om/build make-typed-input m {:opts {:type (om/value t), :key k, :val (k m)
                                           :optional? optional?}}))))

(defn read-string-safe [s default]
  (let [r (r/push-back-reader s)]
    (r/read r false default false)))

(defn update-dynamic-field! [type owner]
  (fn [ev]
    (let [k (om/get-state owner :key)
          t (read-string-safe (om/get-state owner :type) nil)]
      (when (and k t)
        (om/transact! type [4] (fn [optional]
                                 (assoc optional k t)))
        (om/set-state! owner :key nil)
        (om/set-state! owner :type nil)))))

(defn make-dynamic-field [m owner type]
  (om/component
    (dom/div #js {:className "field"}
      (dom/input #js {:type "button", :value "+"
                      :onClick (update-dynamic-field! type owner)})
      (dom/input #js {:type "text", :pattern keyword-pattern, :placeholder ":a.key/word"
                      :value (om/get-state owner :key)
                      :onChange #(om/set-state! owner :key (read-keyword (.. % -target -value)))})
      (dom/input #js {:type "text", :value (om/get-state owner :type)
                      :onChange #(om/set-state! owner :type (.. % -target -value))}))))

(defmethod make-typed-input 'HMap [m owner {type :type}]
  (let [hmap (apply hash-map (rest type))
        required (:mandatory hmap)
        optional (:optional hmap)]
    (om/component
      (dom/div nil
        (dom/span nil "{")
        (into-array (map (make-typed-field m false) required))
        (into-array (map (make-typed-field m true) optional))
        (om/build make-dynamic-field m {:opts type})
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