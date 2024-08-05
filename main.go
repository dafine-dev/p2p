package main

import (
	"net/http"
	"p2p/endpoints"
)

func main() {

	go endpoints.Start()

	fileServer := http.FileServer(http.Dir("./"))

	http.Handle("/", http.StripPrefix("/", fileServer))

	http.ListenAndServe(":7000", nil)
}
