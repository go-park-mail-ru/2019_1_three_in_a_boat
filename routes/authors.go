package routes

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
	"net/http"
)

// stores pointers, assumes they aren't modified anywhere else
type authorsResponse = []*db.AuthorData

type AuthorsHandler struct {
}

func (h *AuthorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	authorRows, err := db.GetAllAuthors(_db)
	if err != nil {
		Handle500(w, r, ErrSqlFailure, "db.GetAllAuthors", err)
		return
	}

	var authors authorsResponse

	for authorRows.Next() {
		a, err := db.AuthorDataFromRow(authorRows)
		if err != nil {
			Handle500(w, r, ErrDbScanFailure, "db.AuthorDataFromRow", err)
			return
		}
		authors = append(authors, a)
	}

	Handle200(w, r, authors, "db.AuthorDataFromRow")
}
