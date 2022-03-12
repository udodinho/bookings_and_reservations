package handlers

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/udodinho/bookings/internal/config"
	"github.com/udodinho/bookings/internal/models"
	"github.com/udodinho/bookings/internal/render"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplate = "./../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {

	//What am I going to put in the session
	gob.Register(models.Reservation{})

	// Change this to true when in production
	app.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("Cannot load template cache")
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := NewRepository(&app)
	NewHandlers(repo)
	render.NewRenderer(&app)

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	//mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/first-class-quarters", Repo.FirstClass)
	mux.Get("/business-class-suites", Repo.BusinessClass)

	mux.Get("/reservations", Repo.Reservation)
	mux.Post("/reservations", Repo.PostReservation)
	mux.Get("/reservations-summary", Repo.ReservationSummary)

	mux.Get("/contact", Repo.Contact)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// NoSurf adds CSRF protection to all POST requests.
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
