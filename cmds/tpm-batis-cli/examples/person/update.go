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
	Lastname           string
	Nickname           sql.NullString
	Age                sql.NullInt32
	Consensus          sql.NullBool
	CreationTm         sql.NullTime
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

type UpdateOp func(u *UpdateData) error

func UpdateWithRowsAffectedWanted(p int64) UpdateOp {
	return func(u *UpdateData) error {
		u.rowsAffectedWanted = p
		return nil
	}
}

func UpdateWithLastname(p string) UpdateOp {
	return func(u *UpdateData) error {
		var s string
		var err error
		if s, err = ValidateLastname(p); err != nil {
			return err
		}
		u.Lastname = s
		u.flagsDirty.Set(LastnameFieldIndex)
		return nil
	}
}

func (uda *UpdateData) IsLastnameDirty() bool {
	return uda.flagsDirty.Test(LastnameFieldIndex)
}

func UpdateWithNickname(p sql.NullString) UpdateOp {
	return func(u *UpdateData) error {
		var ns sql.NullString
		var err error
		if ns, err = ValidateNickname(p); err != nil {
			return err
		}
		u.Nickname = ns
		u.flagsDirty.Set(NicknameFieldIndex)
		return nil
	}
}

func (uda *UpdateData) IsNicknameDirty() bool {
	return uda.flagsDirty.Test(NicknameFieldIndex)
}

func UpdateWithAge(p sql.NullInt32) UpdateOp {
	return func(u *UpdateData) error {
		var ni sql.NullInt32
		var err error
		if ni, err = ValidateAge(p); err != nil {
			return err
		}
		u.Age = ni
		u.flagsDirty.Set(AgeFieldIndex)
		return nil
	}
}

func (uda *UpdateData) IsAgeDirty() bool {
	return uda.flagsDirty.Test(AgeFieldIndex)
}

func UpdateWithConsensus(p sql.NullBool) UpdateOp {
	return func(u *UpdateData) error {
		var nb sql.NullBool
		var err error
		if nb, err = ValidateConsensus(p); err != nil {
			return err
		}
		u.Consensus = nb
		u.flagsDirty.Set(ConsensusFieldIndex)
		return nil
	}
}

func (uda *UpdateData) IsConsensusDirty() bool {
	return uda.flagsDirty.Test(ConsensusFieldIndex)
}

func UpdateWithCreationTm(p sql.NullTime) UpdateOp {
	return func(u *UpdateData) error {
		var nt sql.NullTime
		var err error
		if nt, err = ValidateCreationTm(p); err != nil {
			return err
		}
		u.CreationTm = nt
		u.flagsDirty.Set(CreationTmFieldIndex)
		return nil
	}
}

func (uda *UpdateData) IsCreationTmDirty() bool {
	return uda.flagsDirty.Test(CreationTmFieldIndex)
}

func Update(sqlDbOrTx interface{}, f sqlmapper.Filter, uops ...UpdateOp) (int, error) {

	const semLogContext = "person::update"
	if len(uops) == 0 {
		return 0, nil
	}

	var err error

	ud := UpdateData{rowsAffectedWanted: -1}
	for _, o := range uops {
		err = o(&ud)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return 0, err
		}
	}

	mapp := ud.BuildWithFilter(f)
	sqlStmt, err := mapper.GetMappedStatement("update", mapp)
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
		return 0, fmt.Errorf("update accepts *sqlx.DB or *sqlx.Tx objects, provided %T", sqlDbOrTx)
	}

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
func UpdateByPrimaryKey(sqlDbOrTx interface{}, pk PrimaryKey, uops ...UpdateOp) (int, error) {

	const semLogContext = "person::update-by-primary-key"

	if len(uops) == 0 {
		return 0, nil
	}

	var err error

	ud := UpdateData{rowsAffectedWanted: -1}
	for _, o := range uops {
		err = o(&ud)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return 0, err
		}
	}

	mapp := ud.BuildWithPrimaryKey(pk)
	sqlStmt, err := mapper.GetMappedStatement("updateByPrimaryKey", mapp)
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
		return 0, fmt.Errorf("update-by-primary-key accepts *sqlx.DB or *sqlx.Tx objects, provided %T", sqlDbOrTx)
	}
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
