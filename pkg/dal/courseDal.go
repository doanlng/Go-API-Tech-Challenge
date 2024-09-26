package dal

import (
	"database/sql"
	"errors"
	"log"

	"example.com/model"
)

type CourseDAO struct {
	DB *sql.DB
}

func NewCourseDAO(db *sql.DB) *CourseDAO {
	return &CourseDAO{DB: db}
}

func (db *CourseDAO) Create(course *model.Course) (*model.Course, error) {

	if course == nil {
		return nil, errors.New("passed in a nil course")
	}

	const s = `INSERT INTO course (name) VALUES ($1) RETURNING id`

	stmt, err := db.DB.Prepare(s)
	if err != nil {
		log.Fatal(err)
	}

	var id int64
	err = stmt.QueryRow(course.Name).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	course.ID = int(id)

	return course, nil
}

func (db *CourseDAO) List() ([]*model.Course, error) {
	const s = "SELECT * FROM course"

	rows, err := db.DB.Query(s)
	if err != nil {
		log.Fatal("issue querying database")
	}
	defer rows.Close()

	c := []*model.Course{}
	for rows.Next() {
		e := &model.Course{}
		err := rows.Scan(&e.ID, &e.Name)
		if err != nil {
			log.Fatal("error scanning rows")
		}
		c = append(c, e)
	}

	return c, nil
}

func (db *CourseDAO) Get(id int64) (*model.Course, error) {
	const s = "SELECT * FROM course WHERE id = $1 LIMIT 1"

	row, err := db.DB.Query(s, id)

	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	c := &model.Course{}

	for row.Next() {
		err := row.Scan(&c.ID, &c.Name)
		if err != nil {
			log.Fatal("error scanning rows")
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	return c, nil
}

func (db *CourseDAO) Update(c *model.Course, id int64) (*model.Course, error) {
	const s = "UPDATE course SET name = $1 where id = $2"
	_, err := db.DB.Exec(s, c.Name, id)
	if err != nil {
		log.Fatal(err)
	}

	var nc *model.Course
	nc, err = db.Get(id)
	if err != nil {
		log.Fatal(err)
	}

	return nc, nil
}

func (db *CourseDAO) Delete(id int64) (int64, error) {
	const s = "DELETE FROM course WHERE id = $1"
	_, err := db.DB.Exec(s, id)
	if err != nil {
		log.Fatal(err)
	}

	return id, err
}
