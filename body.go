package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type body struct {
	socket   *websocket.Conn
	response chan []byte
	server   *server
}

func (b *body) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("ServeHttp body")
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	b.socket = socket
	b.response = make(chan []byte, messageBufferSize)
	go b.write()

	b.read()
}

func (b *body) read() {

	for {
		if _, msg, err := b.socket.ReadMessage(); err == nil {
			log.Println("read body")
			b.server.send <- msg
		} else {
			break
		}
	}
	log.Println("Close read body")
	b.socket.Close()
}
func (b *body) write() {
	for response := range b.response {
		if err := b.socket.WriteMessage(websocket.TextMessage, response); err != nil {
			break
		}
	}
	b.socket.Close()
}
