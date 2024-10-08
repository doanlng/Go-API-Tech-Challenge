package dal

import (
	"errors"
	"log"

	"github.com/doanlng/Go-Api-Tech-Challenge/internal/model"
	"gorm.io/gorm"
)

type CourseDaoImpl struct {
	DB *gorm.DB
}

func NewCourseDAO(db *gorm.DB) CourseDao {
	return &CourseDaoImpl{DB: db}
}

func (db *CourseDaoImpl) Create(course *model.Course) (*model.Course, error) {

	if course.Name == "" {
		return nil, errors.New("passed in a course wtih no name")
	}

	res := db.DB.Save(&course)
	if res.Error != nil {
		log.Panic("Error Creating course")
		return nil, errors.New("Unable to create courses")
	}
	return course, nil
}

func (db *CourseDaoImpl) List() ([]*model.Course, error) {
	var courses []*model.Course
	cs := db.DB.Find(&courses)

	if cs.Error != nil {
		return nil, errors.New("unable to pull courses")
	}

	return courses, nil
}

func (db *CourseDaoImpl) Get(id int64) (*model.Course, error) {
	if id <= 0 {
		return nil, errors.New("negative ID passed")
	}
	var c = &model.Course{
		ID: id,
	}
	res := db.DB.First(&c)
	if res.Error != nil {
		return nil, errors.New("unable to get course with ID passed")
	}
	return c, nil
}

// performs updates given posted ID, if posted ID doesn't exist, will upsert
func (db *CourseDaoImpl) Update(c *model.Course, id int64) (*model.Course, error) {
	var nc = &model.Course{
		ID:   id,
		Name: c.Name,
	}

	if c.Name == "" {
		return nil, errors.New("passed in an invalid name")
	}

	res := db.DB.Save(&nc)
	if res.Error != nil {
		return nil, errors.New("course Update Failed")
	}

	return nc, nil
}

func (db *CourseDaoImpl) Delete(id int64) (int64, error) {
	var c = &model.Course{
		ID: id,
	}
	res := db.DB.Delete(&c)
	if res.Error != nil {
		return -1, errors.New("course Delete Failed")
	} else if res.RowsAffected == 0 {
		return -1, errors.New("no Course was affected delete")
	}
	return c.ID, nil
}
