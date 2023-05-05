package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"snippetbox.jaswanthp.com/internal/models"
	"snippetbox.jaswanthp.com/ui"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{

	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		// base retrives the file name (everything after the last "/" i.e home.tmpl.html)
		name := filepath.Base(page)
		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
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
