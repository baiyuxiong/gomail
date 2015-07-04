package http

import (
	"html/template"
	"time"
	"net/http"
	"path/filepath"
)

var templateDir = "templates"

var funcMap = template.FuncMap{
	"formatTime" : func(date time.Time, format string) string {
		return date.Format(format)
	},
	"datetime" : func(date time.Time) string {
		return date.Format("2006-01-02 15:04")
	},
}

// 解析模板
func Render(w http.ResponseWriter, file string, data map[string]interface{}){

	fileName := filepath.Join(templateDir, file)
	t, err := template.New(filepath.Base(fileName)).Funcs(funcMap).ParseFiles(fileName,filepath.Join(templateDir, "header.html"),filepath.Join(templateDir, "footer.html"),filepath.Join(templateDir, "nav.html"))

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	t.Execute(w, data)
}