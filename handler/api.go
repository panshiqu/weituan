package handler

import (
	"io"
	"log"
	"net/http"
	"os"
)

// ServeHTTP .
func ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

// ServeFiles .
func ServeFiles(w http.ResponseWriter, r *http.Request) {
	name := "." + r.URL.Path

	fi, err := os.Stat(name)
	if err != nil {
		log.Println("ServeFile", err)
		http.NotFound(w, r)
		return
	}

	if !fi.Mode().IsRegular() {
		log.Println("ServeFile not regular", name)
		http.NotFound(w, r)
		return
	}

	f, err := os.Open(name)
	if err != nil {
		log.Println("ServeFile", err)
		http.NotFound(w, r)
		return
	}

	defer f.Close()

	io.Copy(w, f)
}
