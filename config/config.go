package config

type Config struct {
	Database Database `mapstructure:",squash"`
	Server   Server   `mapstructure:",squash"`
	Settings Settings `mapstructure:",squash"`
	Redis    Redis    `mapstructure:",squash"`
}

type Database struct {
	Address         string `mapstructure:"DATABASE_ADDRESS"`
	Port            string `mapstructure:"DATABASE_PORT"`
	Username        string `mapstructure:"DATABASE_USERNAME"`
	Password        string `mapstructure:"DATABASE_PASSWORD"`
	DBName          string `mapstructure:"DATABASE_DBNAME"`
	MaxOpenConn     int    `mapstructure:"DATABASE_MAXOPENCONN"`
	MaxIdleConn     int    `mapstructure:"DATABASE_MAXIDLECONN"`
	ConnMaxIdleTime int    `mapstructure:"DATABASE_CONNMAXIDLETIME"`
	ConnMaxLifeTime int    `mapstructure:"DATABASE_CONNMAXLIFETIME"`
}

type Server struct {
	Address string `mapstructure:"SERVER_ADDRESS"`
	Debug   bool   `mapstructure:"SERVER_DEBUG"`
	Name    string `mapstructure:"SERVER_NAME"`
	Timeout string `mapstructure:"SERVER_TIMEOUT"`
}

type Settings struct {
	JWTSecret            string `mapstructure:"SETTINGS_JWTSECRET"`
	TokenDuration        int    `mapstructure:"SETTINGS_TOKENDURATION"`
	RefreshTokenDuration int    `mapstructure:"SETTINGS_REFRESHTOKENDURATION"`
}

type Redis struct {
	Addr     string `mapstructure:"REDIS_ADDR"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}
