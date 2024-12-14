package handler

import (
	"html/template"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./template/loginForm.html")
	t.Execute(w, "")

	//
}
