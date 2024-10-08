package dal

import (
	"log"
	"testing"

	"github.com/doanlng/Go-Api-Tech-Challenge/internal/model"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setUpTestDbAndDaoPerson() PersonDao {
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

	_, err = dao.List()
	if err != nil {
		log.Panic(err)
	}

	return dao
}

func TestList(t *testing.T) {
	tdb := setUpTestDbAndDaoPerson()
	p, err := tdb.List()
	if err != nil {
		log.Panic(err)
	}
	assert.Equal(t, len(p), 3)
}

// func TestCreate(t *testing.T) {
// 	// test course creation
// 	tdb := setUpTestDbAndDaoPerson()

// 	person := &model.Person{
// 		FirstName: "Jon",
// 		LastName:  "Doe",
// 		Type:      "student",
// 		Age:       1000,
// 		Courses:   []int64{1, 2, 3},
// 	}

// 	np, err := tdb.Create(person)
// 	log.Println(np)
// 	assert.Nil(t, err)
// }

// func TestGet(t *testing.T) {
// 	// test Course retrieval
// 	tdb := setUpTestDbAndDaoPerson()

// 	p, err := tdb.Get(nil, nil)

// 	assert.Equal(t, len(p), 5)
// 	assert.Nil(t, err)

// 	// course, err = tdb.Get(-1)
// 	// assert.Nil(t, course)
// 	// assert.NotNil(t, err)
// }

// func TestUpdate(t *testing.T) {
// 	// test Course updating
// 	tdb := setUpTestDbAndDaoPerson()

// 	c := &model.Course{
// 		Name: "NewTestNameForCourse",
// 	}
// 	nc, err := tdb.Update(c, 1)
// 	assert.Equal(t, nc.ID, 1)
// 	assert.Equal(t, nc.Name, "NewTestNameForCourse")
// 	assert.Nil(t, err)

// 	nc, err = tdb.Update(c, 100)
// 	assert.Nil(t, nc)
// 	assert.NotNil(t, err)

// }

// func TestDelete(t *testing.T) {
// 	// test course Deletions
// 	tdb := setUpTestDbAndDaoPerson()

// 	id, err := tdb.Delete(2)
// 	assert.Equal(t, id, int64(2))
// 	assert.Nil(t, err)

// 	courses, err := tdb.List()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	assert.Equal(t, len(courses), 2)

// 	id, err = tdb.Delete(100)
// 	assert.Equal(t, id, int64(-1))
// 	assert.NotNil(t, err)
// }
