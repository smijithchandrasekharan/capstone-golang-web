package handler

import (
	"html/template"
	"net/http"

)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("edit.html")
	t.Execute(w, r.URL.RawPath)

	//
}
