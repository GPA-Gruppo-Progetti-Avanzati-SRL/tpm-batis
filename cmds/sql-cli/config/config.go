package config

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

// embedded contains baked in options
//
//go:embed tpm-batis-cli.yml
var embeddedCfg []byte

type Config struct {
	flagSet *flag.FlagSet

	ConfigFile            string `json:"config-file" yaml:"config-file"`
	Version               string `json:"tmpl-version" yaml:"tmpl-version"`
	FormatCode            bool   `json:"format-code" yaml:"format-code"`
	CollectionDefFile     string `json:"collection-def-file" yaml:"collection-def-file"`
	CollectionDefScanPath string `json:"collection-def-scan-path" yaml:"collection-def-scan-path"`
}

var DefaultConfig = Config{
	flagSet:    nil,
	ConfigFile: "",
	Version:    "v1",
	FormatCode: false,
}

type configBuilder struct {
	configFile string
}

type ConfigBuilder interface {
	Build(ctx context.Context) (*Config, error)
	With(cfgFileName string) ConfigBuilder
}

func NewBuilder() ConfigBuilder {
	bld := &configBuilder{}
	return bld
}

func (bld *configBuilder) With(fileName string) ConfigBuilder {
	bld.configFile = fileName
	return bld
}

func (bld *configBuilder) Build(ctx context.Context) (*Config, error) {
	return newConfig(ctx, bld.configFile)
}

func (cfg *Config) String() string {
	return fmt.Sprintf("%#v\n", cfg)
}

func newConfig(_ context.Context, cfgFileName string) (*Config, error) {

	const semLogContext = "config::new"
	cfg := DefaultConfig

	var err error
	ok, err := cfg.readConfigFromByteArray(embeddedCfg)
	if err != nil {
		return nil, err
	} else {
		if ok {
			log.Info().Msg(semLogContext + " embedded configuration Loaded")
		}
	}

	if cfgFileName != "" {
		cfg.ConfigFile = cfgFileName
		ok, err := cfg.readConfigFromFile(cfgFileName, true)
		if err != nil {
			return nil, err
		} else {
			if ok {
				log.Info().Str("fileName", cfgFileName).Msg(semLogContext + " configuration Loaded")
			}
		}
	}

	log.Info().Msg(semLogContext + " initializing flag set")
	cfg.initializeFlagSet()

	currentConfigFile := cfg.ConfigFile

	log.Info().Msg(semLogContext + " parsing cmd line params")
	if err := cfg.flagSet.Parse(os.Args[1:]); err != nil {
		return &cfg, err
	}

	if len(cfg.flagSet.Args()) != 0 {
		log.Warn().Interface("flag", cfg.flagSet.Arg(0)).Msg(semLogContext + " invalid command line flag")
	}

	log.Debug().Str("cfg", fmt.Sprintf("%+v", cfg)).Msg(semLogContext + " command line parsed")

	if cfg.ConfigFile != currentConfigFile {
		/*
		 * Il caricamento dell'ultimo file disponibile non modifica il flagSet. Eventuali flag per path 'dinamici' eventualmente inseriti sono
		 * censiti tra gli errori....
		 */
		if _, err := cfg.readConfigFromFile(cfg.ConfigFile, true); err != nil {
			return &cfg, err
		}
	}

	log.Debug().Str("config", fmt.Sprintf("%+v", cfg)).Msg("configuration Loaded")

	if err := cfg.checkValid(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (cfg *Config) initializeFlagSet() {
	cfg.flagSet = flag.NewFlagSet("sql-cli", flag.ContinueOnError)

	cfg.flagSet.StringVar(&cfg.ConfigFile, "config-file", cfg.ConfigFile, "path to the configuration file.")

	/*
	 */
	cfg.flagSet.BoolVar(&cfg.FormatCode, "format-code", cfg.FormatCode, "boolean: format code")
	cfg.flagSet.StringVar(&cfg.Version, "tmpl-ver", cfg.Version, "Version of templates (v1)")
	cfg.flagSet.StringVar(&cfg.CollectionDefFile, "collection-def-file", cfg.CollectionDefFile, "collection definition filename")
	cfg.flagSet.StringVar(&cfg.CollectionDefScanPath, "collection-def-scan-path", cfg.CollectionDefScanPath, "scan directory for collection definition")
}

func (cfg *Config) readConfigFromFile(aConfigFileName string, mustExists bool) (bool, error) {

	const semLogContext = "config::read-from-file"

	var err error
	if _, err = os.Stat(aConfigFileName); err != nil {
		if os.IsNotExist(err) {
			if !mustExists {
				log.Info().Str("filename", aConfigFileName).Msg(semLogContext + " config file not found.. using default values")
				err = nil
			} else {
				err = fmt.Errorf("config file %s not found", aConfigFileName)
			}
		} else {
			err = fmt.Errorf("config file %s error", aConfigFileName)
		}

		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
		}

		return false, err
	}

	configContent, err := util.ReadFileAndResolveEnvVars(aConfigFileName)
	if err != nil {
		log.Error().Err(err).Str("filename", aConfigFileName).Msg(semLogContext)
		return false, err
	}

	if _, err = cfg.readConfigFromByteArray(configContent); err != nil {
		log.Error().Err(err).Str("filename", aConfigFileName).Msg(semLogContext)
		return false, err
	}

	return true, nil
}

func (cfg *Config) readConfigFromByteArray(configContent []byte) (bool, error) {

	if len(configContent) == 0 {
		return false, nil
	}

	yerr := yaml.Unmarshal(configContent, cfg)
	return true, yerr
}

func (cfg *Config) checkValid() error {

	if cfg.CollectionDefFile == "" && cfg.CollectionDefScanPath == "" {
		return errors.New("no definition files or scan directories specified in config")
	}

	if cfg.CollectionDefFile != "" {
		if f, err := os.Stat(cfg.CollectionDefFile); err != nil {
			return err
		} else {
			if f.IsDir() {
				return errors.New("collection def file is directory")
			}
		}
	}

	if cfg.CollectionDefScanPath != "" {
		if f, err := os.Stat(cfg.CollectionDefScanPath); err != nil {
			return err
		} else {
			if !f.IsDir() {
				return errors.New("collection def file is not directory")
			}
		}
	}

	//if cfg.TargetDirectory == "" {
	//	return errors.New("missing out-dir config")
	//} else if f, err := os.Stat(cfg.TargetDirectory); err != nil {
	//	return err
	//} else {
	//	if !f.IsDir() {
	//		return errors.New("TargetDirectory def file is not directory")
	//	}
	//}

	if cfg.Version != "v1" {
		return errors.New("just have v1 and v2 versions of templates sorry")
	}
	return nil
}

func (cfg *Config) FindCollectionToProcess() ([]string, error) {

	if cfg.CollectionDefFile != "" {
		if _, err := os.Stat(cfg.CollectionDefFile); err == nil {
			return []string{cfg.CollectionDefFile}, nil
		} else {
			return nil, err
		}
	}

	if cfg.CollectionDefScanPath != "" {

		defs := make([]string, 0)
		err := filepath.Walk(cfg.CollectionDefScanPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filterPath(info.Name(), info.IsDir()) {
				if !info.IsDir() {
					log.Debug().Str("name", path).Msg("visited file or dir")
					defs = append(defs, path)
				}
				return nil
			} else {
				if info.IsDir() {
					// fmt.Printf("skipping dir: %+v \n", info.Name())
					log.Debug().Str("name", info.Name()).Msg("skipping dir")
					return filepath.SkipDir
				}
			}

			return nil
		})

		if err != nil {
			return nil, err
		} else if len(defs) > 0 {
			return defs, nil
		}
	}

	return nil, errors.New("no files specified in config")
}

func filterPath(n string, isDir bool) bool {

	if isDir {
		if strings.HasPrefix(n, ".") && n != "." && n != ".." {
			return false
		}

		return true
	} else {
		if strings.HasSuffix(n, "-sql.yml") {
			return true
		}

		return false
	}
}
