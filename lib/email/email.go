package email

import (
	"strconv"
	"net/smtp"
)

var sender = smtp.SendMail

type Email struct {
	Password string
	Server string
	Port int
}

func (e *Email)FullAddr() string {
	return e.Server + ":" + strconv.Itoa(e.Port)
}

func (e *Email)Send(from, to, subject, body string) error {
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		body

	return sender(e.FullAddr(),
		smtp.PlainAuth("", from, e.Password, e.Server),
		from, []string{to}, []byte(msg))
}

func GenBody(template string, data interface{}) string {
	f, err := os.Open(tempalte)
	if err != nil {
		return ""
	}
	defer f.Close()
	bodytemp, err := ioutils.ReadAll(f)
	if err != nil {
		return ""
	}
	t:= template.Must(template.New("genbody").Parse(bodytemp))
	return data.(string)
}
