package person_test

import (
	"database/sql"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/cmds/sql-cli/examples/person"
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

func TestEntity(t *testing.T) {
	lks, err := sqllks.GetLinkedService("default")
	require.NoError(t, err)

	sqlDb, err := lks.DB()
	require.NoError(t, err)
	sqlDb.MustExec(person.EntityTableDDL)
	defer sqlDb.MustExec(person.EntityTableDropDDL)

	p := person.Entity{
		/*  complete as needed */
		Id:         person.MustToMax20Text("user-id"),
		Lastname:   person.MustToMax20Text("paperino"),
		Nickname:   sql.NullString(person.MustToNullableMax20Text("paolino")),
		Age:        sqlutil.ToSqlNullInt32(61),
		Consensus:  sqlutil.ToSqlNullBool(false),
		CreationTm: sql.NullTime{},
	}

	_, err = person.Insert(sqlDb, &p)
	require.NoError(t, err)

	uopts := []person.UpdateOp{
		person.UpdateWithRowsAffectedWanted(1),
		/*  complete as needed */
		person.UpdateWithLastname(person.MustToMax20Text("paperino")),
		person.UpdateWithNickname(sql.NullString(person.MustToNullableMax20Text("paolino"))),
		person.UpdateWithAge(sqlutil.ToSqlNullInt32(61)),
		person.UpdateWithConsensus(sqlutil.ToSqlNullBool(false)),
		person.UpdateWithCreationTm(sql.NullTime{}),
	}

	_, err = person.Update(sqlDb, person.NewFilterBuilder().Build(), uopts...)
	require.NoError(t, err)
}
