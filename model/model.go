package model

import (
	"time"
	"github.com/baiyuxiong/gomail/constants"
)

type Email struct {
	To      string `json:"to"`
	CC      string `json:"cc"`
	BCC     string `json:"bcc"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type EmailSendCloud struct {
	Email
	MailType          string `json:"mailType"` //trigger batch
	Substitution_vars string `json:"substitution_vars"`
	Template_invoke_name string `json:"template_invoke_name"`
}

type EmailLog struct {
	Status     constants.MailStatus `json:"status"`
	Mail       interface{} `json:"mail"`
	SendTime   time.Time    `json:"sendTime"`
	LogMessage string `json:"logMessage"`
}

type MailResp struct {
	Message string `json:"message"`
}