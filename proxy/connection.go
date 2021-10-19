package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type connection struct {
	Socket           *websocket.Conn
	NextId           *uint32
	Instances        *int32
	ResponseRegistry map[uint32]chan string
	Done             chan bool
}

func NewConnection(socket *websocket.Conn) (connection, error) {
	c := connection{
		Socket:           socket,
		NextId:           new(uint32),
		Instances:        new(int32),
		ResponseRegistry: make(map[uint32]chan string),
		Done:             make(chan bool),
	}
	atomic.AddInt32(c.Instances, 1)
	return c, nil
}

func (c *connection) PingAtInterval(seconds int, wait int) {
	pingPeriod := time.Duration(seconds) * time.Second
	waitPeriod := time.Duration(wait) * time.Second

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := c.Socket.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(waitPeriod)); err != nil {
				fmt.Println("ping", err)
			}
		case <-c.Done:
			return
		}
	}

}

func (c *connection) NewMessageId() uint32 {
	return atomic.AddUint32(c.NextId, 1)
}

func (connection *connection) SendAndSubscribe(request JsonRPCRequest,
	ch chan string) error {

	messageId := connection.NewMessageId()

	request.Id = messageId
	data, err := json.Marshal(request)
	if err != nil {
		return err
	}

	err = connection.RegisterChannel(messageId, ch)
	if err != nil {
		return err
	}

	connection.Socket.WriteMessage(websocket.TextMessage, data)

	return nil
}

func (c *connection) ReceiveAndPublish(messageType int, data []byte, err error) error {
	var response JsonRPCResponse

	if err != nil {
		return err
	}

	if messageType != 1 {
		return errors.New("invalid message type")
	}

	err = json.Unmarshal(data, &response)
	if err != nil {
		return err
	}

	messageId := response.Id

	registeredChannel, err := c.GetChannel(messageId)
	if err != nil {
		return err
	}

	registeredChannel <- string(data)

	err = c.UnregisterChannel(messageId)
	if err != nil {
		return err
	}

	return nil
}

func (c *connection) RegisterChannel(messageId uint32, ch chan string) error {
	c.ResponseRegistry[messageId] = ch
	return nil
}

func (c *connection) GetChannel(messageId uint32) (chan string, error) {
	if c.ResponseRegistry[messageId] == nil {
		return nil, errors.New("channel not found")
	}
	channel := c.ResponseRegistry[messageId]
	return channel, nil
}

func (c *connection) UnregisterChannel(messageId uint32) error {
	delete(c.ResponseRegistry, messageId)
	return nil
}
