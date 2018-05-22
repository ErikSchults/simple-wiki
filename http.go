package main

import (
	"html/template"
	"log"
	"net/http"
)

func viewHandler(res http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/view/"):]
	page, err := loadPage(title)

	if err != nil {
		http.Redirect(res, req, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(res, "view", page)
}

func editHandler(res http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/edit/"):]
	page, err := loadPage(title)

	if err != nil {
		page = &Page{Title: title}
	}

	renderTemplate(res, "edit", page)
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

func saveHandler(res http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/save/"):]
	body := req.FormValue("body")

	page := &Page{Title: title, Body: []byte(body)}

	if err := page.save(); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(res, req, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":3300", nil))
}
