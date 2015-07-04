package http

import (
	"encoding/json"
	"log"
	"net/http"
	_ "net/http/pprof"
	"github.com/baiyuxiong/gomail/config"
)

type JsonResult struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Start() {
	startHttpServer()
}

func startHttpServer() {

	addr := config.Config().HttpAddress
	if addr == "" {
		return
	}

	configRoutes()

	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}

	log.Println("http.startHttpServer ok, listening", addr)
	log.Fatalln(s.ListenAndServe())
}

func RenderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

func RenderOKJson(w http.ResponseWriter, data interface{}) {
	RenderJson(w, &JsonResult{
		Code:    200,
		Data:    data,
		Message: "",
	})
}

func RenderErrorJson(w http.ResponseWriter, msg string) {
	RenderJson(w, &JsonResult{
		Code:    404,
		Data:    nil,
		Message: msg,
	})
}



