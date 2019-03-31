package middleware

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/routes"
	"net/http"
)

func Panic(next routes.Handler) routes.Handler {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				switch err := e.(type) {
				case error:
					handlers.Handle500(w, r, formats.ErrPanic, err)
				case string:
					handlers.Handle500(w, r, formats.ErrPanic, fmt.Errorf(err))
				default:
					handlers.Handle500(w, r, formats.ErrPanic, fmt.Errorf("%v", e))
				}
			}
		}()
		next.ServeHTTP(w, r)
	}, next.Settings())
}
