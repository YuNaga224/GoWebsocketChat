package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/YuNaga224/websocketChat/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

// config.jsonから受け取ったクライアントID,クライアントシークレットを格納する構造体
type GoogleOAuth struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type GitHubOAuth struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type OAuthProvider struct {
	Google GoogleOAuth `json:"google"`
	GitHub GitHubOAuth `json:"github"`
}

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
	//config.jsonから設定情報読み込み
	configBytes, err := ioutil.ReadFile("../config.json")
	if err != nil {
		fmt.Println("config.jsonファイルの読み込みでエラーが発生しました", err)
	}

	var oauthProvider OAuthProvider

	err = json.Unmarshal(configBytes, &oauthProvider)
	if err != nil {
		fmt.Println("configファイルの解析でエラーが発生しました", err)
	}

	var addres = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()

	//Gomniauthセットアップ
	gomniauth.SetSecurityKey("デジタル署名となるランダムな値")
	gomniauth.WithProviders(
		facebook.New("クライアントID", "秘密の値", "http://localohost:8080/auth/callback/facebook"),
		github.New(oauthProvider.GitHub.ClientID, oauthProvider.GitHub.ClientSecret, "http://localhost:8080/auth/callback/github"),
		google.New(oauthProvider.Google.ClientID, oauthProvider.Google.ClientSecret, "http://localhost:8080/auth/callback/google"),
	)
	//ROOMの新規作成
	r := newRoom(UseGravatar)
	r.tracer = trace.New(os.Stdout)
	//ルート
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/uploader", uploaderHandler)
	http.Handle("/room", r)
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	//ファイルサーバを書くユーザのブラウザに提供
	http.Handle("/avatars/",http.StripPrefix("/avatars/",http.FileServer(http.Dir("./avatars"))))
	//チャットルームの開始
	go r.run()
	//webサーバーの開始
	log.Println("webサーバーを開始しますポートは", *addres)
	if err := http.ListenAndServe(*addres, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
