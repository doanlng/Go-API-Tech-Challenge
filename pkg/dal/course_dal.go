package dal

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"example.com/model"
)

type CourseDaoImpl struct {
	DB *sql.DB
}

func NewCourseDAO(db *sql.DB) CourseDao {
	return &CourseDaoImpl{DB: db}
}

func (db *CourseDaoImpl) Create(course *model.Course) (*model.Course, error) {

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

func (db *CourseDaoImpl) List() ([]*model.Course, error) {
	const s = "SELECT * FROM course"

	rows, err := db.DB.Query(s)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	c := []*model.Course{}
	for rows.Next() {
		e := &model.Course{}
		err := rows.Scan(&e.ID, &e.Name)
		if err != nil {
			log.Fatal(err)
		}
		c = append(c, e)
	}

	return c, nil
}

func (db *CourseDaoImpl) Get(id int64) (*model.Course, error) {
	const s = "SELECT * FROM course WHERE id = $1 LIMIT 1"

	row, err := db.DB.Query(s, id)

	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	c := &model.Course{}

	if r := row.Next(); !r {
		return nil, errors.New("no course could be located")
	} else {
		err := row.Scan(&c.ID, &c.Name)
		if err != nil {
			log.Fatal("error scanning rows")
		}
	}
	return c, nil
}

func (db *CourseDaoImpl) Update(c *model.Course, id int64) (*model.Course, error) {
	const s = "UPDATE course SET name = $1 where id = $2"
	result, err := db.DB.Exec(s, c.Name, id)
	if err != nil {
		log.Fatal(err)
	}

	if res, err := result.RowsAffected(); err == nil && res == 0 {
		s := fmt.Sprintf("Attempted to update, but Course with ID: %d could not be found", id)
		return nil, errors.New(s)
	}

	var nc *model.Course
	nc, err = db.Get(id)
	if err != nil {
		log.Fatal(err)
	}

	return nc, nil
}

func (db *CourseDaoImpl) Delete(id int64) (int64, error) {
	const s = "DELETE FROM course WHERE id = $1"
	res, err := db.DB.Exec(s, id)
	if err != nil {
		log.Fatal(err)
	}

	if row, err := res.RowsAffected(); err == nil && row == 0 {
		s := fmt.Sprintf("Attempted to delete, but Course with ID: %d could not be found", id)
		return -1, errors.New(s)
	}

	return id, err
}
