package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/spf13/viper"
)



type ChatCompletion struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
}

type Message struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
}

type ChatCompletionResponse struct {
	Choices   []Choice `json:"choices"`
}

type Choice struct {
	Message   Message `json:"message"`
}

func CompleteChat(prompt string)  {
	apiKey, url, model := fetchOpenAIConfig()

	ChatCompletion := ChatCompletion{
		Model: model,
		Messages: []Message{{Role: "user", Content: prompt}},
	}
	
	requestBody, err := json.Marshal(ChatCompletion)
	if err != nil {
		fmt.Println("Error marshalling json: ", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + apiKey)

	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	
	defer resp.Body.Close()

	var chatCompletionResponse ChatCompletionResponse
	err = json.NewDecoder(resp.Body).Decode(&chatCompletionResponse)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return
	}

	fmt.Println(chatCompletionResponse.Choices[0].Message.Content)

}


func fetchOpenAIConfig() (apiKey, url, model string) {
	configPath, err := getConfigPath()
	if err != nil {
		handleError("Error getting absolute path", err)
		return "", "", ""
	}

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		handleError("Error reading configuration", err)
		return "", "", ""
	}

	apiKey = getStringFromConfig("openai.api-key")
	url = getStringFromConfig("openai.url")
	model = getStringFromConfig("openai.model")

	return apiKey, url, model
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
		fmt.Printf("%s property is not set in the configuration\n", key)
	}
	return value
}