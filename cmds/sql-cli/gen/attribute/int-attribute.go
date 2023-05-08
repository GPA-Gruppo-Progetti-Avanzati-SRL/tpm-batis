package attribute

import "fmt"

// IntAttribute implementation of
type IntAttribute struct {
	AttributeImpl
}

func (a IntAttribute) GoPackageImports() []string {
	return nil
}

func (a IntAttribute) GoSampleValue(pkg string) string {
	if a.GetDefinition().SampleValue != "" {
		return a.GetDefinition().SampleValue
	}
	return "12"
}

// NullIntAttribute implementation of
type NullIntAttribute struct {
	AttributeImpl
}

func (a NullIntAttribute) GoPackageImports() []string {
	return []string{"database/sql"}
}

func (a NullIntAttribute) GoSampleValue(pkg string) string {
	if a.GetDefinition().SampleValue != "" {
		return fmt.Sprintf("sqlutil.ToSqlNullInt32(%s)", a.GetDefinition().SampleValue)
	}
	return "sql.NullInt32{}"
}
