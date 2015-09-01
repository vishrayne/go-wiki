package main

import (
	"html/template"
	"net/http"
)

// caching templates
var templates = template.Must(template.ParseFiles(ViewPath("edit"), ViewPath("view")))

// generic template renderer
func renderTemplate(rw http.ResponseWriter, view string, p *Page) {
	err := templates.ExecuteTemplate(rw, view+".html", p)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

// Request handlers
func viewHandler(rw http.ResponseWriter, req *http.Request) {
	title, err := GetTitle(rw, req)
	if err != nil {
		return
	}

	page, err := LoadPage(title)
	if err != nil {
		http.Redirect(rw, req, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(rw, "view", page)
}

func editHandler(rw http.ResponseWriter, req *http.Request) {
	title, err := GetTitle(rw, req)
	if err != nil {
		return
	}

	page, err := LoadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}

	renderTemplate(rw, "edit", page)
}

func saveHandler(rw http.ResponseWriter, req *http.Request) {
	title, err := GetTitle(rw, req)
	if err != nil {
		return
	}

	body := req.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err = p.Save()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, req, "/view/"+title, http.StatusFound)
}

// Main
func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
}
