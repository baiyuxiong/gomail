package mail

import (
	"github.com/baiyuxiong/gomail/sender"
	"github.com/baiyuxiong/gomail/maillog"
	"github.com/baiyuxiong/gomail/config"
)

var senders = []*sender.Sender{
	sender.SendcloudSender,
	sender.SMTPSender,
}

func send(data string) {
	for _,s := range senders {
		if s.Name == config.Config().Sender{
			l := s.Run(data)
			maillog.PutLog(LogConn,l)
			break;
		}
	}
}
