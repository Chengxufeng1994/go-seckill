package bootstrap

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	var err error
	Logger = log.New(os.Stdout, "", log.LstdFlags)
	Logger.Println("[pkg] bootstrap init.")
	err = initBootstrapConfig()
	if err != nil {
		Logger.Printf("fail to init boostrap config: %v\n", err)
		panic(0)
	}
}

func initBootstrapConfig() error {
	viper.AddConfigPath("./")
	viper.SetConfigName("bootstrap")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	viper.OnConfigChange(func(e fsnotify.Event) {
		Logger.Println("[bootstrap] config file changed:", e.Name)
		if err := viper.Unmarshal(&Conf); err != nil {
			return
		}
	})
	viper.WatchConfig()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		return err
	}

	return nil
}
