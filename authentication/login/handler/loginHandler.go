package handler

import (
	"html/template"
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	//title := r.URL.Path[len("/edit/"):]
	//p, err := loadPage(title)
	//if err != nil {
	//p = &Page{Title: title}
	//}
	t, _ := template.ParseFiles("edit.html")
	t.Execute(w, "he")
}
