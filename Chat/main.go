package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

// templateHandler型を宣言
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// templateHandler型のメソッドとしてServeHTTPを作成
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//sync.Onceにより複数のgoroutineがServeHTTPを呼び出してもt.once.Doの引数に渡した関数が一度しか実行されない
	t.once.Do(func() {
		//テンプレートをコンパイル
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	t.templ.Execute(w, r)

}

func main() {

	var addres = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()
	//ROOMの新規作成
	r := newRoom()
	//ルート
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	//チャットルームの開始
	go r.run()
	//webサーバーの開始
	log.Println("webサーバーを開始しますポートは", *addres)
	if err := http.ListenAndServe(*addres, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
