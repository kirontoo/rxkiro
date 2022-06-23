package config

import "github.com/spf13/viper"

type Config struct {
	BotName    string `mapstructure:"BOT_NAME"`
	Streamer   string `mapstructure:"STREAMER"`
	AuthToken  string `mapstructure:"AUTH_TOKEN"`
	DbUrl      string `mapstructure:"DB_API_URL"`
	DbToken    string `mapstructure:"DB_TOKEN"`
	DbHost     string `mapstructure:"DB_HOST"`
	DbUser     string `mapstructure:"DB_USER"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	DbPort     string `mapstructure:"DB_PORT"`
	DbName     string `mapstructure:"DB_NAME"`
}

func LoadConfig(path string, configName string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(configName)
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
