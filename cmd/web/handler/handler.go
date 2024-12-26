package handler

import (
	"example.com/go-web-base/internal/application"
	"html/template"
	"net/http"
)

type BaseHandler struct {
	App application.Application
}

func (h BaseHandler) IndexPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/base.gohtml", "templates/index.gohtml"))

	err := tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, "Error rendering template :"+err.Error(), http.StatusInternalServerError)
	}
}
