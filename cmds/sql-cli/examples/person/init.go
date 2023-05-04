package person

import (
	_ "embed"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/sqlmapper"
	"github.com/rs/zerolog/log"
)

func init() {
	const semLogContext = "person::init"
	log.Info().Msg(semLogContext)

	_ = LoadMapper()
}

//go:embed mapper.xml
var xmlMapper string
var mapper sqlmapper.Mapper

func LoadMapper() error {
	const semLogContext = "person::load-mapper"
	var err error
	if mapper, err = sqlmapper.NewMapper(xmlMapper, sqlmapper.WithBindStyle(sqlmapper.BINDSTYLE_DOLLAR)); err != nil {
		log.Fatal().Err(err).Msg(semLogContext)
	}

	return nil
}
