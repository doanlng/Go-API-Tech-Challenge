package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Mount("/api/course", courseController{}.Routes())

	http.ListenAndServe(":8000", r)
}
