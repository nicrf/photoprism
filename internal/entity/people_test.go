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
