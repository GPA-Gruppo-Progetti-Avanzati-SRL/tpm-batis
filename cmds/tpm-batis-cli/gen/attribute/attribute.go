package attribute

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/cmds/tpm-batis-cli/schema"
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
	GoSampleValue(pkg string) string
	GoAttributeName4Param() string
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

func (a AttributeImpl) GoAttributeName4Param() string {
	return a.AttrDefinition.Name
}

func (a AttributeImpl) GoAttributeType() string {

	const semLogContext = "attribute::type"
	var s string
	switch a.AttrDefinition.Typ {
	case schema.AttributeTypeString:
		s = "string" // fmt.Sprintf("Max%dText", a.GetDefinition().MaxLength)
	case schema.AttributeTypeInt:
		s = "int32"
	case schema.AttributeTypeBool:
		s = "bool"
	case schema.AttributeTypeTime:
		s = "time.Time"
	case schema.AttributeTypeNullableString:
		s = "sql.NullString"
	case schema.AttributeTypeNullableInt:
		s = "sql.NullInt32"
	case schema.AttributeTypeNullableBool:
		s = "sql.NullBool"
	case schema.AttributeTypeNullableTime:
		s = "sql.NullTime"
	default:
		log.Error().Str("type", string(a.AttrDefinition.Typ)).Msg(semLogContext + " unsupported attribute")
	}

	return s
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

	case schema.AttributeTypeBool:
		if attrDefinition.Nullable {
			attrDefinition.Typ = schema.AttributeTypeNullableBool
			a = NullBoolAttribute{AttributeImpl{AttrDefinition: attrDefinition}}
		} else {
			a = BoolAttribute{AttributeImpl{AttrDefinition: attrDefinition}}
		}

	case schema.AttributeTypeInt:
		if attrDefinition.Nullable {
			attrDefinition.Typ = schema.AttributeTypeNullableInt
			a = NullIntAttribute{AttributeImpl{AttrDefinition: attrDefinition}}
		} else {
			a = IntAttribute{AttributeImpl{AttrDefinition: attrDefinition}}
		}

	case schema.AttributeTypeTime:
		if attrDefinition.Nullable {
			attrDefinition.Typ = schema.AttributeTypeNullableTime
			a = NullTimeAttribute{AttributeImpl{AttrDefinition: attrDefinition}}
		} else {
			a = TimeAttribute{AttributeImpl{AttrDefinition: attrDefinition}}
		}
	default:
		panic(fmt.Errorf("unsupported attribute type %s", attrDefinition.Typ))
	}

	return a, nil
}
