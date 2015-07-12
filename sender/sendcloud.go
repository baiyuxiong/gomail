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
	"strings"
)

var SendcloudSender = &Sender{
	Name:"sendcloud",
	Run:sendEmailBySendcloud,
}

func sendEmailBySendcloud(data string) (l model.EmailLog) {
	log.Println("sendEmailBySendcloud...")

	l = model.EmailLog{
		Status:constants.MAIL_ERR,
		SendTime:time.Now(),
		LogMessage:"",
	}

	var m model.EmailSendCloud
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		l.LogMessage = "handleJob - Unmarshal err: " + err.Error()+ ", data : "+ data
		log.Println("sendEmailBySendcloud err" , l.LogMessage)
		return
	}
	l.Mail = m

	api_user := config.Config().Sendcloud.Api_user_trigger
	if m.MailType == constants.MAIL_TYPE_BATCH{
		api_user = config.Config().Sendcloud.Api_user_batch
	}
	log.Println("sendEmailBySendcloud - Using mail type " , api_user)


	values := url.Values{}
	values["api_user"] = []string{api_user}
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
	if len(m.X_smtpapi) > 0{
		values["x_smtpapi"] = []string{m.X_smtpapi}
	}

	log.Println("sendEmailBySendcloud - values " , values)

	client := &http.Client{}
	req,err := http.NewRequest("POST", config.Config().Sendcloud.Address, strings.NewReader(values.Encode()))
	if err != nil{
		log.Println("sendEmailBySendcloud NewRequest err" , err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp,err := client.Do(req)
	if err != nil {
		log.Println("sendEmailBySendcloud read resp error" , err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	log.Println("sendEmailBySendcloud resp" , string(body))
	l.LogMessage = "server resp: " + string(body)

	if err != nil {
		log.Println("sendEmailBySendcloud err" , l.LogMessage)
	}
	defer resp.Body.Close()

	var result model.MailResp
	err = json.Unmarshal(body,&m)
	if err != nil{
		log.Println("sendEmailBySendcloud Unmarshal resp error : " , err.Error())
		l.Status = constants.MAIL_ERR
		return
	}

	if result.Message == "success"{
		l.Status = constants.MAIL_OK
		log.Println("sendEmailBySendcloud OK")
		return
	}else{
		l.Status = constants.MAIL_ERR
		return
	}
}