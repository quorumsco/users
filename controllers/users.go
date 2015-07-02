package controllers

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	. "github.com/iogo-framework/jsonapi"
	"github.com/quorumsco/users/models"
	"github.com/quorumsco/users/views"

	"github.com/iogo-framework/logs"
)

func sPtr(s string) *string { return &s }

func Register(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		store := models.UserStore(getDB(req))

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

		err = store.Save(u)
		if err != nil {
			logs.Error(err)
			return
		}
	}
}

func CreateUser(w http.ResponseWriter, req *http.Request) {
	db := getDB(req)
	store := models.UserStore(db)

	var u = new(models.User)
	err := Request(&views.User{User: u}, req)
	if err != nil {
		logs.Debug(err)
		Fail(w, req, map[string]interface{}{"contact": err.Error()}, http.StatusBadRequest)
		return
	}

	errs := u.Validate()
	if len(errs) > 0 {
		logs.Debug(errs)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		Fail(w, req, map[string]interface{}{"contact": errs}, http.StatusBadRequest)
		return
	}

	err = store.Save(u)
	if err != nil {
		logs.Error(err)
		Error(w, req, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/%s/%d", "contacts", u.ID))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	Success(w, req, views.User{User: u}, http.StatusCreated)
}
