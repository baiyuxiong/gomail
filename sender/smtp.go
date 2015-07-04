package sender

import (
	"github.com/baiyuxiong/gomail/model"
	"github.com/baiyuxiong/gomail/config"
	"encoding/base64"
	"time"
	"strings"
	"net/smtp"
	"encoding/json"
	"github.com/baiyuxiong/gomail/constants"
	"log"
	"fmt"
)

var SMTPSender = &Sender{
	Name:"smtp",
	Run:sendEmailBySMTP,
}

func sendEmailBySMTP(data string) model.EmailLog{
	log.Println("sendEmailBySMTP ", data)

	l := model.EmailLog{
		Status:constants.MAIL_ERR,
		SendTime:time.Now(),
		LogMessage:"",
	}

	var m model.Email
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		l.LogMessage = "handleJob - Unmarshal err: " + err.Error()+ ", data : "+ data
		return l
	}
	l.Mail = m

	err = SendToMail(m.To,m.Subject,m.Message)

	if err != nil{
		l.LogMessage = err.Error()
		return l
	}

	l.Status = constants.MAIL_OK
	log.Println("sendEmailBySMTP OK")

	return l
}


func SendToMail(to, subject, body string) error {
	hp := strings.Split(config.Config().Smtp.Address, ":")
	auth := smtp.PlainAuth("", config.Config().Smtp.Username, config.Config().Smtp.Password, hp[0])

	send_to := strings.Split(to, ";")
	message := ""
	message += fmt.Sprintf("%s: %s\r\n", "Subject", subject)
	message += fmt.Sprintf("%s: %s\r\n", "Content-Transfer-Encoding", "base64")
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	err := smtp.SendMail(config.Config().Smtp.Address, auth, config.Config().Email.From, send_to, []byte(message))
	if err!= nil{
		log.Println("SendToMail",config.Config().Smtp.Address,"with username",config.Config().Smtp.Username,"err : ", err.Error())
	}
	return err
}