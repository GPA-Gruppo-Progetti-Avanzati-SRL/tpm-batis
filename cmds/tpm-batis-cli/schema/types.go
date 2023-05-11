package schema

type AttributeType string
type DbType string

const (
	AttributeTypeString         AttributeType = "string"
	AttributeTypeInt            AttributeType = "int"
	AttributeTypeBool           AttributeType = "bool"
	AttributeTypeTime           AttributeType = "time"
	AttributeTypeNullableString AttributeType = "nullable-string"
	AttributeTypeNullableInt    AttributeType = "nullable-int"
	AttributeTypeNullableBool   AttributeType = "nullable-bool"
	AttributeTypeNullableTime   AttributeType = "nullable-time"

	DbTypeTable DbType = "table"
	DbTypeView  DbType = "view"
)

type Field struct {
	Name          string        `json:"name,omitempty" yaml:"name,omitempty"`
	DbName        string        `json:"db-name,omitempty" yaml:"db-name,omitempty"`
	Typ           AttributeType `json:"type,omitempty" yaml:"type,omitempty"`
	Tags          []string      `json:"tags,omitempty" yaml:"tags,omitempty"`
	IsPKey        bool          `json:"primary-key,omitempty" yaml:"primary-key,omitempty"`
	WithCriterion bool          `json:"with-criterion,omitempty" yaml:"with-criterion,omitempty"`
	WithUpdate    bool          `json:"with-update,omitempty" yaml:"with-update,omitempty"`
	Nullable      bool          `json:"nullable,omitempty" yaml:"nullable,omitempty"`
	MaxLength     int           `json:"max-length,omitempty" yaml:"max-length,omitempty"`
	Options       string        `json:"options,omitempty" yaml:"options,omitempty"`
	SampleValue   string        `json:"sample-value,omitempty" yaml:"sample-value,omitempty"`
}

type Properties struct {
	DbName     string `json:"db-name,omitempty" yaml:"db-name,omitempty"`
	DbType     DbType `json:"db-type,omitempty" yaml:"db-type,omitempty"`
	FolderPath string `json:"folder-path,omitempty" yaml:"folder-path,omitempty"`
	Package    string `json:"package,omitempty" yaml:"package,omitempty"`
	StructName string `json:"struct-name,omitempty" yaml:"struct-name,omitempty"`
}

type Schema struct {
	Name        string     `json:"name,omitempty" yaml:"name,omitempty"`
	DefFileName string     `json:"-" yaml:"-"`
	Properties  Properties `json:"properties,omitempty" yaml:"properties,omitempty"`
	Fields      []Field    `json:"fields,omitempty" yaml:"fields,omitempty"`
	DDL         string     `json:"ddl,omitempty" yaml:"ddl,omitempty"`
}
