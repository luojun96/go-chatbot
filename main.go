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

type Token struct {
	Chatgpt3 string `yaml:"chatgpt3"`
}

var token Token

// read token from config file `config.yaml`
func readConfig() error {
	log.Print("Reading config file...")
	file, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Fatalf("failed reading data from file: %+v", err)
		return err
	}

	if err = yaml.Unmarshal(file, &token); err != nil {
		log.Fatalf("failed unmarshaling data: %+v", err)
		return err
	}

	log.Print("Reading config file success")
	if token.Chatgpt3 == "" {
		return errors.New("token is empty")
	}

	return nil
}

func interactGPT3(content string) (string, error) {
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
	// read config file
	if err := readConfig(); err != nil {
		log.Fatalf("failed reading config file: %+v", err)
		return
	}

	fmt.Println("=======Welcome to GPT chatbot, type something to start:=======")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Human> ")
	for scanner.Scan() {
		input := scanner.Text()
		if input == "exit" {
			break
		}
		output, _ := interactGPT3(input)

		fmt.Println("-------------------------------------------------------------")
		fmt.Printf("Bot> %s\n", output)
		fmt.Println("==============Type something to continue:====================")
		fmt.Print("> ")
	}

	fmt.Printf("Bye!\n")
}
