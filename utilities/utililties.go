package utilities

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	MONGO_URI string `mapstructure:"MONGO_URI"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

// A function to check error
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
