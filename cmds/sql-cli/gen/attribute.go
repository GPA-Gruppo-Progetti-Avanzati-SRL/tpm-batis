package gen

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/cmds/sql-cli/schema"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/rs/zerolog/log"
	"strings"
)

type Attribute interface {
	GetDefinition() schema.Field
	DbAttributeName() string
	GoAttributeName() string
	GoAttributeType() string
	GoPackageImports() []string
	GoSampleValue() string

	/*
		GetName(qualified bool, prefixed bool) string
		GetGoPackageImports() []string


		GetGoAttributeIsZeroCondition() string
	*/
}

type AttributeImpl struct {
	AttrDefinition schema.Field
}

func (a AttributeImpl) GetDefinition() schema.Field {
	return a.AttrDefinition
}

func (a AttributeImpl) DbAttributeName() string {
	return strings.ToLower(a.AttrDefinition.DbName)
}

func (a AttributeImpl) GoAttributeName() string {
	return util.Capitalize(a.AttrDefinition.Name)
}

func (a AttributeImpl) GoAttributeType() string {

	const semLogContext = "attribute::type"
	var s string
	switch a.AttrDefinition.Typ {
	case schema.AttributeTypeString:
		s = "string"
	case schema.AttributeTypeInt:
		s = "int32"
	case schema.AttributeTypeNullableString:
		s = "sql.NullString"
	case schema.AttributeTypeNullableInt:
		s = "sql.NullInt32"
	default:
		log.Error().Str("type", string(a.AttrDefinition.Typ)).Msg(semLogContext + " unsupported attribute")
	}

	return s
}

type StringAttribute struct {
	AttributeImpl
}

func (a StringAttribute) GoPackageImports() []string {
	return nil
}

func (a StringAttribute) GoSampleValue() string {
	if a.GetDefinition().SampleValue != "" {
		return fmt.Sprintf("\"%s\"", a.GetDefinition().SampleValue)
	}
	return "\"hello\""
}

type NullStringAttribute struct {
	AttributeImpl
}

func (a NullStringAttribute) GoPackageImports() []string {
	return []string{"database/sql"}
}

func (a NullStringAttribute) GoSampleValue() string {
	if a.GetDefinition().SampleValue != "" {
		if strings.Contains(a.GetDefinition().SampleValue, "sql.Null") {
			return a.GetDefinition().SampleValue
		}
		return fmt.Sprintf("sqlutil.ToSqlNullString(\"%s\")", a.GetDefinition().SampleValue)
	}
	return "sql.NullString{}"
}

type IntAttribute struct {
	AttributeImpl
}

func (a IntAttribute) GoPackageImports() []string {
	return nil
}

func (a IntAttribute) GoSampleValue() string {
	if a.GetDefinition().SampleValue != "" {
		return a.GetDefinition().SampleValue
	}
	return "12"
}

type NullIntAttribute struct {
	AttributeImpl
}

func (a NullIntAttribute) GoPackageImports() []string {
	return []string{"database/sql"}
}

func (a NullIntAttribute) GoSampleValue() string {
	if a.GetDefinition().SampleValue != "" {
		return a.GetDefinition().SampleValue
	}
	return "sql.NullInt32{}"
}

func NewAttribute(attrDefinition schema.Field) (Attribute, error) {

	var a Attribute

	// Disable the update operation if is primary key
	if attrDefinition.IsPKey {
		attrDefinition.WithUpdate = false
	}

	switch attrDefinition.Typ {
	case schema.AttributeTypeString:
		if attrDefinition.Nullable {
			attrDefinition.Typ = schema.AttributeTypeNullableString
			a = NullStringAttribute{AttributeImpl{AttrDefinition: attrDefinition}}
		} else {
			a = StringAttribute{AttributeImpl{AttrDefinition: attrDefinition}}
		}

	case schema.AttributeTypeInt:
		if attrDefinition.Nullable {
			attrDefinition.Typ = schema.AttributeTypeNullableInt
			a = NullIntAttribute{AttributeImpl{AttrDefinition: attrDefinition}}
		} else {
			a = IntAttribute{AttributeImpl{AttrDefinition: attrDefinition}}
		}

	default:
		panic(fmt.Errorf("unsupported attribute type %s", attrDefinition.Typ))
	}

	return a, nil
}
