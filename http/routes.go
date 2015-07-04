package http

import (
	"fmt"
	"github.com/toolkits/file"
	"github.com/baiyuxiong/gomail/config"
	"github.com/baiyuxiong/gomail/mail"
	"net/http"
	"strings"
	"github.com/baiyuxiong/gomail/model"
	"encoding/json"
	"log"
)

func configRoutes() {
	// GET
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logs,err := mail.GetMailLog()
		if err != nil{
			w.Write([]byte(err.Error()))
			return
		}
		var mailLogs []model.EmailLog

		for _,l := range logs.([]interface{}){
			var m model.EmailLog
			err := json.Unmarshal(l.([]byte), &m)
			if err != nil {
				log.Println("HandleFunc: / Unmarshal error - " , err.Error())
			}else{
				mailLogs  = append(mailLogs, m)
			}
		}

		Render(w,"home/index.html",map[string]interface{}{"mailLogs":mailLogs});
	})

	// GET
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok\n"))
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%s\n", config.VERSION)))
	})

	http.HandleFunc("/workdir", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%s\n", file.SelfDir())))
	})

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		RenderOKJson(w, config.Config())
	})

	http.HandleFunc("/config/reload", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RemoteAddr, "127.0.0.1") {
			config.Parse(config.ConfigFile)
			RenderOKJson(w, config.Config())
		} else {
			w.Write([]byte("no privilege\n"))
		}
	})
	http.HandleFunc("/public/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
}
