package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewPeople(t *testing.T) {

	t.Run("Nicolas Fournier", func(t *testing.T) {
		dob := time.Now()
		people := NewPeople("Nicolas Fournier", nil, &dob, nil)
		assert.Equal(t, "Nicolas Fournier", people.PeopleFullName)
		assert.Equal(t, dob, *people.PeopleBoD)
	})
}
func TestCreatePeople(t *testing.T) {

	t.Run("Nicolas Fournier", func(t *testing.T) {
		dob := time.Now()
		people := NewPeople("Nicolas Fournier", nil, &dob, nil)
		people.Create()

		peopleValidation := FindPeople(people.PeopleUID, "Nicolas Fournier")
		assert.Equal(t, people.ID, peopleValidation.ID)
		assert.Equal(t, people.PeopleUID, peopleValidation.PeopleUID)
		assert.Equal(t, "Nicolas Fournier", peopleValidation.PeopleFullName)
		//TIME ZONE FROM DB assert.Equal(t, dob, *peopleValidation.PeopleBoD)
	})
}
