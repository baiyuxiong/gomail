package constants

type MailStatus uint

const (
	MAIL_OK MailStatus = iota
	MAIL_ERR
)

var MailStatusStrings = [...]string{
	"OK",
	"ERR",
}

const (
	MAIL_TYPE_TRIGGER string = "trigger" //触发邮件
	MAIL_TYPE_BATCH string = "batch"	//批量邮件
)

func (s MailStatus) String() string{
	return MailStatusStrings[s]
}