package helpers

import (
	"fmt"
	"github.com/udodinho/bookings/internal/config"
	"net/http"
	"runtime/debug"
)

var app *config.AppConfig

// NewHelpers helps set up the app config
func NewHelpers(a *config.AppConfig) {
	app = a
}

// ClientError is a custom error type for client errors
func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of:", status)
	http.Error(w, http.StatusText(status), status)

}

//ServerError is a custom error type for server errors
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
