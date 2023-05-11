package person

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/sqlmapper"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func Select(sqlDbOrTx interface{}, f sqlmapper.Filter) ([]Entity, error) {

	const semLogContext = "person::select"

	var mapp map[string]interface{}
	mapp = map[string]interface{}{
		"filter": f,
	}
	sqlStmt, err := mapper.GetMappedStatement("select", mapp)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	var ents []Entity

	switch db := sqlDbOrTx.(type) {
	case *sqlx.DB:
		err = db.Select(&ents, sqlStmt.GetStatement(), sqlStmt.GetParams()...)
	case *sqlx.Tx:
		err = db.Select(&ents, sqlStmt.GetStatement(), sqlStmt.GetParams()...)
	default:
		return nil, fmt.Errorf("select accepts *sqlx.DB or *sqlx.Tx objects, provided %T", sqlDbOrTx)
	}

	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	return ents, nil
}
