package response

type ParamsPagination struct {
	// "page": 1,
	// "limit": 10,
	// "sort": "id.asc",
	// "search": "name",
	// "no_page": "1",

	Page      int
	Limit     int
	SortQuery string
	Search    string
	NoPage    string
}
