package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/udodinho/bookings/pkg/config"
	"github.com/udodinho/bookings/pkg/handlers"
	"github.com/udodinho/bookings/pkg/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":5000"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {
	// If we crash the go code,we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Llongfile)

	// Change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot load template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepository(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)

	fmt.Printf("Starting server at port: %v\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
