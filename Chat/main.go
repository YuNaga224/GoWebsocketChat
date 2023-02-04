package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/YuNaga224/websocketChat/websocketChat/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
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
	data := map[string]interface{}{
		"Host": r.Host,
	}

	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, data)

}

func main() {

	var addres = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()

	//Gomniauthセットアップ
	gomniauth.SetSecurityKey("デジタル署名となるランダムな値")
	gomniauth.WithProviders(
		facebook.New("クライアントID", "秘密の値", "http://localohost:8080/auth/callback/facebook"),
		github.New("88201b7004c5b9ceab8f", "7f2a5c9bbfce6013f29453f6f32e08d93219e5cd", "http://localhost:8080/auth/callback/github"),
		google.New("515360048655-070u177uq2879ol0jafs95g47ov1eiau.apps.googleusercontent.com", "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", "http://localhost:8080/auth/callback/google"),
	)
	//ROOMの新規作成
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	//ルート
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	//チャットルームの開始
	go r.run()
	//webサーバーの開始
	log.Println("webサーバーを開始しますポートは", *addres)
	if err := http.ListenAndServe(*addres, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
