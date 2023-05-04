package sqllks_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/sqllks"
	_ "github.com/lib/pq"
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
