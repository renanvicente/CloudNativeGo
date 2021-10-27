package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	"github.com/thedevsaddam/renderer"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Host string
	Port uint16
	Tags map[string]string
}

var config Config
var rnd = renderer.New()

func init() {
	updates, errors, err := watchConfig("CloudNativeGo/Chapter10/config-json/decode/sample.json")
	if err != nil {
		panic(err)
	}

	go startListening(updates, errors)
}

func loadConfiguration(filepath string) (Config, error) {
	dat, err := ioutil.ReadFile(filepath) // Ingest file as []byte
	if err != nil {
		return Config{}, err
	}

	config := Config{}

	err = json.Unmarshal(dat, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func startListening(updates <-chan string, errors <-chan error) {
	for {
		select {
		case filepath := <-updates:
			c, err := loadConfiguration(filepath)
			if err != nil {
				log.Println("error loading config:", err)
				continue
			}
			config = c
		case err := <-errors:
			log.Println("error watching config", err)

		}
	}
}

func watchConfig(filepath string) (<-chan string, <-chan error, error) {
	errs := make(chan error)
	changes := make(chan string)
	hash := ""

	go func() {
		ticker := time.NewTicker(time.Second)

		for range ticker.C {
			newhash, err := calculateFileHash(filepath)
			if err != nil {
				errs <- err
				continue
			}

			if hash != newhash {
				hash = newhash
				changes <- filepath
			}

		}
	}()

	return changes, errs, nil

}

func calculateFileHash(filepath string) (string, error) {
	file, err := os.Open(filepath) // Open the file for reading
	if err != nil {
		return "", err
	}
	defer file.Close() // Be sure to close your file!

	hash := sha256.New() // Use the Hash in crypto/sha256

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	sum := fmt.Sprintf("%x", hash.Sum(nil)) // Get encoded hash sum
	return sum, nil

}

func watchConfigNotify(filepath string) (<-chan string, <-chan error, error) {
	changes := make(chan string)

	watcher, err := fsnotify.NewWatcher() // Get an fsnotify.Watcher
	if err != nil {
		return nil, nil, err
	}

	err = watcher.Add(filepath) // Tell watcher to watch
	if err != nil {             // our config file
		return nil, nil, err
	}

	go func() {
		changes <- filepath // First is ALWAYS a change

		for event := range watcher.Events { // Range over watcher events
			if event.Op&fsnotify.Write == fsnotify.Write {
				changes <- event.Name
			}
		}
	}()
	return changes, watcher.Errors, nil
}

func printConfigHandler(w http.ResponseWriter, r *http.Request) {
	//w.WriteHeader(http.StatusOK)
	log.Println(config)
	rnd.JSON(w, http.StatusOK, config)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", printConfigHandler)
	log.Fatal(http.ListenAndServe(":8080", r))
}
