package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// RequestData 画面から送信されてくるデータ構造
type RequestData struct {
	// publicにしないとJSON変換処理から見えなくなる
	// ので先頭は大文字にする
	RequestURL  string
	RequestBody []byte
}

// ResponseData 画面に送信するデータ構造
type ResponseData struct {
	ResponseSumally string
	ResponseBody    []byte
}

type server struct {
	// sumally用ソケット
	socket *websocket.Conn
	// リクエスト待ち受けチャネル
	request chan *RequestData
	// レスポンス待ち受けチャネル
	response chan *ResponseData
}

func (s *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("ServeHttp server")
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	s.socket = socket

	// ソケット待ち受け開始
	go s.write()
	s.read()
}

func (s *server) read() {
	for {
		msg := &RequestData{}
		if err := s.socket.ReadJSON(msg); err == nil {
			log.Printf("read server: %v", msg.RequestURL)
			s.request <- msg
		} else {
			break
		}
	}
	s.socket.Close()
}

func (s *server) write() {
	for response := range s.response {
		if err := s.socket.WriteJSON(response); err != nil {
			break
		}
	}
	s.socket.Close()
}

func (s *server) socketIOStart() {
	// ソケットIO発生時の動作定義
	log.Println("socketIOStart")
	for {
		select {
		case request := <-s.request:
			// POSTリクエスト
			log.Printf("request socketIOStart: %v", request.RequestURL)
			s.response <- request.post()
		}
	}
}

func (req *RequestData) post() (res *ResponseData) {
	//httpからCliantインスタンス取得
	//Cliant構造体にはhttp cliantに必要なメソッドが備わっている
	log.Println("post server")
	cliant := &http.Client{}

	postBody := bytes.NewBuffer(req.RequestBody)
	log.Println(postBody)
	response, err := cliant.Post(req.RequestURL, "application/json;charset=UTF-8", postBody)
	log.Println("ここまできた")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// if res.StatusCode != http.StatusOK {
	// 	log.Fatal(err)
	// }

	// Response body read
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	res.ResponseBody = body
	res.ResponseSumally = "StatusCode:" + response.Status + " Content-Type:" + response.Header.Get("Content-Type")
	log.Printf("ここまできた:%v", res.ResponseSumally)
	return res
}
