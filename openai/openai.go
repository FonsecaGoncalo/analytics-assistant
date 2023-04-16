package openai

import (
	"analytics/config"
	"context"
	"github.com/sashabaranov/go-openai"
	"strings"
)

type Session struct {
	temperature float32
	messages    []openai.ChatCompletionMessage
}

var (
	openAiClient = openai.NewClient(config.GetAPIKey())
)

func NewOpenAISession(systemMessages []string, temperature float32) *Session {
	messages := make([]openai.ChatCompletionMessage, len(systemMessages))

	for i := range systemMessages {
		messages[i] = openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemMessages[i],
		}
	}

	return &Session{
		temperature: temperature,
		messages:    messages,
	}
}

func (s *Session) SystemPrompt(prompt string) string {
	return s.prompt(prompt, openai.ChatMessageRoleSystem)
}

func (s *Session) UserPrompt(prompt string) string {
	return s.prompt(prompt, openai.ChatMessageRoleUser)
}

func (s *Session) prompt(prompt string, role string) string {
	s.messages = append(s.messages, openai.ChatCompletionMessage{
		Role:    role,
		Content: prompt,
	})

	resp, err := openAiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Messages:    s.messages,
			Temperature: s.temperature,
		},
	)

	if err != nil {
		panic("Error prompting open ai")
	}

	if len(resp.Choices) == 0 {
		panic("no choices returned from the OpenAI API")
	}

	s.messages = append(s.messages, resp.Choices[0].Message)

	return strings.TrimSpace(resp.Choices[0].Message.Content)
}
