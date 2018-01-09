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
	var server = &server{
		request:  make(chan *RequestData),
		response: make(chan *ResponseData),
	}
	log.Println("リクエスト！")
	http.Handle("/post", server)
	http.Handle("/", &templateHandler{filename: "index.html"})

	go server.socketIOStart()

	// webサーバー開始
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
