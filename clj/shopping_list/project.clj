(defproject shopping-list "0.0-SNAPSHOT"
  :description "A shopping list in Clojure & ClojureScript"
  :dependencies [[org.clojure/clojure "1.5.0"]
                 [ring "1.2.0-beta2"]
                 [compojure "1.1.5"]]
  :plugins [[lein-ring "0.8.3"]
            [lein-cljsbuild "0.3.0"]]
  :ring {:handler shopping-list.server/shopping-list}
  :cljsbuild {:builds [{:source-paths ["."]
                        :compiler {:output-to "js/shoppinglist.js"
                                   :optimizations :whitespace
                                   :pretty-print true}}]}
  :source-paths ["."])
