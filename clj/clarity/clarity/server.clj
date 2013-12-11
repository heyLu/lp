(ns clarity.server
  (:use compojure.core)
  (:use ring.middleware.params
        ring.middleware.keyword-params)
  (:use [clojure.java.io :only (reader)])
  (:require [clojure.tools.reader.edn :as edn])
  (:use [hiccup.core :only (html)])

  (:require clarity)
  (:require [datomic.api :as d]))

(defonce conn
  (do
    (d/create-database "datomic:mem://self.data")
    (d/connect "datomic:mem://self.data")))

(defn read-tx-data [str]
  (edn/read-string {:readers {'db/id (partial apply d/tempid)}} str))

(defn http-error [status body & {:as attrs}]
  (into {:status status
         :body body}
        attrs))

(defroutes app-routes
  (GET "/" []
       {:status 200
        :headers {"Content-Type" "text/html"}
        :body (html [:form#query {:action "/facts" :accept-charset "utf-8"}
                     [:textarea {:name "q" :cols 80 :rows 20}]
                     [:input {:type "submit" :value "query"}]]
                    [:form#facts {:action "/facts" :method "post"}
                     [:textarea {:name "facts" :cols 80 :rows 20}]
                     [:input {:type "submit" :value "transact!"}]])})
  (GET "/facts" {{query :q} :params}
       (if query
         (pr-str (d/q (edn/read-string query) (d/db conn)))
         (http-error 400 (pr-str {:error "Missing required `q` parameter"}))))
  (POST "/facts" [facts]
        (pr-str (d/transact conn (read-tx-data facts)))))

(def app
  (-> app-routes
      wrap-keyword-params
      wrap-params))
