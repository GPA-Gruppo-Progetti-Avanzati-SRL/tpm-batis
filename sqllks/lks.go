package sqllks

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type LinkedService struct {
	cfg   *Config
	sqlDb *sqlx.DB
}

func NewServiceInstanceWithConfig(cfg *Config) (*LinkedService, error) {
	lks := LinkedService{cfg: cfg}
	return &lks, nil
}

func (lks *LinkedService) Close() {

}

func (lks *LinkedService) DB() (*sqlx.DB, error) {

	var err error
	if lks.sqlDb == nil {
		lks.sqlDb, err = lks.newDB(lks.cfg)
	}

	return lks.sqlDb, err
}

func (lks *LinkedService) newDB(cfg *Config) (*sqlx.DB, error) {

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		cfg.Host, cfg.Port, cfg.UserName, cfg.Password, cfg.DbName)

	sqlDb, err := sqlx.Connect("postgres", psqlInfo)
	return sqlDb, err
}

func (lks *LinkedService) SequenceNextVal(seqName, seqPrefix string) (string, error) {

	db, err := lks.DB()
	if err != nil {
		return "", err
	}

	return PqSequenceNextVal(db, seqName, seqPrefix)
}
