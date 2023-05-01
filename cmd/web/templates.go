package main

import (
	"html/template"
	"path/filepath"
	"time"

	"snippetbox.jaswanthp.com/internal/models"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
	Flash       string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:05")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		// base retrives the file name (everything after the last "/" i.e home.tmpl.html)
		name := filepath.Base(page)
		// add the parsed base template to template set
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}
		// adding all the parsed partials to the template set
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}
		// adding the main parsed page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil

}

// func myMiddleware(next http.Handler) http.Handler {
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		// middleware logic
// 		next.ServeHTTP(w, r)
// 	}
// 	return http.HandlerFunc(fn)
// }

// func myShortMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// middlware logic
// 		next.ServeHTTP(w, r)
// 	})
// }
