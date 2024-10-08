package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/doanlng/Go-Api-Tech-Challenge/internal/model"
	"gorm.io/gorm"

	dal "github.com/doanlng/Go-Api-Tech-Challenge/internal/dal/person"
	"github.com/go-chi/chi/v5"
)

type PersonController struct {
	DAO dal.PersonDao
}

func NewPersonController(conn *gorm.DB) *PersonController {
	dao := dal.NewPersonDAO(conn)
	return &PersonController{DAO: dao}
}

func (pc PersonController) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", pc.Get)
	r.Post("/", pc.Create)

	r.Route("/{name}", func(r chi.Router) {
		r.Get("/", pc.Get)
		r.Put("/", pc.Update)
		r.Delete("/", pc.Delete)
	})

	return r
}

func (pc PersonController) Get(w http.ResponseWriter, r *http.Request) {

	ageStr := r.URL.Query().Get("age")
	var ageParam *int64
	if ageStr != "" {
		age, err := strconv.ParseInt(ageStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid URL Parameters passed", http.StatusBadRequest)
			return
		}

		ageParam = &age
	} else {
		ageParam = nil
	}

	var nameParam *string
	var name string
	if name = chi.URLParam(r, "name"); name != "" {
		nameParam = &name
	} else if name = r.URL.Query().Get("name"); name != "" {
		nameParam = &name
	} else {
		nameParam = nil
	}
	log.Println(name)
	people, err := pc.DAO.Get(nameParam, ageParam)
	if err != nil {
		log.Panic(err)
		http.Error(w, "Issue at Controller level list", http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(people)
	if err != nil {
		http.Error(w, "failed to encode courses to JSON", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (pc PersonController) Create(w http.ResponseWriter, r *http.Request) {
	var p *model.Person
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&p); err != nil {
		http.Error(w, "Invalid values passed into create course", http.StatusInternalServerError)
		return
	}

	nc, err := pc.DAO.Create(p)
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

func (pc PersonController) Update(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		http.Error(w, "Person Not Specified", http.StatusBadRequest)
		return
	}
	var p *model.Person
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var np *model.Person
	np, err := pc.DAO.Update(p, &name)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Person Not Found", http.StatusNotFound)
		return
	}

	result, err := json.Marshal(np)
	if err != nil {
		http.Error(w, "failed to marshal course to JSON", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (pc PersonController) Delete(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	idDel, err := pc.DAO.Delete(name)
	if err != nil {
		http.Error(w, "Error Finding Person To Delete", http.StatusNotFound)
		return
	}

	confirm := fmt.Sprintf("Person with ID:%d has been deleted", idDel)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	if err := json.NewEncoder(w).Encode(confirm); err != nil {
		http.Error(w, "failed to encode deleted course ID to JSON", http.StatusInternalServerError)
	}
}
