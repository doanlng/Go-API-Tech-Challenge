package dal

import (
	"errors"

	"example.com/model"
	"github.com/stretchr/testify/mock"
)

/*
*

	This is a mock that is to be injected into controllers for testing

*
*/
type MockCourseDao struct {
	mock.Mock
}

var _ CourseDao = (*MockCourseDao)(nil)

func NewMockCourseDao() *MockCourseDao {
	return &MockCourseDao{}
}

func (m *MockCourseDao) Create(course *model.Course) (*model.Course, error) {
	args := m.Called(course)                          // record that create method was called with a course object, m.On("create", course) gets its return value put into args
	return args.Get(0).(*model.Course), args.Error(1) // gets the 0th element that is called and cast to model course,
}

func (m *MockCourseDao) List() ([]*model.Course, error) {
	args := m.Called()
	return args.Get(0).([]*model.Course), args.Error(1)
}

func (m *MockCourseDao) Get(id int64) (*model.Course, error) {
	args := m.Called(id)
	course, ok := args.Get(0).(*model.Course)
	if !ok && args.Get(0) != nil {
		// Handle the case where the type is unexpected
		return nil, errors.New("failure to retrieve Course")
	}
	return course, args.Error(1)
}

func (m *MockCourseDao) Update(course *model.Course, id int64) (*model.Course, error) {
	args := m.Called(course, id)
	course, ok := args.Get(0).(*model.Course)
	if !ok && args.Get(0) != nil {
		return nil, errors.New("failure to Update Course")
	}
	return course, args.Error(1)
}

func (m *MockCourseDao) Delete(id int64) (int64, error) {
	args := m.Called(id)
	_, ok := args.Get(0).(*model.Course)
	if !ok && args.Get(0) != nil {
		return -1, errors.New("failure to Delete Course")
	}
	return int64(args.Int(0)), args.Error(1)
}
