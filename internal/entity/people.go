package entity

import (
	"github.com/jinzhu/gorm"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/pkg/rnd"
	"github.com/ulule/deepcopier"
	"sync"
	"time"
)

var peopleMutex = sync.Mutex{}

type People struct {
	ID             uint       `gorm:"primary_key" json:"ID" yaml:"-"`
	PeopleUID      string     `gorm:"type:VARBINARY(42);unique_index;" json:"UID" yaml:"UID"`
	PeopleFullName string     `gorm:"type:VARBINARY(755);unique_index;" json:"FullName" yaml:"FullName"`
	PeopleUserId   *string    `gorm:"type:VARCHAR(42);" json:"UserId" yaml:"UserId,omitempty"`
	PeopleBoD      *time.Time `json:"BoD" yaml:"BoD,omitempty"`
	PeopleDeadDate *time.Time `json:"DeadDate" yaml:"DeadDate,omitempty"`
	PhotoCount     int        `gorm:"default:1" json:"PhotoCount" yaml:"-"`
	PlaceCount     int        `gorm:"default:1" json:"PlaceCount" yaml:"-"`
	CreatedAt      time.Time  `json:"CreatedAt" yaml:"-"`
	UpdatedAt      time.Time  `json:"UpdatedAt" yaml:"-"`
	DeletedAt      *time.Time `sql:"index" json:"DeletedAt" yaml:"DeletedAt,omitempty"`
	//	PhotoPeopleTag []PhotoPeopleTag
}

// UnknownPeople is PhotoPrism's default people.
var UnknownPeople = People{
	PeopleUID:      "zz",
	PeopleFullName: "Unknown",
	PeopleUserId:   nil,
	PeopleBoD:      nil,
	PeopleDeadDate: nil,
	PhotoCount:     0,
	PlaceCount:     0,
	CreatedAt:      time.Time{},
	UpdatedAt:      time.Time{},
	DeletedAt:      nil,
}

// CreateUnknownPeople creates the default people if not exists.
func CreateUnknownPeople() {
	FirstOrCreatePeople(&UnknownPeople)
}

// FindPlace finds a matching place or returns nil.
func FindPeople(uid string, fullname string) *People {
	people := People{}

	if fullname == "" {
		if err := Db().Where("people_uid = ?", uid).First(&people).Error; err != nil {
			log.Debugf("peoples: failed finding %s", uid)
			return nil
		} else {
			return &people
		}
	}

	if err := Db().Where("people_uid = ? OR people_full_name = ?", uid, fullname).First(&people).Error; err != nil {
		return nil
	} else {
		return &people
	}
}

// Find fetches entity values from the database the primary key.
func (m *People) Find() error {
	if err := Db().First(m, "id = ?", m.ID).Error; err != nil {
		return err
	}

	return nil
}

// BeforeCreate creates a random UID if needed before inserting a new row to the database.
func (m *People) BeforeCreate(scope *gorm.Scope) error {
	if rnd.IsUID(m.PeopleUID, 'p') {
		return nil
	}

	return scope.SetColumn("people_uid", rnd.PPID('p'))
}

// Create inserts a new row to the database.
func (m *People) Create() error {
	peopleMutex.Lock()
	defer peopleMutex.Unlock()

	return Db().Create(m).Error
}

// NewPeople creates a new people;
func NewPeople(peopleFullName string, peopleUserID *string, peopleDoB *time.Time, peopleDeadDate *time.Time) *People {
	now := Timestamp()

	result := &People{
		PeopleFullName: peopleFullName,
		PeopleUserId:   peopleUserID,
		PeopleBoD:      peopleDoB,
		PeopleDeadDate: peopleDeadDate,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	return result
}

// Saves the entity using form data and stores it in the database.
func (m *People) SaveForm(f form.People) error {
	if err := deepcopier.Copy(m).From(f); err != nil {
		return err
	}

	return Db().Save(m).Error
}

// Update sets a new value for a database column.
func (m *People) Update(attr string, value interface{}) error {
	return UnscopedDb().Model(m).UpdateColumn(attr, value).Error
}

// FirstOrCreatePeople fetches an existing row, inserts a new row or nil in case of errors.
func FirstOrCreatePeople(m *People) *People {
	if m.PeopleUID == "" {
		log.Errorf("places: people must not be empty (find or create)")
		return nil
	}

	if m.PeopleFullName == "" {
		log.Errorf("peoples: fullanme must not be empty (find or create people %s)", m.PeopleUID)
		return nil
	}

	result := People{}

	if findErr := Db().Where("id = ? OR people_full_name = ?", m.ID, m.PeopleFullName).First(&result).Error; findErr == nil {
		return &result
	} else if createErr := m.Create(); createErr == nil {
		return m
	} else if err := Db().Where("id = ? OR people_full_name = ?", m.ID, m.PeopleFullName).First(&result).Error; err == nil {
		return &result
	} else {
		log.Errorf("peoples: %s (create peoples %s)", createErr, m.PeopleUID)
	}

	return nil
}

// Unknown returns true if this is an unknown people
func (m People) Unknown() bool {
	return m.PeopleUID == "" || m.PeopleUID == UnknownPeople.PeopleUID
}

// FullName returns people  FullName
func (m People) FullName() string {
	return m.FullName()
}

// DoD returns people Birth Date
func (m People) DoD() *time.Time {
	return m.PeopleBoD
}

// DeadDate returns people Dead Date
func (m People) DeadDate() *time.Time {
	return m.PeopleDeadDate
}
