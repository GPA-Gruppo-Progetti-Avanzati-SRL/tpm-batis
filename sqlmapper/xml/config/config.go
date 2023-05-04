package config

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/sqlmapper"
	xml2 "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/sqlmapper/xml"
)

type Config struct {
	XMLName        xml.Name
	ListOfMappers  MappersConfig `xml:"mappers"`
	mapperRegistry map[string]*xml2.MapperRootNode
}

type MappersConfig struct {
	Mappers []MapperConfig `xml:"sqlmapper"`
}

type MapperConfig struct {
	XMLName  xml.Name `xml:"sqlmapper"`
	Resource string   `xml:"resource,attr"`
}

func NewConfig(resolver sqlmapper.ResourceResolver, aResourceName string) (*Config, error) {

	var cfg *Config

	if cfgFileContent, err := resolver.GetResource(aResourceName); err != nil {
		return nil, err
	} else {

		buf := bytes.NewBuffer(cfgFileContent)
		dec := xml.NewDecoder(buf)

		cfg = &Config{mapperRegistry: make(map[string]*xml2.MapperRootNode)}
		err := dec.Decode(cfg)
		if err != nil {
			return nil, err
		}

		for _, mi := range cfg.ListOfMappers.Mappers {
			if mc, err := resolver.GetResource(mi.Resource); err != nil {
				return nil, err
			} else {
				if me, err := xml2.ParseMapper(string(mc)); err != nil {
					return nil, err
				} else {
					cfg.mapperRegistry[me.Namespace] = me
				}
			}
		}
	}

	return cfg, nil
}

func (c *Config) GetMapperConfig(namespace string) (*xml2.MapperRootNode, error) {
	if me, ok := c.mapperRegistry[namespace]; ok {
		return me, nil
	}

	return nil, errors.New(fmt.Sprintf("Mapper Not Found: %s", namespace))
}

func (c *Config) GetMappersConfig() ([]*xml2.MapperRootNode, error) {

	if len(c.mapperRegistry) == 0 {
		return nil, nil
	}

	list := make([]*xml2.MapperRootNode, 0, len(c.mapperRegistry))
	for _, me := range c.mapperRegistry {
		list = append(list, me)
	}

	return list, nil
}
