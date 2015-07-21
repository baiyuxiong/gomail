package http

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/toolkits/file"
	"github.com/baiyuxiong/gomail/config"
	"github.com/baiyuxiong/gomail/mail"
	"net/http"
	"github.com/baiyuxiong/gomail/model"
	"encoding/json"
	"log"
	"net/url"
	"strconv"
	"time"
	"errors"
)

var lastErrorTime time.Time = time.Now().Add(time.Minute*-2)

var store=sessions.NewCookieStore([]byte("70df6669f6414d0ce82cbf46fffe0741"))
var sessionName = "gomail"

func configRoutes() {
	// GET
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if nil != checkLoggedIn(w,r){
			return
		}

		queries, err := url.ParseQuery(r.URL.RawQuery)

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
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		Render(w,"home/login.html",nil);
	})

	http.HandleFunc("/doLogin", func(w http.ResponseWriter, r *http.Request) {
		if !lastErrorTime.IsZero() && (time.Now().Sub(lastErrorTime) < time.Minute){
			w.Write([]byte("Pls try after " + time.Now().Add(time.Minute).String()))
			return
		}

		r.ParseForm();
		if (len(r.Form["password"]) < 1){
			w.Write([]byte("Need password."))
			return
		}
		if (len(r.Form["username"]) < 1){
			w.Write([]byte("Need username."))
			return
		}

		password := r.Form["password"][0]
		username := r.Form["username"][0]

		if password != config.Config().Password || username != config.Config().Username{
			lastErrorTime = time.Now()
			w.Write([]byte("I love you, Pls don't hurt me."))
			return
		}else{
			storeUsername(w,r,username);
			http.Redirect(w,r,"/", http.StatusFound)
			return
		}
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		storeUsername(w,r,"");
		http.Redirect(w,r,"/login", http.StatusFound)
		return
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok\n"))
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		if nil != checkLoggedIn(w,r){
			return
		}
		w.Write([]byte(fmt.Sprintf("%s\n", config.VERSION)))
	})

	http.HandleFunc("/workdir", func(w http.ResponseWriter, r *http.Request) {
		if nil != checkLoggedIn(w,r){
			return
		}
		w.Write([]byte(fmt.Sprintf("%s\n", file.SelfDir())))
	})

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		if nil != checkLoggedIn(w,r){
			return
		}
		RenderOKJson(w, config.Config())
	})

	http.HandleFunc("/config/reload", func(w http.ResponseWriter, r *http.Request) {
		if nil != checkLoggedIn(w,r){
			return
		}
		config.Parse(config.ConfigFile)
		RenderOKJson(w, config.Config())
	})
	http.HandleFunc("/public/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
}

func checkLoggedIn(w http.ResponseWriter, r *http.Request) error{
	session,_:=store.Get(r,sessionName)
	username,ok :=session.Values["username"].(string)
	if !ok || len(username) == 0{
		http.Redirect(w,r,"/login", http.StatusFound)
		return errors.New("Login first")
	}
	return nil
}

func storeUsername(w http.ResponseWriter, r *http.Request,username string) {
	session,_:=store.Get(r,sessionName)
	session.Values["username"] = username
	session.Save(r,w)
}
