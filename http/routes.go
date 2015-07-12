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
	"net/url"
	"strconv"
	"time"
)

var lastErrorTime time.Time = time.Now().Add(time.Minute*-2)

func configRoutes() {
	// GET
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		queries, err := url.ParseQuery(r.URL.RawQuery)

		if !lastErrorTime.IsZero() && (time.Now().Sub(lastErrorTime) < time.Minute){
			w.Write([]byte("Pls try after " + time.Now().Add(time.Minute).String()))
			return
		}

		if (len(queries["password"]) < 1){
			w.Write([]byte("Need password."))
			return
		}

		password := queries["password"][0]
		if password != config.Config().Password{
			lastErrorTime = time.Now()
			w.Write([]byte("I love you, Pls don't hurt me."))
			return
		}

		var perPage int = 20

		var start int = 0
		if err == nil && len(queries["start"]) > 0 {
			start,_ = strconv.Atoi(queries["start"][0])
		}

		var stop int = start+perPage-1
		if err == nil && len(queries["stop"]) > 0 {
			stop,_ = strconv.Atoi(queries["stop"][0])
		}

		var currentPage int = start/perPage+1;

		var prePageStart int= 0
		var prePageStop int= 19
		var nextPageStart int= 20
		var nextPageStop int= 39

		var showPrePage = false;
		if currentPage>1{
			showPrePage= true
			prePageStart = (currentPage-2)*perPage
			prePageStop = prePageStart +perPage-1
		}
		nextPageStart = currentPage*perPage
		nextPageStop = nextPageStart +perPage-1


		logs,err := mail.GetMailLog(start,stop)
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

		data := map[string]interface{}{
			"mailLogs":mailLogs,
			"start":start,
			"stop":stop,
			"currentPage":currentPage,
			"showPrePage":showPrePage,
			"prePageStart":prePageStart,
			"prePageStop":prePageStop,
			"nextPageStart":nextPageStart,
			"nextPageStop":nextPageStop,
		}

		Render(w,"home/index.html",data);
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
