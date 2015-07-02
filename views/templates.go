package views

import "text/template"

func Templates() map[string]*template.Template {
	var T = make(map[string]*template.Template)

	T["users/register"] = template.Must(template.ParseFiles("views/users/base.tmpl", "views/users/register.tmpl"))

	return T
}
