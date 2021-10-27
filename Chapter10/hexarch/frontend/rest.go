package frontend

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/renanvicente/CloudNativeGo/Chapter10/hexarch/core"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// restFrontEnd contains a reference to the core application logic,
// and complies with the contract defined by the FrontEnd interface.
type restFrontEnd struct {
	store *core.KeyValueStore
}

// Set to true if you're working on the new storage backend
var useConsul = os.Getenv("USE_CONSUL")

// keyValuePutHandler expects to be called with a PUT request for
// the "/v1/key/{key}" resource.
func (f *restFrontEnd) keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = f.store.Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	log.Printf("PUT key=%s value=%s\n", key, string(value))
}

func (f *restFrontEnd) keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Retrieve "key" from the request
	key := vars["key"]
	var value string
	var err error
	if useConsul == "true" {
		value, err = f.store.ConsulGet(key)
	} else {
		value, err = f.store.Get(key) // Get value for key
	}
	if errors.Is(err, core.ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(value)) // Write the value to the response
	log.Printf("GET key=%s\n", key)
}
func (f *restFrontEnd) keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	err := f.store.Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Printf("DELETE key=%s\n", key)
}

func (f *restFrontEnd) sampleTestPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = f.store.Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	log.Printf("PUT key=%s value=%s\n", key, string(value))
}

// Start includes the setup and start logic that previously
// lived in a main function.
func (f *restFrontEnd) Start(store *core.KeyValueStore) error {
	// Remember our core application reference.
	f.store = store
	r := mux.NewRouter()
	// Register keyValuePutHandler as the handler function for PUT
	// requests matching "/v1/{key}"
	r.HandleFunc("/v1/{key}", f.keyValuePutHandler).Methods("PUT")
	// Register keyValueGetHandler as the handler function for GET
	// requests matching "/v1/{key}"
	r.HandleFunc("/v1/{key}", f.keyValueGetHandler).Methods("GET")
	// Register keyValueDeleteHandler as the handler function for DELETE
	// requests matching "/v1/{key}"
	r.HandleFunc("/v1/{key}", f.keyValueDeleteHandler).Methods("DELETE")
	//log.Fatal(http.ListenAndServe(":8080",r))

	// Sample Test Post
	r.HandleFunc("/v1/c8/{key}", f.sampleTestPost).Methods("POST")
	return http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", r)

}
