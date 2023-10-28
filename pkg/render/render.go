package render

import (
	"bytes"
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

func AddDefaultDate(td *models.TemplateData) *models.TemplateData{

	// if something needs to be default, we can add it here.
	// for now, there's no default data, so it just returns
	// as it is.
	return td
}

func RenderTemplate(w http.ResponseWriter, t string, td *models.TemplateData) {

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

	td = AddDefaultDate(td)

	err := tmpl.Execute(buf,td)

	if err != nil {
		log.Fatal("Error occured while writing to buffer, aborting.")
	}

	// The first return value is number of bytes written. 
	_, err = buf.WriteTo(w)

	if err != nil {
		log.Fatal("Error occured while rendering, aborting.")
	}
}


func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := make(map[string]*template.Template)

	pages,err := filepath.Glob("./templates/*.page.html")

	if err!=nil{
		return myCache,err
	}

	for _,page := range pages {

		name := filepath.Base(page)
		ts,err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache,err
		}

		layouts,err := filepath.Glob("./templates/*.layout.html")

		if err != nil {
			return myCache,err
		}

		if len(layouts) > 0 {
			ts,err = ts.ParseGlob("./templates/*.layout.html")
			
			if err != nil {
				return myCache,err
			}	
		}

		myCache[name] = ts

	}
	return myCache,nil
}