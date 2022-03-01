package helpers

import (
	"github.com/udodinho/bookings/internal/config"
	"net/http"
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
