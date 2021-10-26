package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func helloMuxHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello gorilla/mux!\n"))
}

func ProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Write([]byte(vars["key"]))
}

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Found me!\n"))
}

func ArticlesCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(""))
}

func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(""))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/products",ProductsHandler).
		Host("www.example.com").  // Only match a specific domain
		Methods("GET","PUT"). // Only match GET+PUT methods
		Schemes("http")        // Only match the http scheme
	r.HandleFunc("/products/{key}", ProductHandler)
	r.HandleFunc("/articles/{category}",ArticlesCategoryHandler)
	r.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler)
	r.HandleFunc("/",helloMuxHandler)
	log.Fatal(http.ListenAndServe(":8080",r))
}
