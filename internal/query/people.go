package query

import (
	"github.com/photoprism/photoprism/internal/entity"
	"time"
)

// AlbumResult contains found albums
type PeopleResult struct {
	ID             uint       `json:"-"`
	PeopleUID      string     `json:"UID"`
	PeopleFullName string     `json:"PeopleFullName"`
	PeopleUserId   *string    `json:"PeopleUserId"`
	PeopleBoD      *time.Time `json:"PeopleBoD"`
	PeopleDeadDate *time.Time `json:"PeopleDeadDate"`
	PhotoCount     int        `json:"PhotoCount"`
	PlaceCount     int        `json:"PlaceCount"`
	CreatedAt      time.Time  `json:"CreatedAt"`
	UpdatedAt      time.Time  `json:"UpdatedAt"`
	DeletedAt      *time.Time `json:"DeletedAt,omitempty"`
}

type PeopleResults []PeopleResult

// PeopleByID returns a People based on the ID.
func PeopleByUID(peopleUID string) (people entity.People, err error) {
	if err := Db().Where("people_uid = ?", peopleUID).First(&people).Error; err != nil {
		return people, err
	}

	return people, nil
}

// PeopleCoverByUID returns a people preview file based on the uid.
func PeopleCoverByUID(peopleUID string) (file entity.File, err error) {
	a := entity.People{}

	if err := Db().Where("people_uid = ?", peopleUID).First(&a).Error; err != nil {
		return file, err
	} /* //TODO
		else if a.AlbumType != entity.AlbumDefault { // TODO: Optimize
		f := form.PhotoSearch{Album: a.AlbumUID, Filter: a.AlbumFilter, Order: entity.SortOrderRelevance, Count: 1, Offset: 0, Merged: false}

		if photos, _, err := PhotoSearch(f); err != nil {
			return file, err
		} else if len(photos) > 0 {
			for _, photo := range photos {
				if err := Db().Where("photo_uid = ? AND file_primary = 1", photo.PhotoUID).First(&file).Error; err != nil {
					return file, err
				} else {
					return file, nil
				}
			}
		}

		return file, fmt.Errorf("found no cover for moment")
	}

	if err := Db().Where("files.file_primary = 1 AND files.file_missing = 0 AND files.file_type = 'jpg' AND files.deleted_at IS NULL").
		Joins("JOIN albums ON albums.album_uid = ?", albumUID).
		Joins("JOIN photos_albums pa ON pa.album_uid = albums.album_uid AND pa.photo_uid = files.photo_uid AND pa.hidden = 0").
		Joins("JOIN photos ON photos.id = files.photo_id AND photos.photo_private = 0 AND photos.deleted_at IS NULL").
		Order("photos.photo_quality DESC, photos.taken_at DESC").
		First(&file).Error; err != nil {
		return file, err
	}
	*/
	return file, nil
}

/* TODO with people tag
// PeoplePhotos returns up to count photos from all photos.
func PeoplePhotos(a entity.People, count int) (results PhotoResults, err error) {
	results, _, err = PhotoSearch(form.PhotoSearch{
		Album:  a.AlbumUID,
		Filter: a.AlbumFilter,
		Count:  count,
		Offset: 0,
	})

	return results, err
}
*/
/*
// AlbumSearch searches albums based on their name.
func AlbumSearch(f form.AlbumSearch) (results AlbumResults, err error) {
	if err := f.ParseQueryString(); err != nil {
		return results, err
	}

	defer log.Debug(capture.Time(time.Now(), fmt.Sprintf("albums: search %s", form.Serialize(f, true))))

	s := UnscopedDb().Table("albums").
		Select("albums.*, cp.photo_count,	cl.link_count").
		Joins("LEFT JOIN (SELECT album_uid, count(photo_uid) AS photo_count FROM photos_albums WHERE hidden = 0 AND missing = 0 GROUP BY album_uid) AS cp ON cp.album_uid = albums.album_uid").
		Joins("LEFT JOIN (SELECT share_uid, count(share_uid) AS link_count FROM links GROUP BY share_uid) AS cl ON cl.share_uid = albums.album_uid").
		Where("albums.album_type <> 'folder' OR albums.album_path IN (SELECT photos.photo_path FROM photos WHERE photos.deleted_at IS NULL)").
		Where("albums.deleted_at IS NULL")

	if f.ID != "" {
		s = s.Where("albums.album_uid IN (?)", strings.Split(f.ID, Or))

		if result := s.Scan(&results); result.Error != nil {
			return results, result.Error
		}

		return results, nil
	}

	if f.Query != "" {
		likeString := "%" + f.Query + "%"
		s = s.Where("albums.album_title LIKE ? OR albums.album_location LIKE ?", likeString, likeString)
	}

	if f.Type != "" {
		s = s.Where("albums.album_type IN (?)", strings.Split(f.Type, Or))
	}

	if f.Category != "" {
		s = s.Where("albums.album_category IN (?)", strings.Split(f.Category, Or))
	}

	if f.Location != "" {
		s = s.Where("albums.album_location IN (?)", strings.Split(f.Location, Or))
	}

	if f.Favorite {
		s = s.Where("albums.album_favorite = 1")
	}

	if (f.Year > 0 && f.Year <= txt.YearMax) || f.Year == entity.YearUnknown {
		s = s.Where("albums.album_year = ?", f.Year)
	}

	if (f.Month >= txt.MonthMin && f.Month <= txt.MonthMax) || f.Month == entity.MonthUnknown {
		s = s.Where("albums.album_month = ?", f.Month)
	}

	if (f.Day >= txt.DayMin && f.Month <= txt.DayMax) || f.Day == entity.DayUnknown {
		s = s.Where("albums.album_day = ?", f.Day)
	}

	switch f.Order {
	case "slug":
		s = s.Order("albums.album_favorite DESC, album_slug ASC")
	default:
		s = s.Order("albums.album_favorite DESC, albums.album_year DESC, albums.album_month DESC, albums.album_day DESC, albums.album_title, albums.created_at DESC")
	}

	if f.Count > 0 && f.Count <= MaxResults {
		s = s.Limit(f.Count).Offset(f.Offset)
	} else {
		s = s.Limit(MaxResults).Offset(f.Offset)
	}

	if result := s.Scan(&results); result.Error != nil {
		return results, result.Error
	}

	return results, nil
}*/

// GetPeoples returns a slice of peoples.
func GetPeoples(offset, limit int) (results []entity.People, err error) {
	err = UnscopedDb().Table("peoples").Select("*").Offset(offset).Limit(limit).Find(&results).Error
	return results, err
}