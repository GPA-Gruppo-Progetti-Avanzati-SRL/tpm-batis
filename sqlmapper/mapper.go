//go:generate stringer -type=ParamBindStyle
package sqlmapper

import (
	"errors"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/sqlmapper/xml"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/system/util"
	"github.com/rs/zerolog/log"
)

type Mapper interface {
	GetStatement(statementId string) (Statement, error)
	GetStatementIds() []string
	GetMappedStatement(statementId string, params map[string]interface{}) (MappedStatement, error)
	GetNamespace() string
	// ExecuteQuery(db *sql.DB, statementId string, params map[string]interface{}) error
}

type MapperRegistry interface {
	GetMapper(name string) Mapper
	AddMapper(aName string, m Mapper) error
}

type ParamBindStyle int

const (
	BINDSTYLE_QUESTION ParamBindStyle = iota
	BINDSTYLE_DOLLAR
)

type Option func(mapperOpts *mapperOptions)

type mapperOptions struct {
	bindStyle ParamBindStyle
}

func WithBindStyle(bindStyle ParamBindStyle) Option {
	return func(mapperOpts *mapperOptions) {
		mapperOpts.bindStyle = bindStyle
	}
}

type mapper struct {
	xmlMapper    *xml.MapperRootNode
	stmtRegistry map[string]Statement
	opts         mapperOptions
}

type mapperRegistry struct {
	registry map[string]Mapper
}

func (r *mapperRegistry) GetMapper(aName string) Mapper {
	if m, ok := r.registry[aName]; ok {
		return m
	}
	return nil
}

func (r *mapperRegistry) AddMapper(aName string, m Mapper) error {
	r.registry[aName] = m
	return nil
}

func NewRegistry(resolver util.ResourceResolver, aResourceConfigName string) (MapperRegistry, error) {

	sf := &mapperRegistry{registry: make(map[string]Mapper)}
	cfg, err := xml.NewConfig(resolver, aResourceConfigName)
	if err != nil {
		return sf, err
	}

	if mappersConfig, err := cfg.GetMappersConfig(); err != nil {
		return sf, err
	} else {
		for _, mp := range mappersConfig {
			if m, e := NewMapperFromParsedXML(mp); e == nil {
				sf.registry[m.GetNamespace()] = m
			} else {
				return sf, e
			}
		}
	}

	return sf, nil
}

func NewMapper(mapperSource string, opts ...Option) (Mapper, error) {
	if mapperXmlTree, err := xml.ParseMapper(string(mapperSource)); err != nil {
		return nil, err
	} else {
		return NewMapperFromParsedXML(mapperXmlTree, opts...)
	}
}

func NewMapperFromParsedXML(mapperNodeConfig *xml.MapperRootNode, opts ...Option) (Mapper, error) {

	mo := mapperOptions{bindStyle: BINDSTYLE_QUESTION}
	for _, o := range opts {
		o(&mo)
	}

	m := mapper{xmlMapper: mapperNodeConfig, stmtRegistry: make(map[string]Statement), opts: mo}

	for _, mapperStmtId := range mapperNodeConfig.GetStatementsIds() {
		if s, e := newStatement(m.xmlMapper, mapperStmtId, mo); e != nil {
			return nil, e
		} else {
			m.stmtRegistry[mapperStmtId] = s
		}
	}

	return &m, nil
}

func (m *mapper) GetNamespace() string {
	return m.xmlMapper.Namespace
}

func (m *mapper) GetStatement(statementId string) (Statement, error) {

	if s, ok := m.stmtRegistry[statementId]; ok {
		return s, nil
	}

	return nil, errors.New(fmt.Sprintf("Statement not found %s", statementId))
}

func (m *mapper) GetStatementIds() []string {

	if len(m.stmtRegistry) == 0 {
		return nil
	}

	ids := make([]string, len(m.stmtRegistry))
	i := 0
	for _, s := range m.stmtRegistry {
		ids[i] = s.(*statement).id
		i++
	}

	return ids
}

type MappedStatement interface {
	GetStatement() string
	GetParams() []interface{}
	GetResultMap() ResultMap
}

type mappedStatement struct {
	stmt      string
	params    []PreparedStatementParam
	resultMap ResultMap
}

func (ms *mappedStatement) GetStatement() string {
	return ms.stmt
}

func (ms *mappedStatement) GetParams() []interface{} {
	var qparams []interface{}
	for _, p := range ms.params {
		qparams = append(qparams, p.Value)
	}

	return qparams
}

func (ms *mappedStatement) LogInfo() {

	logEvent := log.Trace().Str("sql", ms.stmt)
	for i, p := range ms.params {
		logEvent.Str(fmt.Sprintf("$%d", i), fmt.Sprint(p))
	}

	logEvent.Msg("tmp-batis statement")
}

func (ms *mappedStatement) GetResultMap() ResultMap {
	return ms.resultMap
}

func (m *mapper) GetMappedStatement(statementId string, params map[string]interface{}) (MappedStatement, error) {
	var s Statement
	var err error
	if s, err = m.GetStatement(statementId); err != nil {
		return nil, err
	}

	var statement string
	var stmtParams []PreparedStatementParam
	if statement, stmtParams, err = s.ExecuteTemplate(params); err != nil {
		return nil, err
	}

	ms := &mappedStatement{statement, stmtParams, s.GetResultMap()}
	ms.LogInfo()
	return ms, nil
}

/*func (m *sqlmapper)ExecuteQuery2(statementId string, params ...interface{}) error {

	p := make(map[string]interface{})
	p["site"] = params[0]
	if statement, stmtParams, err := m.GetStatement(statementId, p); err == nil {

		queryParams := make([]interface{}, 0, len(stmtParams))
		for _, sp := range stmtParams {
			queryParams = append(queryParams, sp.Value)
		}

		if rows, err := m.db.Query(statement, queryParams...); err == nil {
			defer rows.Close()

			for rows.Next() {

				var site string
				err := rows.Scan(&site)
				if err != nil {
					_ = level.Error(m.logger).Log(tpmlog.DefaultLogMessageField, err)
					return err
				}
			}

		} else {
			level.Error(m.logger).Log(tpmlog.DefaultLogMessageField, err)
			return err
		}

	} else {
		return err
	}

	return nil
}*/
