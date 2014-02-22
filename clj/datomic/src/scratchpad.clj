(ns scratchpad
  "A place for experiments"
  (:require [datomic.api :as d]))

(defn connect-to [db-uri]
  (d/create-database db-uri)
  (d/connect db-uri))

(def conn (connect-to "datomic:free://localhost:4334/scratchpad"))

(defn db [] (d/db conn))

(d/transact conn
  [{:db/id (d/tempid :db.part/db)
    :db/ident :person/name
    :db/valueType :db.type/string
    :db/cardinality :db.cardinality/one
    :db.install/_attribute :db.part/db}])

(def nums [1 2 3 4 5 6])

(def fancy-nums (atom [42]))

(def immutable-nums (deref fancy-nums))
immutable-nums

(swap! fancy-nums conj 3.145)
@fancy-nums

(d/transact conn
  [{:db/id (d/tempid :db.part/user)
    :person/name "Tom"}])

(d/q '[:find ?e
       :where [?e :person/name "Tom"]]
     (db))

(d/transact conn
  [{:db/id (d/tempid :db.part/db)
    :db/ident :place/name
    :db/valueType :db.type/string
    :db/cardinality :db.cardinality/one
    :db.install/_attribute :db.part/db}])

(def db1 (db))

db1

(d/touch (d/entity (db) 17592186045418))

(d/transact conn
  [{:db/id 17592186045418
    :place/name "Leipzig"}])

(d/q '[:find ?e1 ?e2 ?name ?place
       :where [?e1 :person/name ?name]
              [?e2 :place/name ?place]]
     (db))

(d/transact conn
  [{:db/id (d/tempid :db.part/user)
    :person/name "Paul"
    :place/name "Paris"}])

(d/transact conn
  [{:db/id 17592186045418
    :person/name "Klaus"}])

(d/q '[:find ?old-name ?new-name
       :where [?e :person/name ?old-name ?t1]
              [?e :person/name ?new-name ?t2]
              [(not= ?old-name ?new-name)]
              [(< ?t1 ?t2)]]
  (d/history (db)))

(d/q '[:find ?e ?a ?v
       :where [?e ?a ?v]]
  (db))
