package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/cmds/tpm-batis-cli/config"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/cmds/tpm-batis-cli/gen"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/cmds/tpm-batis-cli/schema"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

//go:embed VERSION
var version string

// appLogo contains the ASCII splash screen
//
//go:embed app-logo.txt
var appLogo []byte

func main() {
	const semLogContext = "sql-cli::main"

	fmt.Println(string(appLogo))
	fmt.Printf("Version: %s\n", version)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg, err := config.NewBuilder().Build(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg(semLogContext)
	}

	cDefs, err := cfg.FindCollectionToProcess()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	for _, collDefFile := range cDefs {
		log.Info().Str("def-file", collDefFile).Msg(semLogContext)

		sch, err := schema.ReadSchemaFromFile(collDefFile)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return
		}

		outputPath, err := sch.OutputPath(true)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return
		}

		g, err := gen.NewGenerator(sch, gen.WithOutputFolder(outputPath), gen.WithVersion(cfg.Version), gen.WithFormatCode(cfg.FormatCode))
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return
		}

		err = g.Generate()
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return
		}
	}
}
