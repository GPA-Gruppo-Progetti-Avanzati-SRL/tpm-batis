package sqllks

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"strings"
)

const (
	SeqPQSQLExpression = "select concat('%s', to_char(nextval('%s'), '%sFM'))"
)

/*
	func DBSequenceNextVal(db *sqlx.DB, seqName string, pfix string) (string, error) {
		q := fmt.Sprintf(SeqPQSQLExpression, pfix, seqName, strings.Repeat("0", 12))

		res := ""
		err := db.Get(&res, q)
		return res, err
	}

	func TXSequenceNextVal(db *sqlx.Tx, seqName string, pfix string) (string, error) {
		q := fmt.Sprintf(SeqPQSQLExpression, pfix, seqName, strings.Repeat("0", 12))

		res := ""
		err := db.Get(&res, q)
		return res, err
	}
*/

func PqSequenceNextVal(dbOrTx interface{}, seqName string, pfix string) (string, error) {

	const semLogContext = "sql-lks::seq-next-val"
	q := fmt.Sprintf(SeqPQSQLExpression, pfix, seqName, strings.Repeat("0", 12))

	var res string
	var err error

	switch t := dbOrTx.(type) {
	case *sqlx.DB:
		err = t.Get(&res, q)
	case *sqlx.Tx:
		err = t.Get(&res, q)
	default:
		log.Error().Err(err).Msg(semLogContext)
		err = fmt.Errorf("unsupported type %T: expected *sqlx.DB or *sqlx.Tx", dbOrTx)
	}

	return res, err
}
