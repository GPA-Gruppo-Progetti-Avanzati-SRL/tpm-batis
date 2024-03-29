package person

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/sqlmapper"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/sqlutil"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

const (
	NoSqlxDbOrTxParam = "select accepts *sqlx.DB or *sqlx.Tx objects, provided %T"
)

func Select(sqlDbOrTx interface{}, f sqlmapper.Filter) ([]Entity, error) {

	const semLogContext = "person::select"

	var mapp map[string]interface{}
	mapp = map[string]interface{}{
		"filter": f,
	}
	sqlStmt, err := mapper.GetMappedStatement("select", mapp)
	if err != nil {
		log.Fatal().Err(err).Msg(semLogContext)
	}

	var ents []Entity

	switch db := sqlDbOrTx.(type) {
	case *sqlx.DB:
		err = db.Select(&ents, sqlStmt.GetStatement(), sqlStmt.GetParams()...)
	case *sqlx.Tx:
		err = db.Select(&ents, sqlStmt.GetStatement(), sqlStmt.GetParams()...)
	default:
		return nil, fmt.Errorf(NoSqlxDbOrTxParam, sqlDbOrTx)
	}

	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	return ents, nil
}

func Count(sqlDbOrTx interface{}, f sqlmapper.Filter) (int, error) {

	const semLogContext = "person::count"

	var mapp map[string]interface{}
	mapp = map[string]interface{}{
		"filter": f,
	}
	sqlStmt, err := mapper.GetMappedStatement("count", mapp)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	var numRows int
	switch db := sqlDbOrTx.(type) {
	case *sqlx.DB:
		err = db.Get(&numRows, sqlStmt.GetStatement(), sqlStmt.GetParams()...)
	case *sqlx.Tx:
		err = db.Get(&numRows, sqlStmt.GetStatement(), sqlStmt.GetParams()...)
	default:
		return 0, fmt.Errorf(NoSqlxDbOrTxParam, sqlDbOrTx)
	}

	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return 0, err
	}

	return numRows, nil
}
func SelectByPrimaryKey(sqlDbOrTx interface{}, pk PrimaryKey) (Entity, error) {

	const semLogContext = "person::select-by-primary-key"

	var mapp map[string]interface{}
	mapp = map[string]interface{}{
		"pk": pk,
	}

	sqlStmt, err := mapper.GetMappedStatement("selectByPrimaryKey", mapp)
	if err != nil {
		log.Fatal().Err(err).Msg(semLogContext)
	}

	var ent Entity

	switch db := sqlDbOrTx.(type) {
	case *sqlx.DB:
		err = db.Get(&ent, sqlStmt.GetStatement(), sqlStmt.GetParams()...)
	case *sqlx.Tx:
		err = db.Get(&ent, sqlStmt.GetStatement(), sqlStmt.GetParams()...)
	default:
		return Entity{}, fmt.Errorf(NoSqlxDbOrTxParam, sqlDbOrTx)
	}

	if err != nil {
		err = sqlutil.MapSqlError(err)
		log.Error().Err(err).Msg(semLogContext)
		return ent, err
	}

	return ent, nil
}

func CountByPrimaryKey(sqlDbOrTx interface{}, pk PrimaryKey) (int, error) {

	const semLogContext = "person::count-by-primary-key"

	var mapp map[string]interface{}
	mapp = map[string]interface{}{
		"pk": pk,
	}
	sqlStmt, err := mapper.GetMappedStatement("count", mapp)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	var numRows int
	switch db := sqlDbOrTx.(type) {
	case *sqlx.DB:
		err = db.Get(&numRows, sqlStmt.GetStatement(), sqlStmt.GetParams()...)
	case *sqlx.Tx:
		err = db.Get(&numRows, sqlStmt.GetStatement(), sqlStmt.GetParams()...)
	default:
		return 0, fmt.Errorf(NoSqlxDbOrTxParam, sqlDbOrTx)
	}

	if err != nil {
		log.Error().Err(err).Send()
		return 0, err
	}

	return numRows, nil
}
