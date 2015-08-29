package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

const viewPath string = "src/cmd/wiki/view/"

// Helpers to save and load page
// TODO: move these to a separate package and write test!
type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

// Request handlers
func viewHandler(rw http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/view/"):]

	page, err := loadPage(title)
	if err != nil {
		http.Redirect(rw, req, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(rw, "view", page)
}

func editHandler(rw http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/edit/"):]

	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}

	renderTemplate(rw, "edit", page)
}

func saveHandler(rw http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/save/"):]
	body := req.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, req, "/view/"+title, http.StatusFound)
}

// generic template renderer
func renderTemplate(rw http.ResponseWriter, view string, p *Page) {
	t, err := template.ParseFiles(viewPath + view + ".html")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(rw, p)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

// Main
func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
}
