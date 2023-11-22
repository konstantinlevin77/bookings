package handlers

import (
	"encoding/json"
	"github.com/konstantinlevin77/bookings/internal/config"
	"github.com/konstantinlevin77/bookings/internal/forms"
	"github.com/konstantinlevin77/bookings/internal/models"
	"github.com/konstantinlevin77/bookings/internal/render"
	"log"
	"net/http"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepository(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {

	Repo = r
}

func (m *Repository) HomeHandler(w http.ResponseWriter, r *http.Request) {

	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, r, "home.page.html", &models.TemplateData{})

}

func (m *Repository) AboutHandler(w http.ResponseWriter, r *http.Request) {

	//perform some logic

	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")

	stringMap["remote_ip"] = remoteIP

	// send data to the template.

	render.RenderTemplate(w, r, "about.page.html", &models.TemplateData{StringMap: stringMap})

}

func (m *Repository) GeneralsHandler(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "generals.page.html", &models.TemplateData{})

}

func (m *Repository) MajorsHandler(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "majors.page.html", &models.TemplateData{})
}

func (m *Repository) SearchAvailabilityHandler(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "search-availability.page.html", &models.TemplateData{})

}

func (m *Repository) ContactHandler(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "contact.page.html", &models.TemplateData{})

}

func (m *Repository) ReservationHandler(w http.ResponseWriter, r *http.Request) {

	var emptyReservation models.Reservation
	data := make(map[string]interface{})

	data["reservation"] = emptyReservation

	render.RenderTemplate(w, r, "make-reservation.page.html", &models.TemplateData{
		Form: forms.NewForm(nil),
		Data: data,
	})

}

func (m *Repository) PostReservationHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.NewForm(r.PostForm)

	//form.Has("first_name", r)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, r, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
	}

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
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}
