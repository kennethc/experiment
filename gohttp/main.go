package main

import (
	"io"
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Cache-Control", "public")
	io.WriteString(w, "Hello World!\n")
}

func main() {
	http.HandleFunc("/hello", HelloServer)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
