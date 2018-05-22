package main

import (
	"html/template"
	"log"
	"net/http"
	"regexp"
)

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":3300", nil))
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		match := validPath.FindStringSubmatch(req.URL.Path)
		if match == nil {
			http.NotFound(res, req)
			return
		}
		fn(res, req, match[2])
	}
}

func renderTemplate(res http.ResponseWriter, tmpl string, page *Page) {
	t, err := template.ParseFiles("./templates/" + tmpl + ".html")

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = t.Execute(res, page); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(res http.ResponseWriter, req *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		http.Redirect(res, req, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(res, "view", page)
}

func editHandler(res http.ResponseWriter, req *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}

	renderTemplate(res, "edit", page)
}

func saveHandler(res http.ResponseWriter, req *http.Request, title string) {
	body := req.FormValue("body")
	page := &Page{Title: title, Body: []byte(body)}
	if err := page.save(); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(res, req, "/view/"+title, http.StatusFound)
}
