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

func (s MailStatus) String() string{
	return MailStatusStrings[s]
}