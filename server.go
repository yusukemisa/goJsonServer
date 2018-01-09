package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type server struct {
	// sumally用ソケット
	socket  *websocket.Conn
	host    chan string
	sumally chan []byte
	// リクエスト待ち受けチャネル
	send chan []byte
	body *body
}

func (s *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("ServeHttp server")
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	s.socket = socket

	go s.write()

	s.read()
}

func (s *server) read() {
	for {
		if _, msg, err := s.socket.ReadMessage(); err == nil {
			log.Println("read server")
			s.host <- string(msg)
		} else {
			break
		}
	}
	s.socket.Close()
}
func (s *server) write() {
	for sumally := range s.sumally {
		if err := s.socket.WriteMessage(websocket.TextMessage, sumally); err != nil {
			break
		}
	}
	s.socket.Close()
}

var hostName string

func (s *server) run() {
	for {
		select {
		case host := <-s.host:
			// 送信先設定
			log.Println("host run server")
			hostName = host
			log.Println(hostName)
		case request := <-s.send:
			// POSTリクエスト
			log.Println("request run server")
			resHeader, resBody := post(request)
			log.Printf("After post server:%v", string(resBody))
			s.sumally <- resHeader
			s.body.response <- resBody
		}
	}
}

type Test struct {
	Test string
}

func post(request []byte) (resHeader []byte, resBody []byte) {
	//httpからCliantインスタンス取得
	//Cliant構造体にはhttp cliantに必要なメソッドが備わっている
	log.Println("post server")
	cliant := &http.Client{}

	postBody := bytes.NewBuffer(request)
	log.Println(postBody)
	res, err := cliant.Post(hostName, "application/json;charset=UTF-8", postBody)
	log.Println("ここまできた")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	log.Printf("ここまできた:%v", res.StatusCode)
	// if res.StatusCode != http.StatusOK {
	// 	log.Fatal(err)
	// }
	log.Printf("Content-Type:%v", res.Header.Get("Content-Type"))

	// Response body read
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	// json -> struct
	var test Test
	if err = json.Unmarshal(body, &test); err != nil {
		log.Fatal(err)
	}
	sumally := "Request OK:Content-Type:" + res.Header.Get("Content-Type")
	return []byte(sumally), body
}
