package response

import (
	_ "log"

	"gorm.io/gorm"
)

// type CustomPagination struct {
// 	// "next": null,
// 	// "previous": null,
// 	// "count": 2,
// 	// "total_pages": 1,
// 	// "page": 1,
// 	// "page_size": 10,
// 	// "results": []

// 	QueryParams  string      `json:"queryParams"`
// 	NextPage     *string     `json:"nextPage"`
// 	PreviousPage *string     `json:"previousPage"`
// 	TotalRows    int64       `json:"count"`
// 	Page         int         `json:"currentPage"`
// 	TotalPage    int         `json:"totalPage"`
// 	PageSize     int         `json:"pageSize"`
// 	Items        interface{} `json:"items"`
// }

// func (s CustomPagination) Response(page int, pageSize int, totalPage int, items interface{}, query string) interface{} {
// 	var nextPage, previousPage *string
// 	if page+1 <= totalPage {
// 		params := url.Values{}
// 		params.Add("page", strconv.Itoa(page+1))
// 		params.Add("page_size", strconv.Itoa(pageSize))

// 		splits := strings.Split(s.QueryParams, "&")
// 		for _, prm := range splits {
// 			x := strings.Split(prm, "=")
// 			if len(x) == 2 {
// 				// fmt.Println(x[0])
// 				// fmt.Println(x[1])
// 				params.Add(x[0], string(x[1]))
// 			}
// 		}

// 		link := query + "?" + params.Encode()
// 		nextPage = &link
// 	}
// 	if page-1 > 0 {
// 		params := url.Values{}
// 		params.Add("page", strconv.Itoa(page-1))
// 		params.Add("page_size", strconv.Itoa(pageSize))
// 		splits := strings.Split(s.QueryParams, "&")

// 		for _, prm := range splits {
// 			x := strings.Split(prm, "=")
// 			if len(x) == 2 {
// 				// fmt.Println(x[0])
// 				// fmt.Println(x[1])
// 				params.Add(x[0], string(x[1]))
// 			}
// 		}

// 		link := query + "?" + params.Encode()
// 		previousPage = &link
// 	}

// 	// s.QueryParams = query
// 	s.NextPage = nextPage
// 	s.PreviousPage = previousPage
// 	s.Page = page
// 	s.TotalPage = totalPage
// 	s.PageSize = pageSize
// 	s.Items = items
// 	return s
// }

type Pagination struct {
	// "next": null,
	// "previous": null,
	// "count": 2,
	// "total_pages": 1,
	// "page": 1,
	// "page_size": 10,
	// "results": []

	// NextPage     *string     `json:"next"`
	// PreviousPage *string     `json:"previous"`
	Count      int64       `json:"count"`
	TotalPages int         `json:"total_pages"`
	Page       int         `json:"page,omitempty;query:page"`
	Limit      int         `json:"limit,omitempty;query:limit"`
	Sort       string      `json:"sort,omitempty;query:sort"`
	Data       interface{} `json:"data"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "Id desc"
	}
	return p.Sort
}

func Paginate(value interface{}, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.Count = totalRows
	totalPages := calculateTotalPage(int(totalRows), pagination.GetLimit())
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

func PaginateRow(totalRows int64, pagination *Pagination) func(db *gorm.DB) *gorm.DB {
	pagination.Count = totalRows
	totalPages := calculateTotalPage(int(totalRows), pagination.GetLimit())
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

func calculateTotalPage(totalRows, pageSize int) (totalPage int) {
	totalPage = totalRows / pageSize
	if totalRows%pageSize > 0 {
		totalPage++
	}
	return
}
