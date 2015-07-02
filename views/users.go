package views

import "github.com/quorumsco/users/models"

type Users struct {
	Users []models.User `json:"users"`
}

type User struct {
	User *models.User `json:"user"`
}
