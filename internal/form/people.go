package form

import (
	"github.com/ulule/deepcopier"
	"time"
)

// People represents an people edit form.
type People struct {
	PeopleFullName string     `json:"FullName"`
	PeopleUserId   *string    `json:"UserI"`
	PeopleBoD      *time.Time `json:"BoD "`
	PeopleDeadDate *time.Time `json:"DeadDate "`
}

func NewPeople(m interface{}) (f People, err error) {
	err = deepcopier.Copy(m).To(&f)

	return f, err
}
