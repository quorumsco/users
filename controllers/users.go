package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	. "github.com/quorumsco/jsonapi"
	"github.com/quorumsco/logs"
	"github.com/quorumsco/users/models"
	"github.com/quorumsco/users/views"
)

func sPtr(s string) *string {
	if s == "" {
		return nil
	} else {
		return &s
	}
}

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

func Auth(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logs.Error(err)
		Fail(w, r, map[string]string{"Authentification": "Error"}, http.StatusBadRequest)
		return
	}
	infos := make(map[string]interface{})
	if err := json.Unmarshal(body, &infos); err != nil {
		logs.Error(err)
		Fail(w, r, map[string]string{"Authentification": "Bad parameters"}, http.StatusBadRequest)
		return
	}
	username := infos["username"].(string)
	password := infos["password"].(string)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	var (
		u         = models.User{Mail: &username, Password: sPtr(string(passwordHash))}
		db        = getDB(r)
		userStore = models.UserStore(db)
	)
	if err = userStore.First(&u); err != nil {
		logs.Error(err)
		Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	if u.GroupID == 0 {
		Fail(w, r, map[string]interface{}{"User": "No such user"}, http.StatusBadRequest)
		return
	}
	Success(w, r, views.User{User: &u}, http.StatusOK)
}
