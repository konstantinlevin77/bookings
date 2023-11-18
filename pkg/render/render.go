package render

import (
	"bytes"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/konstantinlevin77/bookings/pkg/config"
	"github.com/konstantinlevin77/bookings/pkg/models"
)

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {

	td.CSRFToken = nosurf.Token(r)

	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, t string, td *models.TemplateData) {

	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	tmpl, isInMap := tc[t]

	if !isInMap {
		log.Fatal("Template couldn't be found in the cache, aborting.")
	}

	buf := new(bytes.Buffer)

	// call it here

	td = AddDefaultData(td, r)

	err := tmpl.Execute(buf, td)

	if err != nil {
		log.Println(err)
		log.Fatal("Error occurred while writing to buffer, aborting.")

	}

	// The first return value is number of bytes written.
	_, err = buf.WriteTo(w)

	if err != nil {
		log.Fatal("Error occurred while rendering, aborting.")
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := make(map[string]*template.Template)

	pages, err := filepath.Glob("./templates/*.page.html")

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {

		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		layouts, err := filepath.Glob("./templates/*.layout.html")

		if err != nil {
			return myCache, err
		}

		if len(layouts) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")

			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts

	}
	return myCache, nil
}
