package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Movies App
type MovieResponse struct {
	Feed           []string
	Recommendation []string
}

func main() {
	http.HandleFunc("/movies", fetchMoviesFeedHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fetchMoviesFeedHandler(w http.ResponseWriter, r *http.Request) {
	mr := MovieResponse{
		Feed: []string{"Transformers", "Fault in our stars", "The Old Boy"},
	}
	rms, err := fetchRecommendations()
	if err != nil {
		w.WriteHeader(500)
	}
	mr.Recommendation = rms
	bytes, err := json.Marshal(mr)
	if err != nil {
		w.WriteHeader(500)
	}
	w.Write(bytes)
}

func fetchRecommendations() ([]string, error) {
	resp, err := http.Get("http://localhost:9090/recommendations")
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}
	var mvsr map[string]interface{}
	err = json.Unmarshal(body, &mvsr)
	if err != nil {
		return []string{}, err
	}
	mvsb, err := json.Marshal(mvsr["movies"])
	if err != nil {
		return []string{}, err
	}
	var mvs []string
	err = json.Unmarshal(mvsb, &mvs)
	if err != nil {
		return []string{}, err
	}
	return mvs, nil
}

