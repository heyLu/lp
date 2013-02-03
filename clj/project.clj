(defproject lp-clj "0.0"
  :description "Playing around with Clojure"
  :dependencies [[org.clojure/clojure "1.4.0"]
                 [compojure "1.1.5"]
                 [ring-middleware-format "0.2.4"]
                 [cheshire "5.0.1"]
                 [hiccup "1.0.2"]]
  :plugins [[lein-ring "0.8.2"]]
  :ring {:handler shame/serve}
  :source-paths ["."])
