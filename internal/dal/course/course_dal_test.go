package dal

import (
	"log"
	"testing"

	"github.com/doanlng/Go-Api-Tech-Challenge/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setUpTestDbAndDaoCourse() CourseDao {
	// Create in memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failure in creating database for testing")
	}
	err = db.AutoMigrate(&model.Course{})
	if err != nil {
		log.Fatal("failure migrating course schema")
	}

	dao := NewCourseDAO(db)
	return dao
}

func TestListCourse(t *testing.T) {
	// test list function
	tdb := setUpTestDbAndDaoCourse()

	courses, err := tdb.List()

	assert.Nil(t, err)
	assert.Equal(t, len(courses), 0)

	tdb.Create(&model.Course{Name: "TestCourse1"})
	tdb.Create(&model.Course{Name: "TestCourse2"})
	tdb.Create(&model.Course{Name: "TestCourse3"})

	courses, err = tdb.List()
	assert.Nil(t, err)
	assert.Equal(t, len(courses), 3)
}

func TestCreateCourse(t *testing.T) {
	// test course creation
	tdb := setUpTestDbAndDaoCourse()

	course := &model.Course{
		Name: "TestCourse4",
	}

	nc, err := tdb.Create(course)
	assert.Equal(t, nc.ID, int64(1))
	assert.Nil(t, err)

	newCourses, err := tdb.List()
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, len(newCourses), 1)

	course = &model.Course{
		Name: "",
	}

	nc, err = tdb.Create(course)
	assert.Nil(t, nc)
	assert.NotNil(t, err)
}

func TestGetCourse(t *testing.T) {
	// test Course retrieval
	tdb := setUpTestDbAndDaoCourse()
	course := &model.Course{
		Name: "TestCourse1",
	}

	_, err := tdb.Create(course)
	if err != nil {
		log.Fatal(err)
	}
	course, err = tdb.Get(1)

	assert.Equal(t, course.ID, int64(1))
	assert.Equal(t, course.Name, "TestCourse1")
	assert.Nil(t, err)

	course, err = tdb.Get(-1)
	assert.Nil(t, course)
	assert.NotNil(t, err)

	course, err = tdb.Get(2)
	assert.Nil(t, course)
	assert.NotNil(t, err)
}

func TestUpdateCourse(t *testing.T) {
	// test Course updating
	tdb := setUpTestDbAndDaoCourse()
	c := &model.Course{
		Name: "TestCourse1",
	}

	_, err := tdb.Create(c)
	if err != nil {
		log.Fatal(err)
	}

	c = &model.Course{
		Name: "NewTestNameForCourse",
	}
	nc, err := tdb.Update(c, 1)
	assert.Equal(t, nc.ID, int64(1))
	assert.Equal(t, nc.Name, "NewTestNameForCourse")
	assert.Nil(t, err)

	//tests upsert
	nc, err = tdb.Update(c, 100)
	assert.Equal(t, nc.ID, int64(100))
	assert.Equal(t, nc.Name, "NewTestNameForCourse")
	assert.Nil(t, err)

	c = &model.Course{
		Name: "",
	}
	nc, err = tdb.Update(c, 1)
	assert.Nil(t, nc)
	assert.NotNil(t, err)
}

func TestDeleteCourse(t *testing.T) {
	// test course Deletions
	tdb := setUpTestDbAndDaoCourse()
	c := &model.Course{
		Name: "TestCourse1",
	}

	_, err := tdb.Create(c)
	if err != nil {
		log.Fatal(err)
	}

	id, err := tdb.Delete(1)
	assert.Equal(t, id, int64(1))
	assert.Nil(t, err)

	courses, err := tdb.List()
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, len(courses), 0)

	id, err = tdb.Delete(100)
	assert.Equal(t, id, int64(-1))
	assert.NotNil(t, err)
}
