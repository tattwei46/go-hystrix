package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", logger(handle))
	log.Println(" Producer listening on :8081")
	http.ListenAndServe(":8081", nil)

}

func handle(w http.ResponseWriter, r *http.Request) {
	var msg = "success"
	resp, _ := json.Marshal(msg)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

// logger is Handler wrapper function for logging
func logger(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path, r.Method)
		fn(w, r)
	}
}
