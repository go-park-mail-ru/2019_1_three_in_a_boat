package routes

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
	"strconv"
	"strings"
)

// Takes string slice out (e.g., from URL.Query()) and makes a valid
// []db.SelectOrder, ignores fields that aren't keys in allowed
func makeOrderList(orderParams []string, allowed map[string]string) []db.SelectOrder {
	order := make([]db.SelectOrder, 0, len(orderParams))

	for _, field := range orderParams {
		if field[0] != '-' {
			field = strings.Title(field)
			if _, ok := allowed[field]; ok {
				order = append(order, db.SelectOrder{field, false})
			}
		} else {
			field = strings.Title(field[1:])
			if _, ok := allowed[field]; ok {
				order = append(order, db.SelectOrder{field, true})
			}
		}
	}

	return order
}

// Takes string slice out (e.g., from URL.Query()) and makes a valid
// page size, s.t. 0 < pageSize <= maxSize, defaulting to defaultSize in case
// the slice is empty or invalid. ignores pageSizeParams[1:], if any.
func makePageSize(pageSizeParams []string, maxSize int, defaultSize int) int {
	var pageSize int
	if len(pageSizeParams) == 0 {
		pageSize = defaultSize
	} else {
		var err error
		pageSize, err = strconv.Atoi(pageSizeParams[0])
		if err != nil {
			pageSize = defaultSize
		} else if pageSize > maxSize {
			pageSize = maxSize
		} else if pageSize <= 0 {
			pageSize = defaultSize
		}
	}

	return pageSize
}

// Takes string slice out (e.g., from URL.Query()) and makes a valid
// page number, s.t. 0 <= nPage, defaulting to 0 in case the slice is empty or
// invalid. ignores pageSizeParams[1:], if any.
func makeNPage(nPageParams []string) int {
	var page int
	if len(nPageParams) == 0 {
		page = 0
	} else {
		var err error
		page, err = strconv.Atoi(nPageParams[0])
		if err != nil || page < 0 {
			page = 0
		}
	}

	return page
}
