package config

import "github.com/spf13/viper"

type Config struct {
	BotName   string `mapstructure:"BOT_NAME"`
	Streamer  string `mapstructure:"STREAMER"`
	AuthToken string `mapstructure:"AUTH_TOKEN"`
	CmdPrefix string `mapstructure:"COMMAND_PREFIX"`
	DbUrl     string `mapstructure:"DB_API_URL"`
	DbToken   string `mapstructure:"DB_TOKEN"`
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
