package main

import (
	"errors"
	"net/http"
	"sync/atomic"

	"github.com/gorilla/pat"
)

type server struct {
	router      *pat.Router
	connections map[string]*connection
}

func NewServer() (*server, error) {
	s := &server{
		router:      pat.New(),
		connections: make(map[string]*connection),
	}
	s.routes()
	return s, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) AddConnection(dongleId string, connection *connection) {
	if _, has := s.connections[dongleId]; has {
		atomic.AddInt32(s.connections[dongleId].Instances, 1)
	}
	s.connections[dongleId] = connection
}

func (s *server) RemoveConnection(dongleId string) {
	if _, has := s.connections[dongleId]; has {
		if *s.connections[dongleId].Instances > 1 {
			atomic.AddInt32(s.connections[dongleId].Instances, -1)
		} else {
			delete(s.connections, dongleId)
		}
	}
}

func (s *server) CleanupConnection(dongleId string, connection *connection) {
	connection.Socket.Close()
	<-connection.Done // stop the pinging
	s.RemoveConnection(dongleId)
}

func (s *server) GetConnectionById(dongleId string) (*connection, error) {
	if _, has := s.connections[dongleId]; has {
		return s.connections[dongleId], nil
	} else {
		return nil, errors.New("no connection found")
	}
}
