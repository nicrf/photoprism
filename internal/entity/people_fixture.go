package entity

import (
	"github.com/photoprism/photoprism/pkg/s2"
	"time"
)

type PeoplesMap map[string]People

func (m PeoplesMap) Get(name string) People {
	if result, ok := m[name]; ok {
		return result
	}

	return UnknownPeople
}

func (m PeoplesMap) Pointer(name string) *People {
	if result, ok := m[name]; ok {
		return &result
	}

	return &UnknownPeople
}

var PeopleFixtures = PeoplesMap{
	"People1": {
		PeopleUID:      s2.TokenPrefix + "85d1ea7d3278",
		PeopleUserId:   new(string),
		PeopleFullName: "Teotihuacán, Mexico, Mexico",
		PeopleBoD:      new(time.Time),
		PeopleDeadDate: nil,
		PhotoCount:     1,
		PlaceCount:     1,
		CreatedAt:      Timestamp(),
		UpdatedAt:      Timestamp(),
	},
}

// CreatePlaceFixtures inserts known entities into the database for testing.
func CreatePeopleFixtures() {
	for _, entity := range PeopleFixtures {
		Db().Create(&entity)
	}
}
