package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/dal"
	"example.com/model"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {

	courses := []*model.Course{
		{Name: "NewCourse1", ID: 1},
		{Name: "NewCourse2", ID: 2},
		{Name: "NewCourse3", ID: 3},
	}

	mock := dal.NewMockCourseDao()
	mock.On("List").Return(courses, nil)

	tcc := &CourseController{DAO: mock}

	req, err := http.NewRequest("GET", "/course", nil)
	assert.Nil(t, err)

	w := httptest.NewRecorder()

	tcc.List(w, req)

	expected := `[{"id":1,"name":"NewCourse1"},{"id":2,"name":"NewCourse2"},{"id":3,"name":"NewCourse3"}]`
	assert.Equal(t, w.Result().StatusCode, http.StatusOK)
	assert.JSONEq(t, expected, w.Body.String())

	mock.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	course := &model.Course{
		Name: "NewCourse1",
	}

	newCourse := &model.Course{
		ID:   1,
		Name: "NewCourse1",
	}

	mock := dal.NewMockCourseDao()
	mock.On("Create", course).Return(newCourse, nil)

	tcc := &CourseController{DAO: mock}
	j, err := json.Marshal(course)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", "/course", bytes.NewBuffer([]byte(j)))
	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()

	tcc.Create(w, req)

	expected := `{"id":1,"name":"NewCourse1"}`
	assert.Equal(t, w.Result().StatusCode, http.StatusCreated)
	assert.Equal(t, expected, w.Body.String())
	mock.AssertExpectations(t)

	invalidCourse := &model.Course{
		Name: "",
	}
	j, err = json.Marshal(invalidCourse)
	if err != nil {
		panic(err)
	}
	mock.On("Create", invalidCourse).Return(nil, nil)
	req, err = http.NewRequest("POST", "/course", bytes.NewBuffer([]byte(j)))
	if err != nil {
		panic(err)
	}
	w = httptest.NewRecorder()
	tcc.Create(w, req)

	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
	assert.Equal(t, w.Body.String(), "Invalid values passed to create course\n")
	mock.AssertExpectations(t)
}

func TestGet(t *testing.T) {
	course := &model.Course{
		ID:   1,
		Name: "NewCourse1",
	}
	mock := dal.NewMockCourseDao()
	tcc := &CourseController{DAO: mock}
	r := chi.NewRouter()
	r.Route("/course", func(r chi.Router) {
		r.Get("/{id}", tcc.Get)
	})

	id := 1
	mock.On("Get", int64(id)).Return(course, nil)
	req, err := http.NewRequest("GET", "/course/1", nil)
	if err != nil {
		panic(err)
	}
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	mock.AssertExpectations(t)

	id = -1
	mock.On("Get", int64(id)).Return(nil, errors.New("failure to retrieve Course"))
	req, err = http.NewRequest("GET", "/course/-1", nil)
	if err != nil {
		panic(err)
	}
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	mock.AssertExpectations(t)

	req, err = http.NewRequest("GET", "/course/string", nil)
	if err != nil {
		panic(err)
	}
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	mock.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	//setup
	mock := dal.NewMockCourseDao()
	tcc := &CourseController{DAO: mock}
	r := chi.NewRouter()
	r.Route("/course", func(r chi.Router) {
		r.Put("/{id}", tcc.Update)
	})

	//case 1
	newCourse := &model.Course{
		Name: "NewCourse2",
	}
	j, err := json.Marshal(newCourse)
	if err != nil {
		panic(err)
	}
	id := 1
	mock.On("Update", newCourse, int64(id)).Return(&model.Course{ID: 1, Name: "NewCourse2"}, nil)
	req, err := http.NewRequest("PUT", "/course/1", bytes.NewBuffer([]byte(j)))
	if err != nil {
		panic(err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	//mock.AssertExpectations(t)

	//case 2
	newCourse2 := &model.Course{
		Name: "",
	}
	j, err = json.Marshal(newCourse2)
	if err != nil {
		panic(err)
	}
	id = 1
	req, err = http.NewRequest("PUT", "/course/1", bytes.NewBuffer([]byte(j)))
	if err != nil {
		panic(err)
	}
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	//case 3
	newCourse3 := &model.Course{
		Name: "NewCourse3",
	}
	j, err = json.Marshal(newCourse3)
	if err != nil {
		panic(err)
	}
	id = -1
	mock.On("Update", newCourse3, int64(id)).Return(nil, errors.New("failure to Update Course"))
	req, err = http.NewRequest("PUT", "/course/-1", bytes.NewBuffer([]byte(j)))
	if err != nil {
		panic(err)
	}
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	mock.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	//setup
	mock := dal.NewMockCourseDao()
	tcc := &CourseController{DAO: mock}
	r := chi.NewRouter()
	r.Route("/course", func(r chi.Router) {
		r.Delete("/{id}", tcc.Delete)
	})

	//case 1
	courseId := 1
	mock.On("Delete", int64(courseId)).Return(1, nil)
	req, err := http.NewRequest("DELETE", "/course/1", nil)
	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, w.Result().StatusCode, http.StatusNoContent)

	confirm := fmt.Sprintf("\"Course with ID:%d has been deleted\"\n", courseId)
	assert.Equal(t, w.Body.String(), confirm)
	mock.AssertExpectations(t)

	//case 2
	courseId = 100
	mock.On("Delete", int64(courseId)).Return(-1, errors.New("Could not locate course"))
	req, err = http.NewRequest("DELETE", "/course/100", nil)
	if err != nil {
		panic(err)
	}

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, w.Result().StatusCode, http.StatusNotFound)
	mock.AssertExpectations(t)
}
