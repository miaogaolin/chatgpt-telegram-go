package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

const telegramEndpoint = "https://api.telegram.org/bot"

var config *Config

func main() {
	router := gin.Default()
	var err error
	config, err = ReadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	router.POST("/chat", handleMessage)
	router.Run(config.Port)
}

func handleMessage(c *gin.Context) {
	var msg Message
	c.BindJSON(&msg)

	response, err := getResponse(msg.Text)
	if err != nil {
		log.Println(err)
	}
	sendMessage(msg.ChatID, response)
}

type Message struct {
	MessageID int    `json:"message_id"`
	Text      string `json:"text"`
	ChatID    int    `json:"chat_id"`
}

type SendMessageRequest struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

func sendMessage(chatID int, message string) {
	sendMessageRequest := &SendMessageRequest{
		ChatID: chatID,
		Text:   message,
	}

	requestBody, _ := json.Marshal(sendMessageRequest)
	http.Post(telegramEndpoint+config.BotToken+"/sendMessage", "application/json", bytes.NewBuffer(requestBody))
}

func getResponse(query string) (string, error) {
	client := openai.NewClient(config.OpenApiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: query,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
