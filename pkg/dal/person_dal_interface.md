package dal

import (
	"example.com/model"
)

type PersonDao interface {
	Create(Person *model.Person) (*model.Person, error)
	List() ([]*model.Person, error)
	Get(id int64) (*model.Person, error)
	Update(c *model.Person, id int64) (*model.Person, error)
	Delete(id int64) (int64, error)
}
