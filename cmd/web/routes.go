package main

import (
	"github.com/konstantinlevin77/bookings/internal/config"
	"github.com/konstantinlevin77/bookings/internal/handlers"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.HomeHandler)
	mux.Get("/about", handlers.Repo.AboutHandler)
	mux.Get("/generals-quarters", handlers.Repo.GeneralsHandler)
	mux.Get("/majors-suite", handlers.Repo.MajorsHandler)
	mux.Get("/search-availability", handlers.Repo.SearchAvailabilityHandler)
	mux.Get("/contact", handlers.Repo.ContactHandler)

	mux.Get("/make-reservation", handlers.Repo.ReservationHandler)
	mux.Post("/make-reservation", handlers.Repo.PostReservationHandler)

	mux.Post("/search-availability", handlers.Repo.PostSearchAvailabilityHandler)

	mux.Post("/search-availability-json", handlers.Repo.SearchAvailabilityJSONHandler)

	// File server returns a http handler which serves the content of the specified root.

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
