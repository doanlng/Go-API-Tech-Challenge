package dal

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"example.com/model"
	"github.com/lib/pq"
)

type PersonDaoImpl struct {
	DB *sql.DB
}

func NewPersonDAO(db *sql.DB) PersonDao {
	return &PersonDaoImpl{DB: db}
}

func (db *PersonDaoImpl) Get(name *string, age *int64) ([]*model.Person, error) {
	s := `
		SELECT p.id, p.first_name, p.last_name, p.type, p.age,
			COALESCE(array_agg(c.id), '{}') AS courses
		FROM person p
		LEFT JOIN person_course pc ON p.id = pc.person_id
		LEFT JOIN course c ON pc.course_id = c.id
		WHERE 1=1
		`
	var qArgs []interface{}

	if name != nil {
		fn, ln := nameHelper(*name)
		if fn == "" && ln == "" {
			return nil, errors.New("invalid name string passed in")
		}
		qArgs = append(qArgs, fn)
		qArgs = append(qArgs, ln)
		s += " AND first_name ILIKE $1 AND last_name ILIKE $2"
	}

	if age != nil {
		qArgs = append(qArgs, age)
		appendStr := fmt.Sprintf(" AND age = $%d", len(qArgs))
		s += appendStr
	}

	s += " GROUP BY p.ID"
	rows, err := db.DB.Query(s, qArgs...)
	if err != nil {
		log.Panic(err)
		return nil, errors.New("could not locate person(s)")
	}
	defer rows.Close()

	p := []*model.Person{}
	for rows.Next() {
		t := &model.Person{}
		err := rows.Scan(&t.ID, &t.FirstName, &t.LastName, &t.Type, &t.Age, pq.Array(&t.Courses))
		if err != nil {
			log.Panic(err)
			return nil, errors.New("could not scan person(s)")
		}
		p = append(p, t)
	}
	return p, nil
}

func (db *PersonDaoImpl) Create(Person *model.Person) (int64, error) {

	if Person == nil {
		return -1, errors.New("passed in a nil Person")
	}

	const s = `INSERT INTO person (first_name, last_name, type, age) VALUES ($1, $2, $3, $4) RETURNING id`

	stmt, err := db.DB.Prepare(s)
	if err != nil {
		log.Fatal(err)
	}

	var id int64
	err = stmt.QueryRow(Person.FirstName, Person.LastName, Person.Type, Person.Age).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	Person.ID = id

	ci := insertCourses(Person.ID, Person.Courses)
	_, err = db.DB.Exec(ci)
	if err != nil {
		log.Fatal(err)
		return -1, errors.New("issue encountered inserting courses")
	}

	return id, nil
}

func (db *PersonDaoImpl) Update(Person *model.Person, name *string) (*model.Person, error) {
	fn, ln := nameHelper(*name)
	if fn == "" && ln == "" {
		return nil, errors.New("invalid name string passed in")
	}

	const s = "UPDATE person SET type = $1, age = $2 WHERE first_name ILIKE $3 AND last_name ILIKE $4"
	_, err := db.DB.Query(s, Person.Type, Person.Age, fn, ln)
	if err != nil {
		log.Fatal(err)
	}

	const s2 = "SELECT id FROM person WHERE first_name ILIKE $1 AND last_name ILIKE $2"
	row, err := db.DB.Query(s2, fn, ln)
	if err != nil {
		log.Panic(err)
		return nil, errors.New("issue encountered updating courses (deletion)")
	}
	var id int64
	row.Next()
	err = row.Scan(&id)
	if err != nil {
		log.Panic(err)
		return nil, errors.New("panic reading ID")
	}
	const s3 = "DELETE from person_course WHERE person_id = $1"
	_, err = db.DB.Query(s3, id)
	if err != nil {
		log.Panic(err)
		return nil, errors.New("issue encountered updating courses (deletion)")
	}

	insertString := insertCourses(id, Person.Courses)
	_, err = db.DB.Query(insertString)
	if err != nil {
		log.Panic(err)
		return nil, errors.New("issue encountered updating courses (insertion)")
	}

	return Person, nil
}

func (db *PersonDaoImpl) Delete(name string) (int64, error) {
	fn, ln := nameHelper(name)
	if fn == "" && ln == "" {
		return -1, errors.New("error splitting name properly")
	}
	const s = "SELECT id FROM Person WHERE first_name ILIKE $1 AND last_name ILIKE $2"
	st, err := db.DB.Query(s, fn, ln)
	if err != nil {
		log.Fatal(err)
		return -1, errors.New("panic querying for person (delete)")
	}

	var id int64
	st.Next()
	err = st.Scan(&id)
	if err != nil {
		log.Fatal(err)
		return -1, errors.New("panic reading ID")
	}
	log.Println(id)
	const s2 = "DELETE from person_course WHERE person_id = $1"
	_, err = db.DB.Exec(s2, id)
	if err != nil {
		log.Fatal(err)
		return -1, errors.New("issue encountered updating courses (deletion)")
	}

	const s3 = "DELETE from person WHERE id = $1"
	_, err = db.DB.Exec(s3, id)
	if err != nil {
		log.Panic(err)
		return -1, errors.New("issue encountered deleting")
	}

	return id, err
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
