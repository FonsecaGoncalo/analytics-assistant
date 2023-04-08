package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func GetAPIKey() string {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = viper.GetString("openai.api_key")
	}
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY env var not set")
		os.Exit(-1)
	}
	return apiKey
}

func GetAnalystSystemMessages() []string {
	return viper.GetStringSlice("openai.analyst.messages.system_messages")
}

func GetAnalystContextMessages() string {
	return viper.GetString("openai.analyst.messages.context_message")
}

func GetAnalystQueryResultsMessage() string {
	return viper.GetString("openai.analyst.messages.query_results_message")
}

func GetQueryParserSystemMessages() []string {
	return viper.GetStringSlice("openai.query_parser.messages.system_messages")
}

func GetQueryParserMessage() string {
	return viper.GetString("openai.query_parser.messages.parse_query_message")
}

func GetAnalystTemperature() float32 {
	return float32(viper.GetFloat64("openai.analyst.temperature"))
}

func GetQueryParserTemperature() float32 {
	return float32(viper.GetFloat64("openai.query_parser.temperature"))
}
