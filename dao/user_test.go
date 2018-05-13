package dao

import (
	"database/sql"
	"fmt"
	"testing"

	models "github.com/bamchoh/study_webapp/models"
	_ "github.com/lib/pq"
)

type Test struct {
	id    string
	name  string
	email string
	pass  string
	conf  string
}

func setup() (*sql.DB, error) {
	var err error
	db, err := sql.Open("postgres", "postgres://postgres:yoshihiko_yamanaka_bamchoh@localhost:5432/test_heroku_db?sslmode=disable")
	if err != nil {
		return nil, err
	}

	_, err = db.Query("TRUNCATE TABLE users")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestCreateNormal(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Fatalf("Error opening database: %q", err)
	}

	tests := []Test{
		{
			id:    "test_id",
			name:  "test_name",
			email: "test@email.com",
			pass:  "test_password",
			conf:  "test_password",
		},
		{
			id:    "test_id2",
			name:  "あいうえおかきくけこさしすせそたちつてと12",
			email: "test2@email.com",
			pass:  "test_password",
			conf:  "test_password",
		},
	}

	builder := &UserDao{db}

	for _, tst := range tests {
		user := models.User{
			ID:    tst.id,
			Name:  tst.name,
			Email: tst.email,
		}
		if err = builder.Create(user, tst.pass, tst.conf); err != nil {
			t.Fatalf("[Create] %q", err)
		}

		got, err := builder.Get(user.Email, tst.pass)
		if err != nil {
			t.Fatalf("[Get] %q", err)
		}

		if user.Email != got.Email {
			t.Errorf("Email is different\nwant:%v\ngot :%v", user.Email, got.Email)
		}

		if user.Name != got.Name {
			t.Errorf("Name is different\nwant:%v\ngot :%v", user.Name, got.Name)
		}
	}
}

func TestCreateError(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Fatalf("Error opening database: %q", err)
	}

	builder := &UserDao{db}

	tests := []Test{
		{
			id:    "test_id",
			name:  "test_name",
			email: "test@email.com",
			pass:  "test_password",
			conf:  "test_password1",
		},
		{
			id:    "test_id",
			name:  "test_name",
			email: "test_email.com",
			pass:  "test_password",
			conf:  "test_password",
		},
		{
			id:    "",
			name:  "test_name",
			email: "test@email.com",
			pass:  "test_password",
			conf:  "test_password",
		},
		{
			id:    "test_id",
			name:  "",
			email: "test@email.com",
			pass:  "test_password",
			conf:  "test_password",
		},
		{
			id:    "test_id1234512345123451234512345",
			name:  "test_name",
			email: "test@email.com",
			pass:  "",
			conf:  "",
		},
		{
			id:    "test_id",
			name:  "test_name",
			email: "",
			pass:  "test_password",
			conf:  "test_password",
		},
		{
			id:    "test_id",
			name:  "32chars_name_abcdefg1234567890123",
			email: "test@email.com",
			pass:  "test_password",
			conf:  "test_password",
		},
	}

	for _, tst := range tests {
		title := fmt.Sprintf("id:%v,name:%v,email:%v,pass:%v,conf:%v",
			tst.id, tst.name, tst.email, tst.pass, tst.conf)
		t.Run(title, func(t *testing.T) {
			user := models.User{
				ID:    tst.id,
				Name:  tst.name,
				Email: tst.email,
			}
			password := tst.pass
			confirmpassword := tst.conf

			if err = builder.Create(user, password, confirmpassword); err == nil {
				t.Fatalf("[Create] Creation was succeeded even though invalid parameters")
			}

			_, err = builder.Get(user.Email, password)
			if err == nil {
				t.Fatalf("[Create] Creation was succeeded even though invalid parameters")
			}
		})
	}
}

func TestDuplicationCreation(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Fatalf("Error opening database: %q", err)
	}

	builder := &UserDao{db}

	user1 := models.User{
		ID:    "test_id",
		Name:  "test_name1",
		Email: "test1@email.com",
	}

	user2 := models.User{
		ID:    "test_id",
		Name:  "test_name2",
		Email: "test2@email.com",
	}

	if err = builder.Create(user1, "testtest", "testtest"); err != nil {
		t.Fatalf("[Create] Creation was succeeded even though invalid parameters")
	}

	if err = builder.Create(user2, "testtest", "testtest"); err == nil {
		t.Fatalf("[Create] Creation was succeeded even though invalid parameters")
	}
}
