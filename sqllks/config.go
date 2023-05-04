package sqllks

type Config struct {
	ServerName      string `mapstructure:"server-name" json:"server-name" yaml:"server-name"`
	ServerType      string `mapstructure:"type" json:"type" yaml:"type"`
	Host            string `mapstructure:"host" json:"host" yaml:"host"`
	Port            int    `mapstructure:"port" json:"port" yaml:"port"`
	DbName          string `mapstructure:"db-name" json:"db-name" yaml:"db-name"`
	UserName        string `mapstructure:"username" json:"username" yaml:"username"`
	Password        string `mapstructure:"password" json:"password" yaml:"password"`
	SslMode         bool   `mapstructure:"ssl-mode"  json:"ssl-mode" yaml:"ssl-mode"`
	EnableMigration bool   `mapstructure:"enable-migration" json:"enable-migration" yaml:"enable-migration"`
	MaxOpenConns    int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"`
	MaxIdleConns    int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"`
	ConnMaxLifetime int    `mapstructure:"conn-max-lifetime" json:"conn-max-lifetime" yaml:"conn-max-lifetime"`
	ConnMaxIdleTime int    `mapstructure:"conn-max-idle-time" json:"conn-max-idle-time" yaml:"conn-max-idle-time"`
}
