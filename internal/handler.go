package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Config struct {
	Host    string
	Port    int
	Enabled bool
	Headers map[string]string
}
type Window struct {
	id     uint
	server *Server
}

type Server struct {
	id         uint // to allow for registration as a window
	app        *application.App
	config     *Config
	srv        *http.Server
	window     Window
	clients    map[string]client
	clientLock sync.Mutex
}
type message struct {
	Type string
	Data string
}

type callback struct {
	ID     string `json:"id"`
	Result string `json:"result"`
}

type client struct {
	address string
	events  chan message
}

func (c *client) Send(msg message) error {
	var err error
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("connection lost")
		}
	}()
	c.events <- msg
	return err
}

func (c client) close() {
	if _, ok := (<-c.events); ok {
		close(c.events)
	}
}

func (s *Server) removeClient(clientID string) {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	delete(s.clients, clientID)
	s.app.Logger.Info(fmt.Sprintf("client %v disconnected", clientID))
	// s.Info(fmt.Sprintf("client %v disconnected", clientID))
}

func (s *Server) sendToClient(requestID string, message message) {
	client, ok := s.clients[requestID]
	if !ok {
		return
	}
	if err := client.Send(message); err != nil {
		s.removeClient(client.address)
	}
}

func (s *Server) sendToAllClients(msg message) {
	if len(s.clients) == 0 {
		return
	}
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	dead := []client{}
	for _, client := range s.clients {
		if err := client.Send(msg); err != nil {
			dead = append(dead, client)
		}
	}
	for _, d := range dead {
		s.removeClient(d.address)
	}
}

func (s *Server) handleClient(rw http.ResponseWriter, req *http.Request) {
	client := client{
		events:  make(chan message, 5),
		address: req.RemoteAddr,
	}
	// s.Info(fmt.Sprintf("client %v connected", client.Identifier()))
	s.app.Logger.Info(fmt.Sprintf("client %v connected", client.address))
	clientID := req.URL.Query().Get("clientId")
	if clientID != "" {
		// we only save if we have an identifier
		s.clientLock.Lock()
		s.clients[clientID] = client
		s.clientLock.Unlock()
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	for header, value := range s.config.Headers {
		rw.Header().Set(header, value)
	}
	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "Connection does not support streaming", http.StatusBadRequest)
		return
	}

	for {
		timeout := time.After(500 * time.Millisecond)
		select {
		case <-req.Context().Done():
			client.close()
			s.removeClient(client.address)
			return
		case msg := <-client.events:
			fmt.Fprintf(rw, "event: %s\n", msg.Type)
			fmt.Fprintf(rw, "data: %v\n\n", msg.Data)
		case <-timeout:
			continue
		}
		flusher.Flush()
	}

}

func TestLoop(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 10; i++ {
		w.Write([]byte("test loop"))
		time.Sleep(time.Second * 2)
	}
}
