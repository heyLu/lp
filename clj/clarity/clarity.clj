(ns clarity
  (:use clojure.core.typed))

; (HMap :mandatory {:a ... :b ...}) -> { :a input :b input }
;   needs nesting, dynamic rules?

; - required/optional? sometimes comes from above, an environment?
; - optional fields: adding dynamically would be better?
;     * but how to realize without creating an abstract ui model?
; - additional attributes (req/opt, hidden, defaults) are a problem
;   in general as they are quite tightly coupled to the implementations,
;   is there a name for this? is this a leaking abstraction?

; what do i want? from a description of the data i want to create an
; input method that only allows inserting valid values.
; sometimes we have existing data present or required so we must
; render that as part of the input tree. we also sometimes want to add
; information dynamically, for example adding key-value-pairs to maps
; and new elements to collections. to be able to do that we need to
; dynamically add new typed input methods on request. we also want to
; delete information sometimes, but only when it is not required.

(defmulti make-input (fn [type & _]
                       (if (seq? type)
                         (first type)
                         type)))

(defmethod make-input 'Value [type]
  (let [x (second type)
        attrs {:value (str x)
               :readonly ""}]
    [:input (into attrs
              (cond
               (instance? Boolean x) {:type "checkbox"}
               (number? x) {:type "number"}
               (or (instance? java.net.URL x) (instance? java.net.URI x)) {:type "url"}
               (keyword? x) {:pattern "^:(\\w+|\\w+(\\.\\w+)*\\/\\w+)$"
                             :type "text"}
               :else {:type "text"}))]))

(defmethod make-input 'Option [[_ type]]
  (make-input type))

(defmethod make-input 'U [[_ & alts]]
  [:select
   (map (fn [[_ x]]
          [:option (str x)])
        alts)])

(defmethod make-input 'Keyword [_]
  [:input {:type "text"
           :placeholder ":my.ns/identifier"
           :pattern "^:(\\w+|\\w+(\\.\\w+)*\\/\\w+)$"}])

(defmethod make-input 'String [_]
  [:input {:type "text"}])

(defmethod make-input 'Boolean [_]
  [:input {:type "checkbox"}])

(defmethod make-input 'Number [_]
  [:input {:type "number"}])

(defmethod make-input 'HVec [[_ & types]]
  ; display existing values as editable *and* allow adding new elements
  ; those elements can be of multiple types -> dynamism required?
  )

(defmethod make-input 'HMap [[ _ & {:keys [mandatory optional]}]]
  (concat '("{")
          (map (fn [[ key val]]
                 [:div.field
                  [:label (str key)]
                  (make-input val)])
               mandatory)
          (map (fn [[key val]]
                 [:div.field.optional
                  [:label (str key)]
                  (make-input val)])
               optional)
          '("}")))
