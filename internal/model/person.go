package model

type Person struct {
	ID        int64   `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Type      string  `json:"type"`
	Age       int     `json:"age"`
	Courses   []int64 `json:"Courses"`
}
