package chat_db

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/settings/shared"
)

func DB() *sql.DB {
	return settings.DB()
}
