package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/panshiqu/weituan/define"
)

func serveLogin(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// ServeHTTP .
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	switch r.URL.Path {
	case "/login":
		err = serveLogin(w, r)

	default:
		err = define.ErrUnsupportedAPI
	}

	if _, ok := err.(*define.MyError); ok {
		fmt.Fprint(w, err)
	}

	log.Println(r.URL.Path, err)
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
