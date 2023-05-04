package person

import (
	"database/sql"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/sqlmapper"
	"github.com/bits-and-blooms/bitset"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type UpdateData struct {
	flagsDirty         bitset.BitSet
	rowsAffectedWanted int64
	Id                 string
	Lastname           string
	Nickname           sql.NullString
	Age                sql.NullInt32
}

func (uda *UpdateData) BuildWithFilter(f sqlmapper.Filter) map[string]interface{} {
	var mapp map[string]interface{}
	mapp = map[string]interface{}{
		"record": uda,
		"filter": f,
	}

	return mapp
}

func (uda *UpdateData) BuildWithPrimaryKey(pk PrimaryKey) map[string]interface{} {
	var mapp map[string]interface{}
	mapp = map[string]interface{}{
		"record": uda,
		"pk":     pk,
	}

	return mapp
}

type UpdateOp func(u *UpdateData)

func UpdateWithRowsAffectedWanted(p int64) UpdateOp {
	return func(u *UpdateData) {
		u.rowsAffectedWanted = p
	}
}

func UpdateWithLastname(p string) UpdateOp {
	return func(u *UpdateData) {
		u.Lastname = p
		u.flagsDirty.Set(LastnameFieldIndex)
	}
}

func (uda *UpdateData) IsLastnameDirty() bool {
	return uda.flagsDirty.Test(LastnameFieldIndex)
}

func UpdateWithNickname(p sql.NullString) UpdateOp {
	return func(u *UpdateData) {
		u.Nickname = p
		u.flagsDirty.Set(NicknameFieldIndex)
	}
}

func (uda *UpdateData) IsNicknameDirty() bool {
	return uda.flagsDirty.Test(NicknameFieldIndex)
}

func UpdateWithAge(p sql.NullInt32) UpdateOp {
	return func(u *UpdateData) {
		u.Age = p
		u.flagsDirty.Set(AgeFieldIndex)
	}
}

func (uda *UpdateData) IsAgeDirty() bool {
	return uda.flagsDirty.Test(AgeFieldIndex)
}

func Update(tx *sqlx.DB, f sqlmapper.Filter, uops ...UpdateOp) (int, error) {

	const semLogContext = "person::update"
	if len(uops) == 0 {
		return 0, nil
	}

	var err error

	ud := UpdateData{rowsAffectedWanted: -1}
	for _, o := range uops {
		o(&ud)
	}

	mapp := ud.BuildWithFilter(f)
	sqlStmt, err := mapper.GetMappedStatement("update", mapp)
	if err != nil {
		return 0, err
	}

	r, err := tx.Exec(sqlStmt.GetStatement(), sqlStmt.GetParams()...)
	if err != nil {
		return 0, err
	}

	n, err := r.RowsAffected()
	if err != nil {
		log.Warn().Err(err).Msg(semLogContext)
		return 0, err
	}

	if ud.rowsAffectedWanted >= 0 {
		if n != ud.rowsAffectedWanted {
			err = fmt.Errorf("mismatch between number of rows affected (%d) and wanted (%d)", n, ud.rowsAffectedWanted)
			log.Error().Err(err).Msg(semLogContext)
			return 0, err
		}
	}

	return int(n), nil
}

func UpdateByPrimaryKey(tx *sqlx.DB, pk PrimaryKey, uops ...UpdateOp) (int, error) {

	const semLogContext = "person::update-by-primary-key"

	if len(uops) == 0 {
		return 0, nil
	}

	var err error

	ud := UpdateData{rowsAffectedWanted: -1}
	for _, o := range uops {
		o(&ud)
	}

	mapp := ud.BuildWithPrimaryKey(pk)
	sqlStmt, err := mapper.GetMappedStatement("updateByPrimaryKey", mapp)
	if err != nil {
		return 0, err
	}

	r, err := tx.Exec(sqlStmt.GetStatement(), sqlStmt.GetParams()...)
	if err != nil {
		return 0, err
	}

	n, err := r.RowsAffected()
	if err != nil {
		log.Warn().Err(err).Msg(semLogContext)
		return 0, err
	}

	if ud.rowsAffectedWanted >= 0 {
		if n != ud.rowsAffectedWanted {
			err = fmt.Errorf("mismatch between number of rows affected (%d) and wanted (%d)", n, ud.rowsAffectedWanted)
			log.Error().Err(err).Msg(semLogContext)
			return 0, err
		}
	}

	return int(n), nil
}
