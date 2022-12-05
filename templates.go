package main

import (
	"html/template"
)

var tmpl = make(map[string]*template.Template)

func init() {
	m := template.Must
	p := template.ParseFiles
	tmpl["index"] = m(p("templates/index.gohtml", "templates/layout.gohtml"))
	tmpl["about"] = m(p("templates/about.gohtml", "templates/layout.gohtml"))
	tmpl["events"] = m(p("templates/events.gohtml", "templates/layout.gohtml"))
	tmpl["create"] = m(p("templates/create.gohtml", "templates/layout.gohtml"))
	tmpl["post-creation"] = m(p("templates/post-creation.gohtml", "templates/layout.gohtml"))
	tmpl["donate"] = m(p("templates/donate.gohtml", "templates/layout.gohtml"))
}
