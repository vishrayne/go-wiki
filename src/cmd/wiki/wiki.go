package main

import (
	"html/template"
	"net/http"
	"regexp"
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

// validating path + using fn literal to remove repetitive code
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		m := validPath.FindStringSubmatch(req.URL.Path)
		if m == nil {
			http.NotFound(rw, req)
			return
		}

		fn(rw, req, m[2])
	}
}

// Request handlers
func viewHandler(rw http.ResponseWriter, req *http.Request, title string) {
	page, err := LoadPage(title)
	if err != nil {
		http.Redirect(rw, req, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(rw, "view", page)
}

func editHandler(rw http.ResponseWriter, req *http.Request, title string) {
	page, err := LoadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}

	renderTemplate(rw, "edit", page)
}

func saveHandler(rw http.ResponseWriter, req *http.Request, title string) {
	body := req.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.Save()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, req, "/view/"+title, http.StatusFound)
}

// Main
func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8080", nil)
}
