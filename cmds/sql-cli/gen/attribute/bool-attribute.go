package attribute

import "fmt"

// BoolAttribute implementation of
type BoolAttribute struct {
	AttributeImpl
}

func (a BoolAttribute) GoPackageImports() []string {
	return nil
}

func (a BoolAttribute) GoSampleValue(pkg string) string {
	if a.GetDefinition().SampleValue != "" {
		return a.GetDefinition().SampleValue
	}
	return "false"
}

// NullBoolAttribute implementation of
type NullBoolAttribute struct {
	AttributeImpl
}

func (a NullBoolAttribute) GoPackageImports() []string {
	return []string{"database/sql"}
}

func (a NullBoolAttribute) GoSampleValue(pkg string) string {
	if a.GetDefinition().SampleValue != "" {
		return fmt.Sprintf("sqlutil.ToSqlNullBool(%s)", a.GetDefinition().SampleValue)
	}
	return "sql.NullBool{}"
}
