package routes

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
	"net/http"
)

// stores pointers, assumes they aren't modified anywhere else
type authorsResponse = []*db.AuthorData

type AuthorsHandler struct {
	db    *sql.DB
	route string
}

func (h *AuthorsHandler) SetDB(_db *sql.DB) {
	h.db = _db
}

func (h *AuthorsHandler) SetRoute(route string) {
	h.route = route
}

func (h *AuthorsHandler) GetRoute() string {
	return h.route
}

func (h *AuthorsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	setDefaultHeaders(res)

	authorRows, err := db.GetAllAuthors(h.db)
	if err != nil {
		handle500(h, "E_DB_FAILURE", "db.GetAllAuthors", err, res)
	}

	var authors authorsResponse

	for authorRows.Next() {
		a, err := db.AuthorDataFromRow(authorRows)
		if err != nil {
			handle500(h, "E_DB_SCAN_FAILURE", "db.AuthorDataFromRow", err, res)
		}
		authors = append(authors, a)
	}

	handle200(h, authors, "db.AuthorDataFromRow", res)
}
