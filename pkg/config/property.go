package config

var Conf Config

type Config struct {
	Postgres Postgres `mapstructure:"postgres"`
	Trace    Trace    `mapstructure:"trace"`
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Db       string `mapstructure:"db"`
}

type Trace struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Url  string `mapstructure:"url"`
}
