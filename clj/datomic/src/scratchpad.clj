(ns scratchpad
  "A place for experiments"
  (:require [datomic.api :as d]))

(defn connect-to [db-uri]
  (d/create-database db-uri)
  (d/connect db-uri))

(def conn (connect-to "datomic:free://localhost:4334/scratchpad"))

(defn db [] (d/db conn))
