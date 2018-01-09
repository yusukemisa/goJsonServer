package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 2048
	messageBufferSize = 512
)

type templateHandler struct {
	once     sync.Once
	filename string
	temple   *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// t.once.Do(func() {
	// 	t.temple = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	// })
	t.temple = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	log.Println("ServeHttp templ")
	t.temple.Execute(w, r)
}

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func main() {
	body := &body{}
	server := &server{
		body:    body,
		send:    make(chan []byte, messageBufferSize),
		sumally: make(chan []byte, messageBufferSize),
		host:    make(chan string, 64),
	}
	body.server = server

	log.Println("リクエスト！")
	http.Handle("/sumally", server)
	http.Handle("/body", body)
	http.Handle("/", &templateHandler{filename: "index.html"})

	go server.run()

	// webサーバー開始
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
