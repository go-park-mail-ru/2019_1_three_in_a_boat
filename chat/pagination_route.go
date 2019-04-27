package main

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/chat/chat_db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/settings/chat"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils/handlers"
)

// this is actually a CHAD handler and he is very alpha
type ChatPaginationHandler struct{}

func (h *ChatPaginationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	msgId := r.URL.Query()["msgId"]
	if len(msgId) != 1 {
		HandleInvalidData(
			w, r, formats.ErrInvalidGetParams, "msgId missing or more than 1")
		return
	}
	msgIdInt, err := strconv.Atoi(msgId[0])
	if err != nil {
		HandleInvalidData(
			w, r, formats.ErrInvalidGetParams, "msgId not an integer")
		return
	}

	rows, err := chat_db.GetNMessagesSince(chat_settings.DB(), 20, msgIdInt)
	if HandleErrForward(w, r, formats.ErrSqlFailure, err) != nil {
		return
	}
	//noinspection GoUnhandledErrorResult
	defer rows.Close()

	var ret []*chat_db.Message
	for rows.Next() {
		m, err := chat_db.MessageFromRow(rows)
		if HandleErrForward(w, r, formats.ErrDbScanFailure, err) != nil {
			return
		}
		ret = append(ret, m)
	}

	Handle200(w, r, ret)

}

func (h *ChatPaginationHandler) Settings() map[string]http_utils.RouteSettings {
	return map[string]http_utils.RouteSettings{
		"GET": {
			AuthRequired:           false,
			CorsAllowed:            true,
			CsrfProtectionRequired: false,
		},
	}
}
