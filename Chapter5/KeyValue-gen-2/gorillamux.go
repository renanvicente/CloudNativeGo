package main

import (
	"errors"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)


// keyValuePutHandler expects to be called with a PUT request for
// the "/v1/key/{key}" resource.
func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	log.Printf("PUT key=%s value=%s\n", key, string(value))
}

func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)			// Retrieve "key" from the request
	key := vars["key"]
	value, err := Get(key)		// Get value for key
	if errors.Is(err, ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(value))		// Write the value to the response
	log.Printf("GET key=%s\n", key)
}
func keyValueDeleteHandler(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	key := vars["key"]
	err := Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Printf("DELETE key=%s\n", key)
}
func main() {
	r := mux.NewRouter()
	// Register keyValuePutHandler as the handler function for PUT
	// requests matching "/v1/{key}"
	r.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	// Register keyValueGetHandler as the handler function for GET
	// requests matching "/v1/{key}"
	r.HandleFunc("/v1/{key}", keyValueGetHandler).Methods("GET")
	// Register keyValueDeleteHandler as the handler function for DELETE
	// requests matching "/v1/{key}"
	r.HandleFunc("/v1/{key}", keyValueDeleteHandler).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080",r))
}
