package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// WriteConsole is a middleware that prints out to stdout every time the page is hit.
func WriteConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Page is hit.")

		// A crucial part it seems, I didn't get it yet.
		next.ServeHTTP(w,r)
	})
}


func NoSurf(next http.Handler) http.Handler {

	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}