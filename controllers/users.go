package controllers

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/quorumsco/logs"
	"github.com/quorumsco/users/models"
)

func sPtr(s string) *string { return &s }

func Register(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.ParseForm()

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.FormValue("password")), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}

		u := &models.User{
			Firstname: sPtr(req.FormValue("firstname")),
			Surname:   sPtr(req.FormValue("surname")),
			Mail:      sPtr(req.FormValue("mail")),
			Password:  sPtr(string(passwordHash)),
		}

		errs := u.Validate()
		if len(errs) > 0 {
			logs.Debug(errs)
			return
		}

		var store = models.UserStore(getDB(req))
		err = store.Save(u)
		if err != nil {
			logs.Error(err)
			return
		}
	}

	templates := getTemplates(req)
	if err := templates["users/register"].ExecuteTemplate(w, "base", nil); err != nil {
		logs.Error(err)
	}
}
