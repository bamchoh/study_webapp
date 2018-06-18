package email

import (
	"net/smtp"
	"testing"
	"github.com/bamchoh/study_webapp/lib/utils"
	gen "github.com/bamchoh/study_webapp/lib/jwtgen"
)

type User struct {
	Name string `json:"name"`
	Age int `json:"age"`
}

type EmailInfo struct {
	Addr string
	From string
	To []string
	Msg []byte
}

func TestSendEmail(t *testing.T) {
	secretkey := utils.SecureRandomStr(64)

	user1 := User {
		Name: "bamchoh",
		Age: 100,
	}

	tokenstring, err := gen.Generate(user1, secretkey)
	if err != nil {
		t.Error(err)
	}

	sender = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		return nil
	}

	email := Email {
		Password: "aivbhjftleabzefv",
		Server: "smtp.gmail.com",
		Port: 587,
	}

	from := "bamchoh@gmail.com"
	to := "bamchoh@gmail.com"
	title := "やぁ"
	email.Send(from, to, title, tokenstring)

	user2 := User{}
	err = gen.Parse(tokenstring, &user2, secretkey)
	if err != nil {
		t.Error(err)
	}

	if !(user1.Name == user2.Name && user1.Age == user2.Age) {
		t.Errorf("user is different between generate and parse:\n Got:%v,%v\nWant:%v,%v",
		user1.Name, user1.Age,
		user2.Name, user2.Age)
	}
}

