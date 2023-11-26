package helpers

import (
	"fmt"
	"github.com/konstantinlevin77/bookings/internal/config"
	"net/http"
	"runtime/debug"
)

var app *config.AppConfig

// NewHelpers populate the app variable.
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {

	app.InfoLogger.Println("Client error with status code:", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {

	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLogger.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}
