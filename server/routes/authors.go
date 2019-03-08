package routes

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/settings"
	"net/http"
)

// stores pointers, assumes they aren't modified anywhere else
type authorsResponse = []*db.AuthorData

type AuthorsHandler struct{}

// Handler for the Authors resource
func (h *AuthorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	rows, err := db.GetAllAuthors(settings.DB())
	if err != nil {
		handlers.Handle500(w, r, formats.ErrSqlFailure, "db.GetAllAuthors", err)
		return
	} else {
		defer rows.Close()
	}

	authors := authorsResponse{}

	for rows.Next() {
		a, err := db.AuthorDataFromRow(rows)
		if err != nil {
			handlers.Handle500(w, r, formats.ErrDbScanFailure,
				"db.AuthorDataFromRow", err)
			return
		}
		authors = append(authors, a)
	}

	if err := rows.Err(); err != nil {
		handlers.Handle500(w, r, formats.ErrDbRowsFailure, "db.GetAllAuthors", err)
		return
	}

	handlers.Handle200(w, r, authors, "db.AuthorDataFromRow")
}

func (h *AuthorsHandler) Methods() map[string]struct{} {
	return map[string]struct{}{"GET": {}}
}

func (h *AuthorsHandler) AuthRequired(method string) bool {
	return false
}

func (h *AuthorsHandler) CorsAllowed(method string) bool {
	return true
}
