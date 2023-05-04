package person

import (
	"database/sql"
)

const (
	IdFieldIndex       = 0
	LastnameFieldIndex = 1
	NicknameFieldIndex = 2
	AgeFieldIndex      = 3
)

type Entity struct {
	Id       string         `db:"id"`
	Lastname string         `db:"lastname"`
	Nickname sql.NullString `db:"nickname"`
	Age      sql.NullInt32  `db:"age"`
}

type PrimaryKey struct {
	Id string `db:"id"`
}
