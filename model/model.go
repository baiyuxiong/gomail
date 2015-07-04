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
	X_smtpapi string `json:"x_smtpapi"`
}

type EmailLog struct {
	Status     constants.MailStatus `json:"status"`
	Mail       interface{} `json:"mail"`
	SendTime   time.Time    `json:"sendTime"`
	LogMessage string `json:"logMessage"`
}