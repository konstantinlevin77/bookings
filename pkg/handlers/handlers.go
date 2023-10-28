package handlers

import (
	"net/http"

	"github.com/konstantinlevin77/bookings/pkg/config"
	"github.com/konstantinlevin77/bookings/pkg/models"
	"github.com/konstantinlevin77/bookings/pkg/render"
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
	m.App.Session.Put(r.Context(),"remote_ip",remoteIP)
	render.RenderTemplate(w, "home.page.html",&models.TemplateData{})

}


func (m *Repository) AboutHandler(w http.ResponseWriter, r *http.Request) { 

	//perform some logic

	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	remoteIP := m.App.Session.GetString(r.Context(),"remote_ip")

	stringMap["remote_ip"] = remoteIP


	// send data to the template.

	render.RenderTemplate(w,"about.page.html",&models.TemplateData{StringMap: stringMap})

}