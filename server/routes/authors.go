package routes

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/settings"
	"net/http"
)

// Handles Authors resource. Only accepts GET requests because creating authors
// is currently not implemented and not needed. Implements routes.Handler
// interface, which extends http.Handler.
type AuthorsHandler struct{}

// Handles GET requests for the authors resource. Assumes method is already
// filtered by the Methods middleware. In case Of a successful request, returns
// []*db.AuthorData
func (h *AuthorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	rows, err := db.GetAllAuthors(settings.DB())
	if HandleErrForward(w, r, formats.ErrSqlFailure, err) != nil {
		return
	} else {
		defer rows.Close()
	}

	authors := make([]*db.AuthorData, 0)

	for rows.Next() {
		a, err := db.AuthorDataFromRow(rows)
		if HandleErrForward(w, r, formats.ErrDbScanFailure, err) != nil {
			return
		}
		authors = append(authors, a)
	}

	if err := rows.Err(); HandleErrForward(w, r, formats.ErrDbRowsFailure, err) != nil {
		return
	}

	Handle200(w, r, authors)
}

func (h *AuthorsHandler) Settings() map[string]RouteSettings {
	return map[string]RouteSettings{
		"GET": {
			AuthRequired:           false,
			CorsAllowed:            true,
			CsrfProtectionRequired: false,
		},
	}
}
