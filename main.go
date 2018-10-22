package main

import (
	"log"
	"net/http"

	"github.com/panshiqu/weituan/handler"
)

func main() {
	http.HandleFunc("/", handler.ServeHTTP)
	http.HandleFunc("/files/", handler.ServeFiles)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
