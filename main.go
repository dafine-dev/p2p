package main

import (
	"net/http"
	"p2p/endpoints"
	"p2p/shared"
)

func main() {

	shared.CalculateAddr()
	go endpoints.Start()

	fileServer := http.FileServer(http.Dir("./"))

	http.Handle("/", http.StripPrefix("/", fileServer))

	http.ListenAndServe(":7000", nil)
}
