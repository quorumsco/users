package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type UserDS interface {
	Save(*User) error
	Delete(*User) error
	First(*User) error
	Find() ([]User, error)
}

type UserSQL struct {
	DB *sqlx.DB
}

func UserStore(db *sqlx.DB) UserDS {
	return &UserSQL{DB: db}
}

func (s *UserSQL) Save(u *User) error {
	if u.ID == 0 {
		result, err := s.DB.NamedExec("INSERT INTO users (firstname, surname, mail, password) VALUES (:firstname, :surname, :mail, :password)", u)
		if err != nil {
			return err
		}

		u.ID, err = result.LastInsertId()
		return err
	}

	_, err := s.DB.NamedExec("UPDATE users SET firstname=:firstname, surname=:surname, mail=:mail, password=:password WHERE id=:id", u)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserSQL) Delete(u *User) error {
	_, err := s.DB.NamedExec("DELETE FROM users WHERE id=:id", u)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserSQL) First(u *User) error {
	if err := s.DB.Get(u, "SELECT * FROM users WHERE id=? LIMIT 1", u.ID); err != nil {
		return err
	}
	return nil
}

func (s *UserSQL) Find() ([]User, error) {
	var users []User
	if err := s.DB.Select(&users, "SELECT id, firstname, surname FROM users ORDER BY surname DESC"); err != nil {
		if err == sql.ErrNoRows {
			return users, nil
		}
		return nil, err
	}
	return users, nil
}
