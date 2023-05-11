package person

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func Insert(tx *sqlx.DB, e *Entity) (int, error) {

	if e == nil {
		panic("trying to insert null entity in cpx-package")
	}

	var err error
	var mapp map[string]interface{}
	mapp = map[string]interface{}{
		"record": e,
	}
	sqlStmt, err := mapper.GetMappedStatement("insert", mapp)
	if err != nil {
		return 0, err
	}

	r, err := tx.Exec(sqlStmt.GetStatement(), sqlStmt.GetParams()...)
	if err != nil {
		return 0, err
	}

	n, err := r.RowsAffected()
	if err != nil {
		log.Warn().Err(err).Send()
	}

	return int(n), nil
}