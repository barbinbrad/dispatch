package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func (s *server) websocketHandler() http.HandlerFunc {

	var (
		upgrader   = websocket.Upgrader{}
		pingPeriod = 5
		waitPeriod = 10
	)

	return func(w http.ResponseWriter, r *http.Request) {
		dongleId := r.URL.Query().Get(":dongle_id")
		socket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Websocket handshake", http.StatusPreconditionFailed)
		}

		connection, _ := NewConnection(socket)

		s.AddConnection(dongleId, &connection)

		defer s.CleanupConnection(dongleId, &connection)
		go connection.PingAtInterval(pingPeriod, waitPeriod)

		for {
			err := connection.ReceiveAndPublish(socket.ReadMessage())

			if err != nil {
				fmt.Println(err.Error())
				break
			}
		}
	}
}
