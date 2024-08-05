package endpoints

import (
	"fmt"
	"log"
	"net/http"
	"p2p/actions"
	"p2p/shared"
	"strings"

	"github.com/gorilla/websocket"
)

var ACTIONS *actions.Actions
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	defer conn.Close()

	for {
		msgType, msg, err := conn.ReadMessage()

		data := strings.Split(string(msg), ";")
		if data[0] == "INSERT" {
			ACTIONS.InsertFile(string(data[1]))
		} else if data[0] == "SEARCH" {
			ACTIONS.GetFile(data[1])
		}

		if err != nil {
			log.Println(err)
			return
		}

		err = conn.WriteMessage(msgType, []byte("Files;"+ACTIONS.Files()))
		if err != nil {
			log.Println(err)
			break
		}

		err = conn.WriteMessage(msgType, []byte("Locations;"+ACTIONS.ListLocations()))
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func Start() {
	ACTIONS = actions.New(shared.ParseAddr(0, 0, 0, 0), "./server")
	go ACTIONS.Run(true)

	fmt.Println(ACTIONS.Files())

	ACTIONS.Connect()
	http.HandleFunc("/ws", handler)
	serverAddr := "localhost:8080"
	fmt.Printf("WebSocket server started at ws://%s\n", serverAddr)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
