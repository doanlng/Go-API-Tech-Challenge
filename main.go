package main

import (
	"net/http"

	dbconn "example.com/dbconn"

	"example.com/controller"
	"github.com/go-chi/chi/v5"
)

func main() {
	db := dbconn.NewDbConn()
	connection := db.Connect()

	r := chi.NewRouter()
	c := controller.NewCourseController(connection)
	p := controller.NewPersonController(connection)
	r.Mount("/api/course", c.Routes())
	r.Mount("/api/person", p.Routes())
	http.ListenAndServe(":8000", r)
}
