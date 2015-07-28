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
		var err error
		if s.DB.DriverName() == "postgres" {
			var result *sqlx.Rows
			result, err = s.DB.NamedQuery("INSERT INTO users (firstname, surname, mail, password) VALUES (:firstname, :surname, :mail, :password) RETURNING id", u)
			result.Scan(&u.ID)
		} else {
			var result sql.Result
			result, err = s.DB.NamedExec("INSERT INTO users (firstname, surname, mail, password) VALUES (:firstname, :surname, :mail, :password)", u)
			u.ID, err = result.LastInsertId()
		}
		return err
	}

	_, err := s.DB.NamedExec("UPDATE users SET firstname=:firstname, surname=:surname, mail=:mail, password=:password WHERE id=:id", u)
	return err
}

func (s *UserSQL) Delete(u *User) error {
	_, err := s.DB.NamedExec("DELETE FROM users WHERE id=:id", u)
	return err
}

func (s *UserSQL) First(u *User) error {
	err := s.DB.Get(u, s.DB.Rebind("SELECT * FROM users WHERE id=? LIMIT 1"), u.ID)
	return err
}

func (s *UserSQL) Find() ([]User, error) {
	var users []User
	err := s.DB.Select(&users, "SELECT id, firstname, surname FROM users ORDER BY surname DESC")
	if err == sql.ErrNoRows {
		return users, nil
	}
	return users, err
}
