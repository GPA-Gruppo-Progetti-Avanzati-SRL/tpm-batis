//go:generate stringer -type=StatementParamType
package xml

import (
	"regexp"
	"strings"
)

const dollarPattern = `\$\{([^\$\}\{]+)\}`
const dashPattern = `\#\{([^\$\}\{]+)\}`

type StatementParamType int

const (
	DollarParam StatementParamType = iota
	DashParam
)

type StatementParam struct {
	Type     StatementParamType
	Expr     string
	MetaInfo string
	Options  []string
}

func (v StatementParam) String() string {
	var sb strings.Builder

	switch v.Type {
	case DollarParam:
		sb.WriteString("${")
	case DashParam:
		sb.WriteString("#{")
	}

	sb.WriteString(v.Expr)
	if v.MetaInfo != "" {
		sb.WriteString(",")
		sb.WriteString(v.MetaInfo)
	}

	for _, o := range v.Options {
		sb.WriteString(",")
		sb.WriteString(o)
	}

	sb.WriteString("}")
	return sb.String()
}

type StatementParamParser struct {
	dollarVarRegExp *regexp.Regexp
	dashVarRegExp   *regexp.Regexp
}

func NewStatementParamParser() (StatementParamParser, error) {

	v := StatementParamParser{}
	var err error
	v.dollarVarRegExp, err = regexp.Compile(dollarPattern)

	if err != nil {
		return v, err
	}

	v.dashVarRegExp, err = regexp.Compile(dashPattern)
	if err != nil {
		return v, err
	}

	return v, nil
}

func (rsolv *StatementParamParser) GetDollarVariables(aLine string) ([]StatementParam, error) {
	return rsolv.parse(rsolv.dollarVarRegExp, DollarParam, aLine)
}

func (rsolv *StatementParamParser) GetDashVariables(aLine string) ([]StatementParam, error) {
	return rsolv.parse(rsolv.dashVarRegExp, DashParam, aLine)
}

func (rsolv *StatementParamParser) parse(parseRegExp *regexp.Regexp, paramType StatementParamType, aLine string) ([]StatementParam, error) {

	var vars []StatementParam
	matches := parseRegExp.FindAllStringSubmatchIndex(aLine, -1)
	for _, m := range matches {
		v := StatementParam{Type: paramType}
		segments := strings.Split(aLine[m[2]:m[3]], ",")
		v.Expr = segments[0]
		if len(segments) >= 2 {
			v.MetaInfo = segments[1]
			if len(segments) > 2 {
				v.Options = segments[2:]
			}
		}
		vars = append(vars, v)
	}

	return vars, nil
}
