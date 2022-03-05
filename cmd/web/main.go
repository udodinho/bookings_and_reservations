package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/udodinho/bookings/helpers"
	"github.com/udodinho/bookings/internal/config"
	"github.com/udodinho/bookings/internal/driver"
	"github.com/udodinho/bookings/internal/handlers"
	"github.com/udodinho/bookings/internal/models"
	"github.com/udodinho/bookings/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = ":5000"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main application function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

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

func run() (*driver.DB, error) {
	// If we crash the go code,we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Llongfile)

	//What am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// Change this to true when in production
	app.InProduction = false
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//Connect to the database
	log.Println("Connecting to the database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 user=mac dbname=bookings")
	if err != nil {
		log.Fatal("Can't connect to the database! dying...: ")
	}
	log.Println("Connected to the database!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot load template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepository(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
