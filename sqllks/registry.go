package sqllks

import (
	"errors"
	"github.com/rs/zerolog/log"
)

type LinkedServices []*LinkedService

var theRegistry LinkedServices

func Initialize(cfgs []Config) (LinkedServices, error) {

	const semLogContext = "sql-registry::initialize"
	if len(cfgs) == 0 {
		log.Info().Msg(semLogContext + " no config provided....skipping")
		return nil, nil
	}

	if len(theRegistry) != 0 {
		log.Warn().Msg(semLogContext + " registry already configured.. overwriting")
	}

	log.Info().Int("no-linked-services", len(cfgs)).Msg(semLogContext)

	var r LinkedServices
	for _, kcfg := range cfgs {
		lks, err := NewServiceInstanceWithConfig(&kcfg)
		if err != nil {
			return nil, err
		}

		r = append(r, lks)
		log.Info().Str("server-name", kcfg.ServerName).Msg(semLogContext + " instance configured")

	}

	theRegistry = r
	return r, nil
}

func Close() {
	const semLogContext = "sql-registry::close"
	log.Info().Msg(semLogContext)
	for _, lks := range theRegistry {
		lks.Close()
	}
}

func GetLinkedService(brokerName string) (*LinkedService, error) {

	const semLogContext = "sql-registry::get-lks"

	log.Trace().Str("server-name", brokerName).Msg(semLogContext)

	for _, k := range theRegistry {
		if k.cfg.ServerName == brokerName {
			return k, nil
		}
	}

	err := errors.New("linked service not found by name " + brokerName)
	log.Error().Err(err).Str("server-name", brokerName).Msg(semLogContext)
	return nil, err
}
