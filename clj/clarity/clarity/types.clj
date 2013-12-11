(ns clarity.types
  "Soft type checks using the core.typed type-syntax.")

(defmulti friendly-check (fn [form type]
                           (if (seq? type)
                             (first type)
                             type)))

(defmethod friendly-check 'Value [val [_ expected]]
  (if (= val expected)
    true
    {:error (str "Expected value " expected ", but found " val ".")}))

(defn def-friendly-check [sym pred]
  (defmethod friendly-check sym [val _]
    (if (pred val)
      true
      {:error (str "Expected value of type " (name sym) ", but got value " val ".")})))

(def-friendly-check 'String string?)
(def-friendly-check 'Keyword keyword?)
(def-friendly-check 'Boolean #(or (true? %1) (false? %1)))

(defmethod friendly-check 'U [val [_ & types]]
  (if (some true? (map (partial friendly-check val) types))
    true
    {:error (str "Expected one of " types ", but found " val ".")}))

(defn friendly-check-keys [m mandatory _]
  (map (fn [[key type]]
         (if-let [val (get m key)]
           (let [check (friendly-check val type)]
             (if (true? check)
               true
               {key (:error check)}))
           {key (str "No value found, but expected one of type " type ".")}))
       mandatory))

(defmethod friendly-check 'HMap [val [_ & {:keys [mandatory optional]}]]
  (if (map? val)
    (let [key-checks (friendly-check-keys val mandatory optional)]
      (if (every? true? key-checks)
        true
        {:error (filter (comp not true?) key-checks)}))
    {:error (str "Expected value of type Map, but got value " val ".")}))

(def datomic-attr-type
  '(HMap :mandatory {:db/id (Value "#db/id[db.part/db]")
                     :db/ident Keyword
                     :db/valueType (U (Value :db.type/keyword)
                                      (Value :db.type/string)
                                      (Value :db.type/boolean)
                                      (Value :db.type/long)
                                      (Value :db.type/bigint)
                                      (Value :db.type/float)
                                      (Value :db.type/double)
                                      (Value :db.type/bigdec)
                                      (Value :db.type/ref)
                                      (Value :db.type/instant)
                                      (Value :db.type/uuid)
                                      (Value :db.type/uri)
                                      (Value :db.type/bytes))
                     :db/cardinality (U (Value :db.cardinality/one)
                                        (Value :db.cardinality/many))
                     :db.install/_attribute (Value :db.part/db)}
         :optional {:db/doc String
                    :db/unique (Option (U (Value :db.unique/value)
                                          (Value :db.unique/identity)))
                    :db/index Boolean ; what about defaults?
                    :db/fulltext Boolean
                    :db/isComponent (Value :db.type/ref) ; FIXME: needs a custom type (dynamic even, because only valid refs should be allowed)
                    :db/noHistory Boolean}))
