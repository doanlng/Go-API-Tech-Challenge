package main

import (
	"net/http"

	dbconn "github.com/doanlng/Go-Api-Tech-Challenge/internal/db_conn"
	"github.com/go-chi/chi"
)

func main() {
	db := dbconn.NewDbConn()
	connection := db.Connect()

	r := chi.NewRouter()
	c := cc.NewCourseController(connection)
	p := pc.NewPersonController(connection)
	r.Mount("/api/course", c.Routes())
	r.Mount("/api/person", p.Routes())
	http.ListenAndServe(":8000", r)

}
