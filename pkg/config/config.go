package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	Logger = log.New(os.Stdout, "", log.LstdFlags)
	Logger.Println("[pkg] config init.")
}

func LoadLocalConfig() error {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	viper.OnConfigChange(func(e fsnotify.Event) {
		Logger.Println("[config] config file changed:", e.Name)
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
