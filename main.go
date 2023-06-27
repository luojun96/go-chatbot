package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
	yaml "gopkg.in/yaml.v3"
)

// Token configure service token for GPT
type Token struct {
	Chatgpt3 string `yaml:"chatgpt3"`
}

const configPath = "./config.yaml"

var token Token

// read token from config file `config.yaml`
func fetchToken() error {
	log.Print("Start to fetch token from configure file.")
	file, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("reading configure file failed: %v", err)
	}

	if err = yaml.Unmarshal(file, &token); err != nil {
		return fmt.Errorf("unmarshaling token failed: %v", err)
	}

	log.Print("Fetch token from configure file successfully.")
	if token.Chatgpt3 == "" {
		return errors.New("empty token")
	}

	return nil
}

func interactGPT(content string) (string, error) {
	client := openai.NewClient(token.Chatgpt3)

	resp, err := client.CreateChatCompletion(context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("failed creating chat completion: %+v", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func main() {
	if err := fetchToken(); err != nil {
		log.Fatalf("Reading configure file failed: %v", err)
	}

	fmt.Println("=======Welcome to GPT chatbot, type something to start:=======")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Human> ")
	for scanner.Scan() {
		input := scanner.Text()
		if input == "exit" {
			break
		}
		output, _ := interactGPT(input)

		fmt.Println("-------------------------------------------------------------")
		fmt.Printf("Bot> %s\n", output)
		fmt.Println("==============Type something to continue:====================")
		fmt.Print("> ")
	}

	fmt.Printf("Bye!\n")
}
