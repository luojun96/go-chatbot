package main

import (
	"context"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
	yaml "gopkg.in/yaml.v3"
)

type Token struct {
	Chatgpt3 string `yaml:"chatgpt3"`
}

var token Token

// read token from config file `config.yaml`
func readConfig() error {
	log.Print("Reading config file...")
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Fatalf("failed reading data from file: %+v", err)
		return err
	}

	if err = yaml.Unmarshal(data, &token); err != nil {
		log.Fatalf("failed unmarshaling data: %+v", err)
		return err
	}

	log.Print("Reading config file success")
	return nil
}

func interactGPT3(token string) {
	client := openai.NewClient(token)
	resp, err := client.CreateChatCompletion(context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello, I'm a human",
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("failed creating chat completion: %+v", err)
	}

	log.Printf("response: %+v", resp)
}

func main() {
	// read config file
	if err := readConfig(); err != nil {
		log.Fatalf("failed reading config file: %+v", err)
	}

	// interact with GPT-3
	interactGPT3(token.Chatgpt3)
}
