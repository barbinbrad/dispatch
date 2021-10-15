package main

import (
	"encoding/json"
	"net/http"
)

func (server *server) restHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var (
			request  JsonRPCRequest
			response JsonRPCResponse
		)

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dongleId := r.URL.Query().Get(":dongle_id")
		connection, err := server.GetConnectionById(dongleId)
		if err != nil {
			http.Error(w, "Device not available", http.StatusPreconditionFailed)
			return
		}

		ch := make(chan string, 1)

		err = connection.SendAndSubscribe(request, ch)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		data := <-ch
		close(ch)

		err = json.Unmarshal([]byte(data), &response)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
