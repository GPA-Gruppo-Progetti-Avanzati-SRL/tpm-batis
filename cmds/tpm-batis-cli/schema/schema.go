package schema

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

func ReadSchemaFromFile(fn string) (*Schema, error) {

	const semLogContext = "schema::read-from-file"

	log.Info().Str("filename", fn).Msg(semLogContext)

	b, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	sch, err := readSchema(b)
	if err != nil {
		return nil, err
	}

	sch.DefFileName = fn
	return sch, nil
}

func readSchema(def []byte) (*Schema, error) {

	const semLogContext = "schema::read-from-byte-array"
	log.Info().Msg(semLogContext)

	schema := &Schema{}

	err := yaml.Unmarshal(def, &schema)
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func (sch *Schema) OutputPath(createIfNotExists bool) (string, error) {

	contentPath := filepath.Dir(sch.DefFileName)
	if sch.Properties.FolderPath != "" && sch.Properties.FolderPath != "." {
		contentPath = filepath.Join(contentPath, sch.Properties.FolderPath)
	}

	if _, err := os.Stat(contentPath); os.IsNotExist(err) {
		if createIfNotExists {
			if err = os.MkdirAll(contentPath, 0755); err != nil {
				return "", err
			}
		}
	}

	return contentPath, nil
}

func (sch *Schema) PrimaryKey() []Field {
	var pk []Field
	for _, f := range sch.Fields {
		if f.IsPKey {
			pk = append(pk, f)
		}
	}

	return pk
}

func (sch *Schema) HasPrimaryKey() bool {
	for _, f := range sch.Fields {
		if f.IsPKey {
			return true
		}
	}

	return false
}

func (sch *Schema) PackageName() string {
	const semLogContext = "schema::package-name"
	ndx := strings.LastIndex(sch.Properties.Package, "/")
	if ndx >= 0 {
		return sch.Properties.Package[ndx+1:]
	}

	log.Error().Err(fmt.Errorf("malformed package value (%s)", sch.Properties.Package)).Msg(semLogContext)
	return "malformedPackage"
}
