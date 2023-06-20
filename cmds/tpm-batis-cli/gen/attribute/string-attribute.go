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
		return fmt.Sprintf("%s.MustValidate%s(\"%s\")", pkg, a.GoAttributeName(), a.GetDefinition().SampleValue)
	}
	return fmt.Sprintf("%s.MustValidate%s(\"%s\")", pkg, a.GoAttributeName(), "hello")
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
		return fmt.Sprintf("%s.MustValidate%s(\"%s\")", pkg, a.GoAttributeName(), a.GetDefinition().SampleValue)
	}
	return "sql.NullString{}"
}
