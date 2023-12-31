package handlers

import (
	"encoding/json"
	"github.com/konstantinlevin77/bookings/helpers"
	"github.com/konstantinlevin77/bookings/internal/config"
	"github.com/konstantinlevin77/bookings/internal/driver"
	"github.com/konstantinlevin77/bookings/internal/forms"
	"github.com/konstantinlevin77/bookings/internal/models"
	"github.com/konstantinlevin77/bookings/internal/render"
	"github.com/konstantinlevin77/bookings/internal/repository"
	"github.com/konstantinlevin77/bookings/internal/repository/dbrepo"
	"log"
	"net/http"
	"strconv"
	"time"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepository(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewHandlers(r *Repository) {

	Repo = r
}

func (m *Repository) HomeHandler(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "home.page.html", &models.TemplateData{})

}

func (m *Repository) AboutHandler(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "about.page.html", &models.TemplateData{})

}

func (m *Repository) GeneralsHandler(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "generals.page.html", &models.TemplateData{})

}

func (m *Repository) MajorsHandler(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "majors.page.html", &models.TemplateData{})
}

func (m *Repository) SearchAvailabilityHandler(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "search-availability.page.html", &models.TemplateData{})

}

func (m *Repository) ContactHandler(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "contact.page.html", &models.TemplateData{})

}

func (m *Repository) ReservationHandler(w http.ResponseWriter, r *http.Request) {

	var emptyReservation models.Reservation
	data := make(map[string]interface{})

	data["reservation"] = emptyReservation

	render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
		Form: forms.NewForm(nil),
		Data: data,
	})

}

func (m *Repository) PostReservationHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	form := forms.NewForm(r.PostForm)

	//form.Has("first_name", r)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
	}

	roomRestriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(roomRestriction)
	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

func (m *Repository) PostSearchAvailabilityHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("A post request has been sent."))

}

type jsonResponse struct {
	OK      bool   `json:ok`
	Message string `json:message`
}

func (m *Repository) SearchAvailabilityJSONHandler(w http.ResponseWriter, r *http.Request) {

	resp := jsonResponse{
		OK:      true,
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

func (m *Repository) ReservationSummaryHandler(w http.ResponseWriter, r *http.Request) {

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		log.Println("Couldn't get reservation from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data: data,
	})

}
