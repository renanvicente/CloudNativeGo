package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

var throttled = Throttle(getHostname, 1, 1, time.Second)

func getHostname(ctx context.Context) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	return os.Hostname()
}

func throttledHandler(w http.ResponseWriter, r *http.Request) {
	ok, hostname, err := throttled(r.Context(), r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if !ok {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(hostname))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hostname",throttledHandler)
	log.Fatal(http.ListenAndServe(":8080", r))
}