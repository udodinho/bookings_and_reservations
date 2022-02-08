package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/udodinho/bookings/pkg/config"
	"github.com/udodinho/bookings/pkg/models"
	"github.com/udodinho/bookings/pkg/render"
	"log"
	"net/http"
)

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// Repo the repository used by the handlers
var Repo *Repository

// NewRepo creates a new Repository
func NewRepository(a *config.AppConfig) *Repository {
	return &Repository{App: a}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home function
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, r, "home.page.gohtml", &models.TemplateData{})
}

// About is the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := map[string]string{}
	stringMap["test"] = "Hello, again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.RenderTemplate(w, r, "about.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})

}

// Contact is the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.gohtml", &models.TemplateData{})
}

// FirstClass is the first class page
func (m *Repository) FirstClass(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "first-class.page.gohtml", &models.TemplateData{})
}

// BusinessClass is the business page
func (m *Repository) BusinessClass(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "business-class.page.gohtml", &models.TemplateData{})
}

// Reservation renders the make reservation page and handles the form submission
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.gohtml", &models.TemplateData{})
}

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("Start date is %s and End date is %s. ", start, end)))
}

type responseJSON struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles requests for the availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := responseJSON{
		OK:      false,
		Message: "Available",
	}

	out, err := json.MarshalIndent(resp, "", "   ")
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}
