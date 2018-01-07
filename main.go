package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
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
	t.temple.Execute(w, r)
}
func main() {
	http.Handle("/", &templateHandler{filename: "index.html"})

	// webサーバー開始
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
