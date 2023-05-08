package attribute

import "fmt"

// TimeAttribute implementation of
type TimeAttribute struct {
	AttributeImpl
}

func (a TimeAttribute) GoPackageImports() []string {
	return nil
}

func (a TimeAttribute) GoSampleValue(pkg string) string {
	if a.GetDefinition().SampleValue != "" {
		return a.GetDefinition().SampleValue
	}
	return "12"
}

// NullTimeAttribute implementation of
type NullTimeAttribute struct {
	AttributeImpl
}

func (a NullTimeAttribute) GoPackageImports() []string {
	return []string{"database/sql"}
}

func (a NullTimeAttribute) GoSampleValue(pkg string) string {
	if a.GetDefinition().SampleValue != "" {
		return fmt.Sprintf("sqlutil.ToSqlNullTime(%s)", a.GetDefinition().SampleValue)
	}
	return "sql.NullTime{}"
}
