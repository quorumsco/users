package models

import "github.com/jinzhu/gorm"

type UserDS interface {
	Save(*User) error
	Delete(*User) error
	First(*User) error
	Find() ([]User, error)
}

type UserSQL struct {
	DB *gorm.DB
}

func UserStore(db *gorm.DB) UserDS {
	return &UserSQL{DB: db}
}

func (s *UserSQL) Save(u *User) error {
	if u.ID == 0 {
		s.DB.Create(u)

		return s.DB.Error
	}

	s.DB.Save(u)

	return s.DB.Error
}

func (s *UserSQL) Delete(u *User) error {
	s.DB.Delete(u)

	return s.DB.Error
}

func (s *UserSQL) First(u *User) error {
	s.DB.Find(u)

	return s.DB.Error
}

func (s *UserSQL) Find() ([]User, error) {
	var users []User
	s.DB.Find(&users)
	if s.DB.Error != nil {
		return users, nil
	}
	return users, s.DB.Error
}
