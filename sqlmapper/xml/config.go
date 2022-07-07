package xml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/system/util"
)

type Config struct {
	XMLName        xml.Name
	ListOfMappers  MappersConfig `xml:"mappers"`
	mapperRegistry map[string]*MapperRootNode
}

type MappersConfig struct {
	Mappers []MapperConfig `xml:"sqlmapper"`
}

type MapperConfig struct {
	XMLName  xml.Name `xml:"sqlmapper"`
	Resource string   `xml:"resource,attr"`
}

func NewConfig(resolver util.ResourceResolver, aResourceName string) (*Config, error) {

	var cfg *Config

	if cfgFileContent, err := resolver.GetResource(aResourceName); err != nil {
		return nil, err
	} else {

		buf := bytes.NewBuffer(cfgFileContent)
		dec := xml.NewDecoder(buf)

		cfg = &Config{mapperRegistry: make(map[string]*MapperRootNode)}
		err := dec.Decode(cfg)
		if err != nil {
			return nil, err
		}

		for _, mi := range cfg.ListOfMappers.Mappers {
			if mc, err := resolver.GetResource(mi.Resource); err != nil {
				return nil, err
			} else {
				if me, err := ParseMapper(string(mc)); err != nil {
					return nil, err
				} else {
					cfg.mapperRegistry[me.Namespace] = me
				}
			}
		}
	}

	return cfg, nil
}

func (c *Config) GetMapperConfig(namespace string) (*MapperRootNode, error) {
	if me, ok := c.mapperRegistry[namespace]; ok {
		return me, nil
	}

	return nil, errors.New(fmt.Sprintf("Mapper Not Found: %s", namespace))
}

func (c *Config) GetMappersConfig() ([]*MapperRootNode, error) {

	if len(c.mapperRegistry) == 0 {
		return nil, nil
	}

	list := make([]*MapperRootNode, 0, len(c.mapperRegistry))
	for _, me := range c.mapperRegistry {
		list = append(list, me)
	}

	return list, nil
}
