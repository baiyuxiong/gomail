package sender

import (
	"github.com/baiyuxiong/gomail/model"
	"time"
	"net/url"
	"net/http"
	"github.com/baiyuxiong/gomail/config"
	"encoding/json"
	"github.com/baiyuxiong/gomail/constants"
	"io/ioutil"
	"log"
)

var SendcloudSender = &Sender{
	Name:"sendcloud",
	Run:sendEmailBySendcloud,
}

func sendEmailBySendcloud(data string) model.EmailLog {
	log.Println("sendEmailBySendcloud...")

	l := model.EmailLog{
		Status:constants.MAIL_ERR,
		SendTime:time.Now(),
		LogMessage:"",
	}

	var m model.EmailSendCloud
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		l.LogMessage = "handleJob - Unmarshal err: " + err.Error()+ ", data : "+ data
		log.Println("sendEmailBySendcloud err" , l.LogMessage)
		return l
	}
	l.Mail = m

	values := url.Values{}
	values["api_user"] = []string{config.Config().Sendcloud.Api_user}
	values["api_key"] = []string{config.Config().Sendcloud.Api_key}
	values["fromname"] = []string{config.Config().Email.Fromname}
	values["replyto"] = []string{config.Config().Email.Replyto}
	values["from"] = []string{config.Config().Email.From}
	values["to"] = []string{m.To}
	values["subject"] = []string{m.Subject}
	values["html"] = []string{m.Message}
	values["bcc"] = []string{m.BCC}
	values["cc"] = []string{m.CC}

	//http://sendcloud.sohu.com/doc/email/#x-smtpapi
	values["x_smtpapi"] = []string{m.X_smtpapi}

	resp, err := http.PostForm(config.Config().Sendcloud.Address, values)

	body, err := ioutil.ReadAll(resp.Body)
	log.Println("sendEmailBySendcloud resp" , string(body))

	if err != nil {
		l.LogMessage = "err: " + err.Error() + ", resp: " + string(body)
		log.Println("sendEmailBySendcloud err" , l.LogMessage)
	}
	defer resp.Body.Close()

	l.Status = constants.MAIL_OK
	log.Println("sendEmailBySendcloud OK")
	return l
}