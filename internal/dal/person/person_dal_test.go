package dal

import (
	"log"
	"testing"

	dal "github.com/doanlng/Go-Api-Tech-Challenge/internal/dal/course"
	"github.com/doanlng/Go-Api-Tech-Challenge/internal/model"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setUpTestDbAndDaoPerson() (PersonDao, dal.CourseDao) {
	// Create in memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failure in creating database for testing")
	}
	err = db.AutoMigrate(&model.Person{}, &model.Course{}, &model.PersonCourse{})
	if err != nil {
		log.Fatal("failure migrating Person schema")
	}

	dao := NewPersonDAO(db)
	cdao := dal.NewCourseDAO(db)

	c1 := &model.Course{
		Name: "C1",
		ID:   int64(1),
	}

	c2 := &model.Course{
		Name: "C2",
		ID:   int64(2),
	}

	c3 := &model.Course{
		Name: "C3",
		ID:   int64(3),
	}

	cdao.Create(c1)
	cdao.Create(c2)
	cdao.Create(c3)

	person := &model.Person{
		FirstName: "Jon",
		LastName:  "Doe",
		Type:      "student",
		Age:       21,
		Courses:   []model.Course{*c1, *c2, *c3},
	}
	_, err = dao.Create(person)
	if err != nil {
		log.Panic(err)
	}

	person = &model.Person{
		FirstName: "jane",
		LastName:  "Ed",
		Type:      "professor",
		Age:       50,
		Courses:   []model.Course{*c1, *c2},
	}
	_, err = dao.Create(person)
	if err != nil {
		log.Panic(err)
	}
	person = &model.Person{
		FirstName: "Jonny",
		LastName:  "Do",
		Type:      "student",
		Age:       18,
		Courses:   []model.Course{*c2, *c3},
	}
	_, err = dao.Create(person)
	if err != nil {
		log.Panic(err)
	}

	return dao, cdao
}

func TestList(t *testing.T) {
	tdb, _ := setUpTestDbAndDaoPerson()
	p, err := tdb.List()
	if err != nil {
		log.Panic(err)
	}
	assert.Equal(t, len(p), 3)
}

func TestCreate(t *testing.T) {
	// test course creation
	tdb, cdb := setUpTestDbAndDaoPerson()

	nc := &model.Course{
		Name: "New Course",
		ID:   int64(4),
	}

	person := &model.Person{
		FirstName: "Steve",
		LastName:  "Jobs",
		Type:      "professor",
		Age:       30,
		Courses:   []model.Course{*nc},
	}

	np, err := tdb.Create(person)
	assert.NotNil(t, err) //course enrolled in doesn't exist
	assert.Equal(t, int64(-1), np)

	cdb.Create(nc)
	np, err = tdb.Create(person)
	assert.Nil(t, err) //course enrolled in doesn't exist
	assert.Equal(t, np, int64(4))

	p, err := tdb.List()
	if err != nil {
		log.Panic(err)
	}
	assert.Equal(t, len(p), int64(4))

}

func TestGet(t *testing.T) {
	// test Course retrieval
	tdb, _ := setUpTestDbAndDaoPerson()

	name := "Jon Doe"
	p, err := tdb.Get(&name, nil)
	assert.Equal(t, p.ID, int64(1))
	assert.Equal(t, p.FirstName, "Jon")
	assert.Equal(t, p.LastName, "Doe")
	assert.Equal(t, len(p.Courses), 3)
	assert.Nil(t, err)

	var age int64 = 18
	p2, err := tdb.Get(nil, &age)
	assert.Equal(t, p2.ID, int64(3))
	assert.Equal(t, p2.FirstName, "Jonny")

	assert.Nil(t, err)
}

func TestUpdate(t *testing.T) {
	// test Course updating
	tdb, _ := setUpTestDbAndDaoPerson()

	p := &model.Person{
		FirstName: "Gormy",
		LastName:  "Gorm",
	}

	var name = "Jon Doe"
	up, err := tdb.Update(p, &name)
	assert.Nil(t, err)
	assert.Equal(t, up.FirstName, "Gormy")
	assert.Equal(t, up.LastName, "Gorm")

	var nName = "Gormy Gorm"
	np, err := tdb.Get(&nName, nil)
	assert.Nil(t, err)
	assert.Equal(t, np.Age, 21)
	assert.Equal(t, np.ID, int64(1))

	var noName = "name nothere"
	up, err = tdb.Update(p, &noName)
	assert.NotNil(t, err)
	assert.Nil(t, up)
}

func TestDelete(t *testing.T) {
	tdb, _ := setUpTestDbAndDaoPerson()

	var name = "Jon Doe"
	p, err := tdb.Delete(&name)
	assert.Equal(t, p, int64(1))
	assert.Nil(t, err)

	var noName = "name nothere"
	p2, err := tdb.Delete(&noName)
	assert.Equal(t, p2, int64(-1))
	assert.NotNil(t, err)
}
