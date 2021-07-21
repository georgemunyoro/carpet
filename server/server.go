package server

import (
	"log"
	"net/http"
)

func Serve(filepath string) {
	http.Handle("/", http.FileServer(http.Dir(filepath)))
	log.Fatal(http.ListenAndServe(":8090", nil))
}
