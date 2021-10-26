package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var logger TransactionLogger

func initializeTransactionLog() error {
	var err error

	//logger, err = NewFileTransactionLogger("transaction.log")
	logger, err = NewPostgresTransactionLogger(PostgresDBParams{
		dbName: os.Getenv("DB_NAME"), // kvs
		host: os.Getenv("DB_HOST"), //db
		user: os.Getenv("DB_USER"), // test
		password: os.Getenv("DB_PASSWORD"),// kvstest
		//transactionTable: "transactions",
	})
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errors := logger.ReadEvents()
	e, ok := Event{}, true
	for ok && err == nil {
		select {
		case err, ok = <-errors:		// Retrieve any errors
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete:			// Got a DELETE event!
				err = Delete(e.Key)
			case EventPut:				// Got a PUT event!
				err = Put(e.Key,e.Value)
			}
		}
	}
	logger.Run()
	return err
}

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
	logger.WritePut(key, string(value))
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
	logger.WriteDelete(key)
	w.WriteHeader(http.StatusOK)
	log.Printf("DELETE key=%s\n", key)
}

func sampleTestPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = Put(key,string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.WritePut(key, string(value))
	w.WriteHeader(http.StatusCreated)
	log.Printf("PUT key=%s value=%s\n", key, string(value))
}
func main() {
	// Initializes the transaction log and loads existing data, if any.
	// Blocks until all data is read.
	err := initializeTransactionLog()
	if err != nil {
		panic(err)
	}
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
	//log.Fatal(http.ListenAndServe(":8080",r))

	// Sample Test Post
	r.HandleFunc("/v1/c8/{key}", sampleTestPost).Methods("POST")
	log.Fatal(http.ListenAndServeTLS(":8080","cert.pem","key.pem",r))

}
