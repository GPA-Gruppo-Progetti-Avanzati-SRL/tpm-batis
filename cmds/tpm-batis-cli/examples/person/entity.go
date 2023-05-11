package person

import (
	"database/sql"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/sqlutil"
)

const (
	IdFieldIndex         = 0
	LastnameFieldIndex   = 1
	NicknameFieldIndex   = 2
	AgeFieldIndex        = 3
	ConsensusFieldIndex  = 4
	CreationTmFieldIndex = 5
)

type Entity struct {
	Id         Max20Text      `db:"id"`
	Lastname   Max20Text      `db:"lastname"`
	Nickname   sql.NullString `db:"nickname"`
	Age        sql.NullInt32  `db:"age"`
	Consensus  sql.NullBool   `db:"consensus"`
	CreationTm sql.NullTime   `db:"creation_tm"`
}

type PrimaryKey struct {
	Id Max20Text `db:"id"`
}

// isLengthRestrictionValid utility function for Max??Text types
func isLengthRestrictionValid(s string, length, minLength, maxLength int) bool {
	if length > 0 && len(s) != length {
		return false
	}

	if minLength > 0 && len(s) < minLength {
		return false
	}

	if maxLength > 0 && len(s) > maxLength {
		return false
	}

	return true
}

/*
 * Max20Text Type Def
 */

type Max20Text string

const (
	Max20TextZero      = ""
	Max20TextLength    = 0
	Max20TextMinLength = 1
	Max20TextMaxLength = 20
)

// IsValid checks if Max105Text of type String is valid
func (t Max20Text) IsValid() bool {
	return isLengthRestrictionValid(t.ToString(), Max20TextLength, Max20TextMinLength, Max20TextMaxLength)
}

// ToString method for easy conversion
func (t Max20Text) ToString() string {
	return string(t)
}

// ToMax20Text  method for easy conversion with application of restrictions
func ToMax20Text(i interface{}) (Max20Text, error) {

	s := ""
	switch ti := i.(type) {
	case fmt.Stringer:
		s = ti.String()
	case string:
		s = ti
	default:
		return "", fmt.Errorf("")
	}
	if !isLengthRestrictionValid(s, Max20TextLength, Max20TextMinLength, Max20TextMaxLength) {
		return "", fmt.Errorf("cannot satisfy length restriction for %s of type Max20Text", s)
	}

	return Max20Text(s), nil
}

// MustToMax20Text  method for easy conversion with application of restrictions. Panics on error.
func MustToMax20Text(s interface{}) Max20Text {
	v, err := ToMax20Text(s)
	if err != nil {
		panic(err)
	}

	return v
}

/*
 * NullableMax20Text Type Def
 */

type NullableMax20Text sql.NullString

const (
	NullableMax20TextZero      = ""
	NullableMax20TextLength    = 0
	NullableMax20TextMinLength = 0
	NullableMax20TextMaxLength = 20
)

// IsValid checks if Max105Text of type String is valid
func (t NullableMax20Text) IsValid() bool {
	return isLengthRestrictionValid(t.ToString(), NullableMax20TextLength, NullableMax20TextMinLength, NullableMax20TextMaxLength)
}

// ToString method for easy conversion
func (t NullableMax20Text) ToString() string {
	if t.Valid {
		return t.String
	}
	return ""
}

// ToNullableMax20Text  method for easy conversion with application of restrictions
func ToNullableMax20Text(i interface{}) (NullableMax20Text, error) {

	s := ""
	switch ti := i.(type) {
	case sql.NullString:
		if ti.Valid {
			s = ti.String
		}

	case fmt.Stringer:
		s = ti.String()
	case string:
		s = ti
	default:
		return NullableMax20Text(sql.NullString{}), fmt.Errorf("")
	}

	if !isLengthRestrictionValid(s, NullableMax20TextLength, NullableMax20TextMinLength, NullableMax20TextMaxLength) {
		return NullableMax20Text(sql.NullString{}), fmt.Errorf("cannot satisfy length restriction for %s of type NullableMax20Text", s)
	}

	return NullableMax20Text(sqlutil.ToSqlNullString(s)), nil
}

// MustToNullableMax20Text  method for easy conversion with application of restrictions. Panics on error.
func MustToNullableMax20Text(s interface{}) NullableMax20Text {
	v, err := ToNullableMax20Text(s)
	if err != nil {
		panic(err)
	}

	return v
}
