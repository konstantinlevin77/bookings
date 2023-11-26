package main

import (
	"encoding/gob"
	"github.com/konstantinlevin77/bookings/helpers"
	"github.com/konstantinlevin77/bookings/internal/config"
	"github.com/konstantinlevin77/bookings/internal/handlers"
	"github.com/konstantinlevin77/bookings/internal/models"
	"github.com/konstantinlevin77/bookings/internal/render"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

const PORT = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLogger *log.Logger
var errorLogger *log.Logger

func main() {

	err := run()
	if err != nil {
		log.Fatal(err)
	}

	//http.HandleFunc("/", handlers.Repo.HomeHandler)
	//http.HandleFunc("/about",handlers.Repo.AboutHandler)

	log.Println("Listening and serving on ", PORT)
	server := &http.Server{
		Addr:    PORT,
		Handler: routes(&app),
	}

	err = server.ListenAndServe()
	log.Fatal(err.Error())

	//http.ListenAndServe(PORT,nil)

}

func run() error {

	gob.Register(models.Reservation{})

	// set it to true when production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	infoLogger = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLogger = infoLogger

	errorLogger = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLogger = errorLogger

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Println("Error occurred while creating template cache, aborting.")
		log.Fatal(err.Error())
		return err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepository(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)
	helpers.NewHelpers(&app)

	return nil
}
