package dal

import (
	"github.com/doanlng/Go-Api-Tech-Challenge/internal/model"
)

type PersonDao interface {
	List() ([]*model.Person, error)
	Create(Person *model.Person) (int64, error)
	Get(name *string, age *int64) (*model.Person, error)
	Update(c *model.Person, name *string) (*model.Person, error)
	Delete(name *string) (int64, error)
}
