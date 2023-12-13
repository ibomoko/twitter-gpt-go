package twitter

import (
	"fmt"
	"context"
	"path/filepath"
	"github.com/spf13/viper"
	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
)

type Credentials struct {
	accessToken string
	accessTokenSecret string
}

func Tweet(text string) {
	credentials := fetchTwitterConfig()
	in := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:          credentials.accessToken,
		OAuthTokenSecret:     credentials.accessTokenSecret,
	}

	c, err := gotwi.NewClient(in)
	if err != nil {
		fmt.Printf("New client set error:\n%s", err)
		return
	}

	p := &types.CreateInput{
		Text: gotwi.String(text),
	}

	_, err = managetweet.Create(context.Background(), c, p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func fetchTwitterConfig() *Credentials {
	configPath, err := getConfigPath()
	if err != nil {
		handleError("Error getting absoulte path", err)
		return nil
	}

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		handleError("Error reading configuration", err)
		return nil
	}
	
	return &Credentials{
		accessToken: getStringFromConfig("twitter.access-token"),
		accessTokenSecret: getStringFromConfig("twitter.access-token-secret"),
	}
}

func getConfigPath() (string, error) {
	configPath, err := filepath.Abs("./config/config.yml")
	return configPath, err
}

func handleError(message string, err error) {
	fmt.Printf("%s: %v\n", message, err)
}

func getStringFromConfig(key string) string {
	value := viper.GetString(key)
	if value == "" {
		fmt.Printf("%s property is not set in the configuration file\n", key)
	}
	return value
}
