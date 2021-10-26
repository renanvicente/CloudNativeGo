package main

import (
	"io/ioutil"
	"net/http"
	"os"
)

func healthShallowHandler(w http.ResponseWriter, r *http.Request) {
	// Create our test file.
	// This will create a filename like /tmp/shallow-123456
	tmpFile, err := ioutil.TempFile(os.TempDir(), "shallow-")
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer os.Remove(tmpFile.Name())

	// Make sure that we can write to the file.
	text := []byte("Check.")
	if _, err = tmpFile.Write(text); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	// Make sure that we can close the file.
	if err := tmpFile.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}
