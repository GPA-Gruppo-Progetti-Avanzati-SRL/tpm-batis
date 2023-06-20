package person_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/cmds/tpm-batis-cli/examples/person"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/sqlmapper"

	"database/sql"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/sqllks"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/sqlutil"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"os"
	"strconv"
	"testing"
)

const (
	SQLHOSTENVVAR   = "SQL_HOST"
	SQLDBENVVAR     = "SQL_DB"
	SQLUSERENVVAR   = "SQL_USER"
	SQLPASSWDENVVAR = "SQL_PASSWD"
	SQLPORTENVVAR   = "SQL_PORT"
)

func TestMain(m *testing.M) {

	p, err := strconv.Atoi(os.Getenv(SQLPORTENVVAR))
	if err != nil {
		panic(err)
	}

	cfg := sqllks.Config{
		ServerName:      "default",
		ServerType:      "pq",
		Host:            os.Getenv(SQLHOSTENVVAR),
		Port:            p,
		DbName:          os.Getenv(SQLDBENVVAR),
		UserName:        os.Getenv(SQLUSERENVVAR),
		Password:        os.Getenv(SQLPASSWDENVVAR),
		SslMode:         false,
		EnableMigration: false,
		MaxOpenConns:    0,
		MaxIdleConns:    0,
		ConnMaxLifetime: 0,
		ConnMaxIdleTime: 0,
	}

	_, err = sqllks.Initialize([]sqllks.Config{cfg})
	if err != nil {
		panic(err)
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}

const doDDL = false

// This is a simple test skeleton that has to be adjusted to do real useful work. please complete with desired values. */
func TestEntity(t *testing.T) {
	lks, err := sqllks.GetLinkedService("default")
	require.NoError(t, err)

	sqlDb, err := lks.DB()
	require.NoError(t, err)

	if doDDL {
		t.Log("ddl execution")
		sqlDb.MustExec(person.EntityTableDDL)
		defer sqlDb.MustExec(person.EntityTableDropDDL)
	}
	t.Log("insert statement")

	tx := sqlDb.MustBegin()

	p := person.Entity{
		/*  complete as needed */
		Id:         person.MustValidateId("user-id"),
		Lastname:   person.MustValidateLastname("paperino"),
		Nickname:   person.MustValidateNickname("paolino"),
		Age:        sqlutil.ToSqlNullInt32(61),
		Consensus:  sqlutil.ToSqlNullBool(false),
		CreationTm: sql.NullTime{},
	}

	_, err = person.Insert(tx, &p)
	require.NoError(t, sqlutil.OnErrorTxClose(tx, err))

	t.Log("update statement")
	uopts := []person.UpdateOp{
		person.UpdateWithRowsAffectedWanted(1),
		/*  complete as needed */
		person.UpdateWithLastname(person.MustValidateLastname("paperino")),
		person.UpdateWithNickname(person.MustValidateNickname("paolino")),
		person.UpdateWithAge(sqlutil.ToSqlNullInt32(61)),
		person.UpdateWithConsensus(sqlutil.ToSqlNullBool(false)),
		person.UpdateWithCreationTm(sql.NullTime{}),
	}

	_, err = person.Update(tx, person.NewFilterBuilder().Build(), uopts...)
	require.NoError(t, sqlutil.OnErrorTxClose(tx, err))

	t.Log("select statement")
	f := sqlmapper.NewFilterBuilder().Limit(2).Offset(0)
	l, err := person.Select(tx, f.Build())
	require.NoError(t, sqlutil.OnErrorTxClose(tx, err))
	for i, e := range l {
		t.Logf("row: %d - %v", i, e)
	}
	t.Log("select-by-primary-key statement")
	e, err := person.SelectByPrimaryKey(tx, person.PrimaryKey{Id: person.MustValidateId("user-id")})
	require.NoError(t, sqlutil.OnErrorTxClose(tx, err))
	t.Log(e)

	fdel := person.NewFilterBuilder()
	// customize with the proper filters ...
	// fdel.Or().And...
	_, err = person.Delete(tx, fdel.Build())
	require.NoError(t, sqlutil.OnErrorTxClose(tx, err))

	err = sqlutil.TxClose(tx, nil)
	require.NoError(t, err)
}
