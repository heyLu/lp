(ns shame
  (:use [clojure.data.json :as json :only []])
  (:use compojure.core)
  (:require [compojure.handler :as handler]
            [compojure.route :as route])
  (:use ring.util.response)
  (:use ring.middleware.reload)
  (:use ring.middleware.stacktrace)
  (:use [ring.middleware.format-response :only [wrap-restful-response]])
  (:use [ring.middleware.format-params :only [wrap-restful-params]])
  (:use hiccup.core))


; roadmap:
;  1. implement add-item, resurrect-item, close-item
;  2. use refs
;  3. persist to disk
;  4. serve as local web-page
;      caveats: client-side needed (add/close/resurrect, index, note)
;  5. publish

;; todo manipulation

(defn get-by [map-coll k v]
  (first (drop-while #(not= (k %) v) map-coll)))

(defn eq-by-key [key m1 m2]
  (= (key m1) (key m2)))

(defn same-name [m1 m2]
  (eq-by-key :name m1 m2))

(defn transplant [item association from to & [eql-rel]]
  "Transplants an item in a map of vectors from one key to another."
  (let [eql-rel (or eql-rel =)])
  (assoc association
         from (vec (filter #(not ((or eql-rel =) item %)) (from association)))
         to   (conj (to association) item)))

(defn add-item [item shaming]
  "Try adding an item to the shaming. If the maximum number of items is
   exceeeded, return the original shaming."
  (let [c (count (:current shaming))]
    (if (< c (get-in shaming [:config :number-of-items]))
      (assoc-in shaming [:current c] item)
      shaming))) ; FIXME: is that the clojure way of doing it?

(defn change-item [item-name changes shaming]
  (assoc shaming
         :current (mapv #(if (= item-name (:name %))
                           (into % changes)
                           %)
                        (:current shaming))))

(defn close-item [item-name status shaming]
  (let [item (get-by (:current shaming) :name item-name)
        item (assoc item :closed-at (java.util.Date.))]
    (transplant item shaming :current :past same-name)))

(defn resurrect-item [item-name shaming]
  (let [item (get-by (:past shaming) :name item-name)
        item (assoc item :started-at (java.util.Date.))]
    (transplant item shaming :past :current same-name)))

;; references (for "mutating" the todo-list)

(def ^:dynamic *shaming*
  (ref nil))

;; "backend"

(defn read-shame [filename]
  (json/read-str (slurp filename) :key-fn keyword))

(defn write-shame [filename shaming]
  (spit filename (json/write-str shaming)))

;; web service

(dosync
  (ref-set *shaming* (read-shame "self.todo.json")))

(defn take-n [start amount coll]
  (take amount (drop start coll)))

(defn int-param [param & [default]]
  "Parse an integer from a string, with an optional default value (which defaults to 0)."
  (try (Integer/valueOf (or param default 0))
    (catch NumberFormatException n (or default 0))))

(defn filter-keys [map keys]
  (reduce-kv #(if (some (fn [k] (= k %2)) keys)
                (assoc %1 %2 %3)
                %1)
             {} map))

(defroutes shame-api
  (GET "/" [] (response @*shaming*))
  (POST "/" [item] (dosync (ref-set *shaming*
                                      (add-item item @*shaming*))))
  (GET "/current" [start limit]
       (response
         (take-n (int-param start) (int-param limit 10) (:current @*shaming*))))
  (GET "/past" [start limit]
       (response
         (take-n (int-param start) (int-param limit 10) (:past @*shaming*))))
  (GET "/:item-name" [item-name] (response (get-by (:current @*shaming*)
                                         :name item-name)))
  (PUT "/:item-name" [item-name :as {params :params}] ; changes should be in the body?
       (dosync (ref-set *shaming*
                        (change-item item-name (filter-keys params [:name])  @*shaming*))))
  (DELETE "/:item-name" [item-name status]
          (dosync (ref-set *shaming*
                           (close-item item-name (or status "failed") @*shaming*))))
  (route/not-found "404 - Alternate Reality Monsters"))

(defn render-item [item]
  [:div.shame.current
   [:h2 (:name item)]
   (when (count (:notes item))
     [:ul (for [note (:notes item)]
            [:li (if (string? note)
                   note
                   (:content note))])])])

(defn shame-index []
  (html
    (let [title (str "shame on me, I'm " (count (:current @*shaming*)) " things behind!")]
      [:html
       [:head
        [:title title]
        [:style {:rel "stylesheet", :href "todo.css"}]]
       [:body
        [:h1 title]
        (for [item (:current @*shaming*)]
          (render-item item))]])))

(defroutes shame-routes
  (context "/todos" [] shame-api)
  (GET "/" []
       (content-type
         (response (shame-index))
         "text/html"))
  (route/not-found "404 - Alternate Reality Monsters (in App)"))

(def serve
  (-> (handler/site shame-routes)
    (wrap-reload {:dirs ["."]})
    (wrap-stacktrace)
    (wrap-restful-params)
    (wrap-restful-response)))
