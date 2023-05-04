package sqlutil

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func TxRollback(tx *sqlx.Tx) {
	err := tx.Rollback()
	if err != nil {
		log.Error().Err(err).Msg("transaction rollback")
	}
}

func TxClose(tx *sqlx.Tx, err error) error {
	if err == nil {
		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("transaction commit")
		}

		return err
	}

	TxRollback(tx)
	return err
}
