package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func (server *server) websocketHandler() http.HandlerFunc {

	var (
		upgrader     = websocket.Upgrader{}
		pingInterval = 5
		waitInterval = 10
	)

	return func(w http.ResponseWriter, r *http.Request) {
		dongleId := r.URL.Query().Get(":dongle_id")
		socket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
		}

		done := make(chan struct{})
		connection, _ := NewConnection(socket, done)

		server.AddConnection(dongleId, &connection)
		defer server.CleanupConnection(dongleId, &connection)

		go connection.PingAtInterval(pingInterval, waitInterval, done)

		for {
			err := connection.ReceiveAndPublish(socket.ReadMessage())

			if err != nil {
				fmt.Println(err.Error())
				break
			}
		}
	}
}
