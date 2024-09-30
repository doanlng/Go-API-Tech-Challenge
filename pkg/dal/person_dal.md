package dal

import (
	"database/sql"
	"errors"
	"log"

	"example.com/model"
)

// type PersonDaoImpl struct {
	DB *sql.DB
}

func NewPersonDAO(db *sql.DB) PersonDao {
	return &PersonDaoImpl{DB: db}
}

func (db *PersonDaoImpl) Create(Person *model.Person) (*model.Person, error) {

	if Person == nil {
		return nil, errors.New("passed in a nil Person")
	}

	const s = `INSERT INTO Person (name) VALUES ($1) RETURNING id`

	stmt, err := db.DB.Prepare(s)
	if err != nil {
		log.Fatal(err)
	}

	var id int64
	err = stmt.QueryRow(Person.Name).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	Person.ID = int(id)

	return Person, nil
}

func (db *PersonDaoImpl) List() ([]*model.Person, error) {
	const s = "SELECT * FROM Person"

	rows, err := db.DB.Query(s)
	if err != nil {
		log.Fatal("issue querying database")
	}
	defer rows.Close()

	c := []*model.Person{}
	for rows.Next() {
		e := &model.Person{}
		err := rows.Scan(&e.ID, &e.FirstName, &e.LastName, &e.Type, &e.Age, &e.Courses)
		if err != nil {
			log.Fatal("error scanning rows")
		}
		c = append(c, e)
	}

	return c, nil
}

func (db *PersonDaoImpl) Get(id int64) (*model.Person, error) {
	const s = "SELECT * FROM Person WHERE id = $1 LIMIT 1"

	row, err := db.DB.Query(s, id)

	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	p := &model.Person{}

	for row.Next() {
		err := row.Scan(&p.ID, &p.FirstName, &p.LastName, &p.Type, &p.Age, &p.Courses)
		if err != nil {
			log.Fatal("error scanning rows")
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	return p, nil
}

func (db *PersonDaoImpl) Update(c *model.Person, id int64) (*model.Person, error) {
	const s = "UPDATE Person SET name = $1 where id = $2"
	_, err := db.DB.Exec(s, c.Name, id)
	if err != nil {
		log.Fatal(err)
	}

	var nc *model.Person
	nc, err = db.Get(id)
	if err != nil {
		log.Fatal(err)
	}

	return nc, nil
}

func (db *PersonDaoImpl) Delete(id int64) (int64, error) {
	const s = "DELETE FROM Person WHERE id = $1"
	_, err := db.DB.Exec(s, id)
	if err != nil {
		log.Fatal(err)
	}

	return id, err
}
