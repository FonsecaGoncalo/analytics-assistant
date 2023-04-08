package openai

import (
	"analytics/config"
	"context"
	"fmt"
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

func (s *Session) SystemPrompt(prompt string) (string, error) {
	return s.prompt(prompt, openai.ChatMessageRoleSystem)
}

func (s *Session) UserPrompt(prompt string) (string, error) {
	return s.prompt(prompt, openai.ChatMessageRoleUser)
}

func (s *Session) prompt(prompt string, role string) (string, error) {
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
		fmt.Println("Error prompting open ai")
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from the OpenAI API")
	}

	s.messages = append(s.messages, resp.Choices[0].Message)

	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

func GenerateText(apiKey, prompt string, ctx []string) (string, error) {
	fmt.Println(prompt)

	c := openai.NewClient(apiKey)

	messages := make([]openai.ChatCompletionMessage, len(ctx))
	for i, m := range ctx {
		messages[i] = openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: "Context that may help you answer the question: " + m,
		}
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: "You are a data analyst",
	})

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})

	resp, err := c.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)

	if err != nil {
		fmt.Printf("Completion error: %v\n", err)
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from the OpenAI API")
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}
