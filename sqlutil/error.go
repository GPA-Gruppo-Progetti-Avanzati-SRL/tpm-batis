package sqlutil

import (
	"database/sql"
	"errors"
)

const NoRowsInResultSet = "sql: no rows in result set"

var EntityNotFoundError = errors.New("entity not found")

func MapSqlError(err error) error {

	if err == sql.ErrNoRows {
		return EntityNotFoundError
	}

	return err
}
