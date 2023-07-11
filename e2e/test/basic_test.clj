(ns basic-test
  (:require [clojure.test :refer [deftest is testing]]
            [babashka.http-client :as http]))

(def endpoint "http://127.0.0.1:8080/")

(defn request [headers] (http/get (str endpoint "/cars") {:headers headers :throw false}))

(deftest authorized-cars
  (testing "bob is authorized to list all cars"
    (let [status (:status (request {"Authorization" "bob"}))]
      (is (= 200 status)))))
 
(deftest unauthorized-cars
  (testing "alice is not authorized to list all cars"
    (let [status (:status (request {"Authorization" "alice"}))]
      (is (= 403 status)))))