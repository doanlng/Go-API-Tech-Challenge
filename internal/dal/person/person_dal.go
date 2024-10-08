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
	var fn, ln *string
	if name != nil {
		first, last := nameHelper(*name)
		if first == "" && last == "" {
			return nil, errors.New("unable to parse string name")
		}
		fn = &first
		ln = &last
	}

	var a *int
	if age != nil {
		n := int(*age)
		a = &n
	}

	var p = &model.Person{
		FirstName: *fn,
		LastName:  *ln,
		Age:       *a,
	}
	res := db.DB.Preload("Courses").First(p)
	if res.Error != nil || res.RowsAffected != 1 {
		return nil, errors.New("error retrieving course, could it not exist?")
	}

	return p, nil
}

func (db *PersonDaoImpl) Create(person *model.Person) (int64, error) {
	res := db.DB.Save(&person)
	if res.Error != nil || res.RowsAffected != 1 {
		return -1, errors.New("error creating a person")
	}

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
	return person.ID, nil
}

func (db *PersonDaoImpl) Update(Person *model.Person, name *string) (*model.Person, error) {
	return nil, nil
}

func (db *PersonDaoImpl) Delete(name string) (int64, error) {
	return -1, nil
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

func insertCourses(id int64, cids []int64) string {
	tups := []string{}
	for _, course := range cids {
		tups = append(tups, fmt.Sprintf("(%d, %d)", id, course))
	}
	var vals = strings.Join(tups, ", ")
	var s = `INSERT INTO person_course (person_id, course_id) VALUES ` + vals

	return s
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
