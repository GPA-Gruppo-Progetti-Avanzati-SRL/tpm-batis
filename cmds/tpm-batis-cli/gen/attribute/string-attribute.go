package attribute

import (
	"fmt"
)

// StringAttribute implementation of
type StringAttribute struct {
	AttributeImpl
}

func (a StringAttribute) GoPackageImports() []string {
	return nil
}

func (a StringAttribute) GoSampleValue(pkg string) string {
	if a.GetDefinition().SampleValue != "" {
		return fmt.Sprintf("%s.MustToMax%dText(\"%s\")", pkg, a.GetDefinition().MaxLength, a.GetDefinition().SampleValue)
	}
	return fmt.Sprintf("%s.MustToMax%dText(\"%s\")", pkg, a.GetDefinition().MaxLength, "hello")
}

// NullStringAttribute implementation of
type NullStringAttribute struct {
	AttributeImpl
}

func (a NullStringAttribute) GoPackageImports() []string {
	return []string{"database/sql"}
}

func (a NullStringAttribute) GoSampleValue(pkg string) string {
	if a.GetDefinition().SampleValue != "" {
		return fmt.Sprintf("sql.NullString(%s.MustToNullableMax%dText(\"%s\"))", pkg, a.GetDefinition().MaxLength, a.GetDefinition().SampleValue)
	}
	return "sql.NullString{}"
}
