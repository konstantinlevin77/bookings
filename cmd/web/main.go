package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/konstantinlevin77/bookings/pkg/config"
	"github.com/konstantinlevin77/bookings/pkg/handlers"
	"github.com/konstantinlevin77/bookings/pkg/render"
)


const PORT = ":8080"
var app config.AppConfig
var session *scs.SessionManager


func main() {

	// set it to true when production
	app.InProduction = false	

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	
	if err != nil {
		log.Println("Error occured while creating template cache, aborting.")
		log.Fatal(err.Error())
	}

	app.TemplateCache= tc
	app.UseCache = false

	repo := handlers.NewRepository(&app)
	handlers.NewHandlers(repo)


	render.NewTemplates(&app)

	
	//http.HandleFunc("/", handlers.Repo.HomeHandler)
	//http.HandleFunc("/about",handlers.Repo.AboutHandler)


	log.Println("Listening and serving on ",PORT)
	server := &http.Server{
		Addr: PORT,
		Handler: routes(&app),
	}

	err = server.ListenAndServe()
	log.Fatal(err.Error())


	//http.ListenAndServe(PORT,nil)
	
}