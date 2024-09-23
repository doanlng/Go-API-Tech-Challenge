package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type courseController struct{}

func (cc courseController) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", cc.List)
	r.Post("/", cc.Create)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("", cc.Get)
		r.Put("", cc.Update)
		r.Delete("", cc.Delete)
	})

	return r
}

func (cc courseController) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hit api/course/ list"))
}

func (cc courseController) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hit api/course/ create"))
}

func (cc courseController) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hit api/course/{id} get"))
}

func (cc courseController) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hit api/course/{id} list"))
}

func (cc courseController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hit api/course/{id} list"))
}
