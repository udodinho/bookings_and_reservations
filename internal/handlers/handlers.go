package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/udodinho/bookings/helpers"
	"github.com/udodinho/bookings/internal/config"
	"github.com/udodinho/bookings/internal/driver"
	"github.com/udodinho/bookings/internal/form"
	"github.com/udodinho/bookings/internal/models"
	"github.com/udodinho/bookings/internal/render"
	"github.com/udodinho/bookings/repository"
	"github.com/udodinho/bookings/repository/dbrepo"
	"net/http"
	"strconv"
	"time"
)

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Repo the repository used by the handlers
var Repo *Repository

// NewRepository creates a new Repository
func NewRepository(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresDbRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home function
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.gohtml", &models.TemplateData{})
}

// About is the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.gohtml", &models.TemplateData{})
}

// Contact is the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.gohtml", &models.TemplateData{})
}

// FirstClass is the first class page
func (m *Repository) FirstClass(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "first-class.page.gohtml", &models.TemplateData{})
}

// BusinessClass is the business page
func (m *Repository) BusinessClass(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "business-class.page.gohtml", &models.TemplateData{})
}

// Reservation renders the make reservation page and handles the form submission
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from session"))
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := map[string]string{}
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.gohtml", &models.TemplateData{
		Form:      form.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation handles the form submission for posting.
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from session"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	//sd := r.Form.Get("start_date")
	//ed := r.Form.Get("end_date")

	// 01-03-2022

	//layout := "2006-01-02"
	//startDate, err := time.Parse(layout, sd)
	//if err != nil {
	//	helpers.ServerError(w, err)
	//	return
	//}
	//
	//endDate, err := time.Parse(layout, ed)
	//if err != nil {
	//	helpers.ServerError(w, err)
	//	return
	//}
	//
	//roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	//if err != nil {
	//	helpers.ServerError(w, err)
	//	return
	//}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	//reservation := models.Reservation{
	//	FirstName: r.Form.Get("first_name"),
	//	LastName:  r.Form.Get("last_name"),
	//	Email:     r.Form.Get("email"),
	//	Phone:     r.Form.Get("phone"),
	//	StartDate: startDate,
	//	EndDate:   endDate,
	//	RoomID:    roomID,
	//}

	forms := form.New(r.PostForm)
	//forms.Has("first_name", r)
	forms.Required("first_name", "last_name", "email")
	forms.MinLength("first_name", 3)
	forms.IsEmail("email")

	if !forms.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.gohtml", &models.TemplateData{
			Form: forms,
			Data: data,
		})
		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{
		ReservationID: newReservationID,
		RoomID:        reservation.RoomID,
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservations-summary", http.StatusSeeOther)

}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get error from session")
		m.App.Session.Put(r.Context(), "error", "Cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := map[string]string{}
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.gohtml", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})

}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.gohtml", &models.TemplateData{})
}

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	// 01/03/2022

	layout := "02/01/2006"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No Availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.gohtml", &models.TemplateData{
		Data: data,
	})

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
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/reservations", http.StatusSeeOther)
}
