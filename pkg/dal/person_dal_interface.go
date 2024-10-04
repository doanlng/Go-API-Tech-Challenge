package dal

import (
	"example.com/model"
)

type PersonDao interface {
	Create(Person *model.Person) (int64, error)
	Get(name *string, age *int64) ([]*model.Person, error)
	Update(c *model.Person, name *string) (*model.Person, error)
	Delete(name string) (int64, error)
}
