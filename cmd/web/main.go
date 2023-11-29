package main

import (
	"encoding/gob"
	"github.com/konstantinlevin77/bookings/helpers"
	"github.com/konstantinlevin77/bookings/internal/config"
	"github.com/konstantinlevin77/bookings/internal/driver"
	"github.com/konstantinlevin77/bookings/internal/handlers"
	"github.com/konstantinlevin77/bookings/internal/models"
	"github.com/konstantinlevin77/bookings/internal/render"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const PORT = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLogger *log.Logger
var errorLogger *log.Logger

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

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

func run() (*driver.DB, error) {

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

	log.Println("Connecting to the database...")
	db, err := driver.NewDB("host=localhost port=5432 dbname=bookings user=mehmettekman password=")
	if err != nil {
		log.Println(err)
		log.Fatal("Couldn't connect to database, aborting...")

	}
	log.Println("Connected to the database!")

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Println("Error occurred while creating template cache, aborting.")
		log.Fatal(err.Error())
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepository(&app, db)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
