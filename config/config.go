package config

import "github.com/spf13/viper"

type Config struct {
	BotName   string `mapstructure:"BOT_NAME"`
	Streamer  string `mapstructure:"STREAMER"`
	AuthToken string `mapstructure:"AUTH_TOKEN"`
	CmdPrefix string `mapstructure:"COMMAND_PREFIX"`
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
