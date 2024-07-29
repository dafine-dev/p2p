package endpoints

import (
	"net/http"
	"p2p/actions"
)

var ACTIONS *actions.Actions

func HandleBegin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ACTIONS.Connect()
	} else {
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func HandleFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid", http.StatusBadRequest)
			return
		}

		filename := r.FormValue("filename")
		ACTIONS.InsertFile(filename)

	} else if r.Method == http.MethodGet {
		filename := r.URL.Query().Get("filename")
		ACTIONS.GetFile(filename)

	} else {
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func HandleLeave(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

	} else {
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}
