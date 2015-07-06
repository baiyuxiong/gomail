package mail

import (
	"github.com/baiyuxiong/gomail/sender"
)

var senders = []*sender.Sender{
	sender.SendcloudSender,
	sender.SMTPSender,
}