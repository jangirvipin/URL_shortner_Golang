package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"errors"

	"time"
)

type DATA struct {
	ID          string    `json:"id"`
	OriginalURL string    `json:"originalurl"`
	ShortURL    string    `json:"shorturl"`
	CreatedAt   time.Time `json:"createdat"`
}

var Database = make(map[string]DATA)

func hashURL(originalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(originalURL))
	data := hasher.Sum(nil)
	hash := hex.EncodeToString(data)
	return hash[:8]
}

func Creation(originalURL string) string {
	shorturl := hashURL(originalURL)
	id := shorturl
	createdAt := time.Now()
	Database[id] = DATA{
		ID:          id,
		OriginalURL: originalURL,
		ShortURL:    shorturl,
		CreatedAt:   createdAt,
	}
	return shorturl
}

func GetURL(id string) (string, error) {
	url, ok := Database[id]

	if !ok {
		return "", errors.New("URL not found")
	}

	shortenURL:= url.OriginalURL
	return shortenURL,nil
}

func ShortUrlhandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Url string `json:"url"`
	}
	var short struct {
		ShortURL string `json:"shorturl"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shortURL := Creation(data.Url)
	short.ShortURL = shortURL

	json.NewEncoder(w).Encode(short)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Path[len("/redirect/"):]

	redirectedURL,err:= GetURL(url)

	if err != nil {
		http.Error(w, "Invalid", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, redirectedURL, http.StatusFound)
}

func main() {
	port := ":8080"
	
	http.HandleFunc("/shorten", ShortUrlhandler)
	http.HandleFunc("/redirect/", RedirectHandler)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("`Error starting server: ", err)
	}
}
