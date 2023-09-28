package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"hash/crc32"
	"html/template"
	"log"
	"net/http"
	"path"
)

var urls = make(map[string]string)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/{shortUrl}", redirect).Methods("GET")
	r.HandleFunc("/shorten", shorten).Methods("POST")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})
	log.Fatalln(http.ListenAndServe(":8080", r))
}

func redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortUrl := vars["shortUrl"]
	longUrl := urls[shortUrl]
	http.Redirect(w, r, longUrl, 301)
}

func shorten(w http.ResponseWriter, r *http.Request) {
	longUrl := r.FormValue("url")
	shortUrl := r.FormValue("slug")
	if shortUrl == "" {
		shortUrl = fmt.Sprintf("%x", crc32.ChecksumIEEE([]byte(longUrl)))
	}
	urls[shortUrl] = longUrl

	fp := path.Join("static", "shortened.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Url    string
		Prefix string
	}{
		Url:    shortUrl,
		Prefix: r.Host,
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
