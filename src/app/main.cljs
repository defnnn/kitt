(ns app.main)

(defonce token 0)

(defn reload! []
  (println (str "value of token " token)))

(defn main! []
  (println "started!"))

(comment
  (prn token))
