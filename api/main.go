package api

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
)

const API_URL = "https://api.twitch.tv"

//TODO: a way to add auth token/headers to all queries

var (
	oauth2Config *clientcredentials.Config
	authToken    *oauth2.Token
)

type config struct {
	ClientID     string `mapstructure:"TWT_CLIENT_ID"`
	ClientSecret string `mapstructure:"TWT_CLIENT_SECRET"`
}

func main() {

	envTokens, err := loadConfig("..", "bot")
	if err != nil {
		log.Fatal(err)
	}

	oauth2Config = &clientcredentials.Config{
		ClientID:     envTokens.ClientID,
		ClientSecret: envTokens.ClientSecret,
		TokenURL:     twitch.Endpoint.TokenURL,
	}

	token, err := oauth2Config.Token(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	authToken = token

	fmt.Printf("Access token: %s\n", token.AccessToken)
}

func loadConfig(path string, configName string) (config config, err error) {
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
