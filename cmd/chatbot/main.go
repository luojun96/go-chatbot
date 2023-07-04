package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/luojun96/chatgpt-opt/conf"
	openai "github.com/sashabaranov/go-openai"
)

func getEnvOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func interactGPT(content string, o *conf.OpenAI) (string, error) {
	client := openai.NewClient(o.Token)

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
	yamlFile := flag.String("f", getEnvOrDefault("CHATBOT_CONFIG", "config.yaml"), "configuration file")
	flag.Parse()

	c, err := conf.New(yamlFile)
	if err != nil {
		log.Fatalf("Reading the YAML configuration file failed: %v", err)
	}

	fmt.Println("=======Welcome to GPT chatbot, type something to start:=======")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Human> ")
	for scanner.Scan() {
		input := scanner.Text()
		if input == "exit" {
			break
		}
		output, _ := interactGPT(input, &c.OpenAI)

		fmt.Println("-------------------------------------------------------------")
		fmt.Printf("Bot> %s\n", output)
		fmt.Println("==============Type something to continue:====================")
		fmt.Print("> ")
	}

	fmt.Printf("Bye!\n")
}
