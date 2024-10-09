package dal

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/doanlng/Go-Api-Tech-Challenge/internal/model"

	"gorm.io/gorm"
)

type PersonDaoImpl struct {
	DB *gorm.DB
}

func NewPersonDAO(db *gorm.DB) PersonDao {
	return &PersonDaoImpl{DB: db}
}

func (db *PersonDaoImpl) List() ([]*model.Person, error) {
	var p []*model.Person
	res := db.DB.Find(&p)
	if res.Error != nil || res.RowsAffected == 0 {
		return nil, errors.New(res.Error.Error())
	}

	return p, nil
}

func (db *PersonDaoImpl) Get(name *string, age *int64) (*model.Person, error) {

	var p = &model.Person{}
	q := db.DB.Debug().Preload("Courses")
	if name != nil {
		first, last := nameHelper(*name)
		if first == "" && last == "" {
			return nil, errors.New("unable to parse string name")
		}

		q = q.Where("first_name = ? AND last_name = ?", first, last)

	}

	if age != nil {
		q = q.Where("age = ?", *age)
	}

	res := q.First(&p)
	if res.Error != nil || res.RowsAffected != 1 {
		return nil, errors.New("error retrieving course, could it not exist?")
	}

	return p, nil
}

func (db *PersonDaoImpl) Create(person *model.Person) (int64, error) {
	for _, course := range person.Courses {

		//check if the course exists
		if err := db.DB.First(&course, course.ID).Error; err != nil {
			return -1, errors.New("tried to create a person enrolled in a nonexistent course")
		}

		personCourse := model.PersonCourse{
			PersonID: person.ID,
			CourseID: course.ID,
		}

		res := db.DB.Create(&personCourse)
		if res.Error != nil {
			return -1, errors.New(res.Error.Error())
		}

	}
	res := db.DB.Save(&person)
	if res.Error != nil || res.RowsAffected != 1 {
		return -1, errors.New("error creating a person")
	}
	return person.ID, nil
}

func (db *PersonDaoImpl) Update(person *model.Person, name *string) (*model.Person, error) {

	first, last := nameHelper(*name)
	if first == "" && last == "" {
		return nil, errors.New("invalid name passed in request")
	}
	if f := db.DB.Debug().First(&model.Person{}, "first_name = ? AND last_name = ?", first, last); f.Error != nil {
		return nil, errors.New("couldn't find person, is this a person in the database?")
	}

	var p2 = &model.Person{}
	if err := db.DB.Model(p2).Where("first_name = ? AND last_name = ?", first, last).Updates(person); err.Error != nil && err.RowsAffected != 1 {
		return nil, errors.New("error updating person")
	}

	return person, nil
}

func (db *PersonDaoImpl) Delete(name *string) (int64, error) {

	first, last := nameHelper(*name)
	if first == "" && last == "" {
		return -1, errors.New("invalid name passed in request")
	}

	var p model.Person

	if err := db.DB.Where("first_name = ? AND last_name = ?", first, last).First(&p).Error; err != nil {
		return -1, errors.New("couldn't find person, is this a person in the database?")
	}

	if err := db.DB.Delete(&p); err.Error != nil || err.RowsAffected != 1 {
		return -1, errors.New("error deleting person")
	}

	return p.ID, nil
}

func nameHelper(s string) (string, string) {
	var firstName string
	var lastName string
	var nameArr = splitWhitespace(s)
	if len(nameArr) != 2 {
		log.Panic("could not split name properly")
		return "", ""
	} else {
		firstName = nameArr[0]
		lastName = nameArr[1]
	}

	return firstName, lastName
}

func splitWhitespace(str string) []string {
	ds, err := url.QueryUnescape(str)
	if err != nil {
		fmt.Println("Error decoding URL:", err)
		return nil
	}

	return strings.FieldsFunc(ds, func(r rune) bool {
		return r == ' ' || r == '+'
	})
}
