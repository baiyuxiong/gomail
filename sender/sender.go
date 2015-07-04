package sender

import (
	"github.com/baiyuxiong/gomail/model"
)

type Sender struct {
	Name string
	Run func(string) model.EmailLog
}