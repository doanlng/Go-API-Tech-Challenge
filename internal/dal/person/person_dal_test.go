package dal

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setUpTestDbAndDaoPerson() PersonDao {
	// Create in memory database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	CREATE TABLE person
	(
		id         SERIAL PRIMARY KEY,
		first_name TEXT                                          NOT NULL,
		last_name  TEXT                                          NOT NULL,
		type       TEXT CHECK (type IN ('professor', 'student')) NOT NULL,
		age        INTEGER                                       NOT NULL
	);

	CREATE TABLE course
	(
		id   SERIAL PRIMARY KEY,
		name TEXT NOT NULL
	);
	
	INSERT INTO course (name)
	VALUES ('Programming'),
		   ('Databases'),
		   ('UI Design');

	INSERT INTO person (first_name, last_name, type, age)
	VALUES ('Steve', 'Jobs', 'professor', 56),
		   ('Jeff', 'Bezos', 'professor', 60),
		   ('Larry', 'Page', 'student', 51),
		   ('Bill', 'Gates', 'student', 67),
		   ('Elon', 'Musk', 'student', 52);

	CREATE TABLE person_course
	(
		person_id INTEGER NOT NULL,
		course_id INTEGER NOT NULL,
		PRIMARY KEY (person_id, course_id),
		FOREIGN KEY (person_id) REFERENCES person (id),
		FOREIGN KEY (course_id) REFERENCES course (id)
	);
	
	INSERT INTO person_course (person_id, course_id)
	VALUES (1, 1),
			(1, 2),
			(1, 3),
			(2, 1),
			(2, 2),
			(2, 3),
			(3, 1),
			(3, 2),
			(3, 3),
			(4, 1),
			(4, 2),
			(4, 3),
			(5, 1),
			(5, 2),
			(5, 3);
    `)

	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	dao := NewPersonDAO(db)
	return dao
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

func TestGet(t *testing.T) {
	// test Course retrieval
	tdb := setUpTestDbAndDaoPerson()

	p, err := tdb.Get(nil, nil)

	assert.Equal(t, len(p), 5)
	assert.Nil(t, err)

	// course, err = tdb.Get(-1)
	// assert.Nil(t, course)
	// assert.NotNil(t, err)
}

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
