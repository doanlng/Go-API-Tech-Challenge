package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"example.com/dal"
	"example.com/model"
	"github.com/go-chi/chi/v5"
)

type CourseController struct {
	DAO dal.CourseDao
}

func NewCourseController(conn *sql.DB) *CourseController {
	dao := dal.NewCourseDAO(conn)
	return &CourseController{DAO: dao}
}

func (cc CourseController) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", cc.List)
	r.Post("/", cc.Create)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", cc.Get)
		r.Put("/", cc.Update)
		r.Delete("/", cc.Delete)
	})

	return r
}

func (cc CourseController) List(w http.ResponseWriter, r *http.Request) {
	courses, err := cc.DAO.List()
	if err != nil {
		http.Error(w, "Issue at Controller level list", http.StatusInternalServerError)
	}

	// if err := json.NewEncoder(w).Encode(courses); err != nil {
	// 	http.Error(w, "failed to encode courses to JSON", http.StatusInternalServerError)
	// }

	result, err := json.Marshal(courses)
	if err != nil {
		http.Error(w, "failed to encode courses to JSON", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (cc CourseController) Create(w http.ResponseWriter, r *http.Request) {
	var c *model.Course
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&c); err != nil {
		http.Error(w, "Invalid values passed into create course", http.StatusInternalServerError)
	}
	log.Println(c.Name)
	if c.Name == "" {
		http.Error(w, "Invalid values passed to create course", http.StatusBadRequest)
		return
	}

	nc, err := cc.DAO.Create(c)
	if err != nil {
		http.Error(w, "panic at controller level for create", http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(nc)
	if err != nil {
		http.Error(w, "failed to marshal newly added course to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

func (cc CourseController) Get(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	var c *model.Course
	c, err = cc.DAO.Get(id)
	if err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}
	result, err := json.Marshal(c)
	if err != nil {
		http.Error(w, "failed to marshal course to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (cc CourseController) Update(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	var c *model.Course
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if c.Name == "" {
		http.Error(w, "Invalid name passed for course", http.StatusBadRequest)
		return
	}

	var nc *model.Course
	nc, err = cc.DAO.Update(c, id)
	if err != nil {
		http.Error(w, "Course Not Found", http.StatusNotFound)
	}

	result, err := json.Marshal(nc)
	if err != nil {
		http.Error(w, "failed to marshal course to JSON", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (cc CourseController) Delete(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	idDel, err := cc.DAO.Delete(id)
	if err != nil {
		http.Error(w, "Error Finding Course To Delete", http.StatusNotFound)
		return
	}

	confirm := fmt.Sprintf("Course with ID:%d has been deleted", idDel)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	if err := json.NewEncoder(w).Encode(confirm); err != nil {
		http.Error(w, "failed to encode deleted course ID to JSON", http.StatusInternalServerError)
	}
}
