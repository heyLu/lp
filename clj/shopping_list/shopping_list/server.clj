(ns shopping-list.server
  (:require [shopping-list :as sl])
  (:use ring.util.response)
  (:use compojure.core)
  (:require [compojure.route :as route])
  (:use hiccup.core))


(defn render-groups [groups]
  (for [[group-name group-items] groups]
    [:div
     [:h2 (str group-name)]
     [:ul
      (for [item group-items]
        (let [{:keys [name price weight category]} item]
          [:li (str name)]))]]))

(defn index []
  (html
   [:html
    [:head [:title "BUY ALL THE THINGS"]]
    [:body
     [:div {:id "main"}
      [:h1 "Hi there"]
      [:p "Mhh yeah, what's up?"]
      (render-groups (sl/sort-map (comp identity compare) (sl/regroup-by-name sl/example-articles)))]]]))

(defroutes shopping-list
  (GET "/" [] (index))
  (route/not-found "<h1>oops. alternate reality monsters</h1>"))
