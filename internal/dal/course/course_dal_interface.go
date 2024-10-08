package dal

import (
	"github.com/doanlng/Go-Api-Tech-Challenge/internal/model"
)

type CourseDao interface {
	Create(course *model.Course) (*model.Course, error)
	List() ([]*model.Course, error)
	Get(id int64) (*model.Course, error)
	Update(c *model.Course, id int64) (*model.Course, error)
	Delete(id int64) (int64, error)
}
