package model

type PersonCourse struct {
	PersonID int64 `gorm:"column:person_id;not null" json:"person_id"`
	CourseID int64 `gorm:"column:course_id;not null" json:"course_id"`
}
