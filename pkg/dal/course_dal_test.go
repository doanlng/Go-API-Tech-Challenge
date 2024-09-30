package dal

import (
	"database/sql"
	"log"
	"testing"

	"example.com/model"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setUpTestDbAndDao() CourseDao {
	// Create in memory database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	CREATE TABLE course
	(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);
	INSERT INTO course (name)
	VALUES ('TestCourse1'),
		   ('TestCourse2'),
		   ('TestCourse3');
    `)

	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	dao := NewCourseDAO(db)
	return dao
}

func TestList(t *testing.T) {
	// test list function
	tdb := setUpTestDbAndDao()

	courses, err := tdb.List()

	assert.Nil(t, err)
	assert.Equal(t, len(courses), 3)
}

func TestCreate(t *testing.T) {
	// test course creation
	tdb := setUpTestDbAndDao()

	course := &model.Course{
		Name: "TestCourse4",
	}

	nc, err := tdb.Create(course)
	assert.Equal(t, nc.ID, 4)
	assert.Nil(t, err)

	newCourses, err := tdb.List()
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, len(newCourses), 4)
}

func TestGet(t *testing.T) {
	// test Course retrieval
	tdb := setUpTestDbAndDao()

	course, err := tdb.Get(1)

	assert.Equal(t, course.ID, 1)
	assert.Equal(t, course.Name, "TestCourse1")
	assert.Nil(t, err)

	course, err = tdb.Get(-1)
	assert.Nil(t, course)
	assert.NotNil(t, err)
}

func TestUpdate(t *testing.T) {
	// test Course updating
	tdb := setUpTestDbAndDao()

	c := &model.Course{
		Name: "NewTestNameForCourse",
	}
	nc, err := tdb.Update(c, 1)
	assert.Equal(t, nc.ID, 1)
	assert.Equal(t, nc.Name, "NewTestNameForCourse")
	assert.Nil(t, err)

	nc, err = tdb.Update(c, 100)
	assert.Nil(t, nc)
	assert.NotNil(t, err)

}

func TestDelete(t *testing.T) {
	// test course Deletions
	tdb := setUpTestDbAndDao()

	id, err := tdb.Delete(2)
	assert.Equal(t, id, int64(2))
	assert.Nil(t, err)

	courses, err := tdb.List()
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, len(courses), 2)

	id, err = tdb.Delete(100)
	assert.Equal(t, id, int64(-1))
	assert.NotNil(t, err)
}
