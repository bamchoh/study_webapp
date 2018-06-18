package dao

import (
	"database/sql"
	"errors"

	models "github.com/bamchoh/study_webapp/models"

	validator "gopkg.in/go-playground/validator.v9"
)

type UserDao struct {
	DB *sql.DB
}

func (m *UserDao) ValidateForCreation(user models.User, password, confirmpassword string) (err error) {
	v := validator.New()

	if err := v.Struct(user); err != nil {
		return err
	}

	if err := v.VarWithValue(password, confirmpassword, "eqfield"); err != nil {
		return err
	}

	if err := v.Var(password, "min=8"); err != nil {
		return err
	}

	return nil
}

func (m *UserDao) Create(user models.User, password, confirmpassword string) (err error) {
	err = m.ValidateForCreation(user, password, confirmpassword)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (id, name, email, password_hash) VALUES ($1, $2, $3, crypt($4, gen_salt('bf')))`
	_, err = m.DB.Query(stmt, user.ID, user.Name, user.Email, password)
	return err
}

func (m *UserDao) Get(email, pass string) (user models.User, err error) {
	var hash string
	hash, err = m.getPasswordHash(email, pass)
	if err != nil || hash == "" {
		err = errors.New("user name or password are invalid")
		return user, err
	}

	rows, err := m.DB.Query("SELECT id,name,email FROM users WHERE email = $1 AND password_hash = $2 AND activated = true", email, hash)
	if err != nil {
		return user, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		return user, err
	}
	return user, err
}

func (m *UserDao) getPasswordHash(email, pass string) (hash string, err error) {
	rows, err := m.DB.Query(`SELECT crypt($1, password_hash) = password_hash as matched, password_hash FROM users WHERE email = $2`, pass, email)
	if err != nil {
		return hash, err
	}
	defer rows.Close()

	for rows.Next() {
		var result string
		if err := rows.Scan(&result, &hash); err != nil {
			return hash, err
		}
		if result == "true" {
			return hash, nil
		}
	}
	return "", nil
}
