(defproject clarity "0.0-SNAPSHOT"
  :dependencies [[org.clojure/clojure "1.5.1"]
                 [org.clojure/core.typed "0.2.19"]

                 [com.datomic/datomic-free "0.9.4331"]

                 [ring "1.2.1"]
                 [compojure "1.1.6"]

                 [hiccup "1.0.4"]]
  :source-paths ["."]
  :plugins [[lein-ring "0.8.8"]]
  :ring {:handler clarity.server/app})
