package form

// AlbumSearch represents search form fields for "/api/v1/albums".
type PeopleSearch struct {
	Query    string `form:"q"`
	ID       string `form:"id"`
	UID      string `form:"uid"`
	FullName string `form:"full_name"`
	Category string `form:"category"`
	Slug     string `form:"slug"`
	BoDYear  int    `json:"bod_year"`
	BoDMonth int    `json:"bod_month"`
	BoDDay   int    `json:"bod_day"`
	Age      int    `json:"age"`
	Count    int    `form:"count" serialize:"-"`
	Offset   int    `form:"offset" serialize:"-"`
	Order    string `form:"order" serialize:"-"`
}

func (f *PeopleSearch) GetQuery() string {
	return f.Query
}

func (f *PeopleSearch) SetQuery(q string) {
	f.Query = q
}

func (f *PeopleSearch) ParseQueryString() error {
	return ParseQueryString(f)
}

func NewPeopleSearch(query string) PeopleSearch {
	return PeopleSearch{Query: query}
}
