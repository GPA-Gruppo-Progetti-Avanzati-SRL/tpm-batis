package person

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func Insert(sqlDbOrTx interface{}, e *Entity) (int, error) {

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

	var r sql.Result
	switch db := sqlDbOrTx.(type) {
	case *sqlx.DB:
		r, err = db.Exec(sqlStmt.GetStatement(), sqlStmt.GetParams()...)
	case *sqlx.Tx:
		r, err = db.Exec(sqlStmt.GetStatement(), sqlStmt.GetParams()...)
	default:
		return 0, fmt.Errorf("insert accepts *sqlx.DB or *sqlx.Tx objects, provided %T", sqlDbOrTx)
	}

	if err != nil {
		return 0, err
	}

	n, err := r.RowsAffected()
	if err != nil {
		log.Warn().Err(err).Send()
	}

	return int(n), nil
}
