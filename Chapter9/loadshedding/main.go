package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const MaxQueueDepth = 2

var count500 = 1

func CurrentQueueDepth(c *int) int {
	*c = *c +1
	return *c
}
// Middleware function, which will be called for each request.
// If queue depth is exceeded, it returns HTTP 503 (service unavailable).
func loadSheddingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CurrentQueueDepth is fictional and for example purposes only.
		current := CurrentQueueDepth(&count500)
		fmt.Println(current)
		if current > MaxQueueDepth {
			log.Println("load shedding engaged")
			http.Error(w, "load shedding engaged", http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getHostname(w http.ResponseWriter, r *http.Request) {
	//hostname, err := os.Hostname()
	//fmt.Println(hostname)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Test"))
}

func main() {
	r := mux.NewRouter()
	// Register middleware
	r.Use(loadSheddingMiddleware)
	r.HandleFunc("/hostname",getHostname)
	log.Fatal(http.ListenAndServe(":8080",r))
}