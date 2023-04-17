package server

import (
	"github.com/go-chi/chi/v5"
	"story-time-server/handler"
)

func UserKidsRoutes(r chi.Router) {
	r.Group(func(kid chi.Router) {
		kid.Post("/", handler.AddKidsProfile)
		//kid.Get("/", handler.GetKidsProfile)
		//kid.Put("/", handler.UpdateKidsProfile)
		//kid.Delete("/", handler.DeleteKidsProfile)
	})
}
