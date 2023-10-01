package bootstrap

var Conf Config

type Config struct {
	Environment string   `mapstructure:"environment"`
	Http        Http     `mapstructrue:"http"`
	Rpc         Rpc      `mapstructure:"rpc"`
	Discover    Discover `mapstructure:"discover"`
}

type Http struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Rpc struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Discover struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	ServiceName string `mapstructure:"serviceName"`
	InstanceId  string `mapstructure:"instanceId"`
	Weight      int    `mapstructure:"weight"`
}
