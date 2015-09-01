package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) Save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func LoadPage(title string) (*Page, error) {
	filename := title + ".txt"

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

const absViewPath string = "src/cmd/wiki/view/"

func ViewPath(verb string) string {
	return absViewPath + verb + ".html"
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func GetTitle(rw http.ResponseWriter, req *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(req.URL.Path)
	if m == nil {
		http.NotFound(rw, req)
		return "", errors.New("Title not found")
	}

	return m[2], nil
}
