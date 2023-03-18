package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	yaml "gopkg.in/yaml.v3"
)

// Token is the token for the discord bot and chatgpt
type Token struct {
	Discord string `yaml:"discord"`
	ChatGPT string `yaml:"chat_gpt"`
}

var token Token

// readConfig reads the config file and unmarshals it into the config variable
func readConfig() error {
	fmt.Println("Reading config file...")
	file, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
		return err
	}

	err = yaml.Unmarshal(file, &token)
	if err != nil {
		log.Fatalf("failed to decode config file: %v", err)
	}

	fmt.Println("Config file read successfully")
	return nil
}

func start() error {
	dg, err := discordgo.New("Bot " + token.Discord)
	if err != nil {
		log.Fatalf("failed to create Discord session: %v", err)
		return err
	}
	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageHandler)
	dg.Identify.Intents = discordgo.IntentGuildMessages
	err = dg.Open()
	if err != nil {
		log.Fatalf("failed to open connection: %v", err)
		return err
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	select {
	case <-done:
		fmt.Println("received the exit signal, exiting...")
	}
	fmt.Println("closing Discord session...")
	dg.Close()
	return nil
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Ignore all messages that don't mention the bot
	mentioned := false
	for _, u := range m.Mentions {
		if u.ID == s.State.User.ID {
			mentioned = true
			break
		}
	}
	if !mentioned {
		fmt.Printf("Not mentioned in message: [%s] %s\n", m.Author.ID, m.Content)
		return
	}

	message := fmt.Sprintf("Message: %s, Author: %s", m.Content, m.Author.ID)

	fmt.Println(message)
	chatGPTResponse, err := callChatGPT(m.Content)
	if err != nil {
		fmt.Println(err.Error())
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	s.ChannelMessageSendReply(m.ChannelID, chatGPTResponse, m.Reference())
}

type ChatGPTResponse struct {
	ID     string `json:"id"`
	Object string `json:"object"`
	Model  string `json:"model"`
	Usage  struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

func callChatGPT(msg string) (string, error) {
	api := "https://api.openai.com/v1/chat/completions"
	body := []byte(`{
		"model": "gpt-3.5-turbo",
		"messages" : [
			{
				"role":"user",
				"content": "` + JSONEscape(msg) + `",
			}
		]
	}`)
	req, err := http.NewRequest("POST", api, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.ChatGPT)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response", err)
		return "", err
	}

	chatGPTData := ChatGPTResponse{}
	err = json.Unmarshal(body, &chatGPTData)
	if err != nil {
		fmt.Println("Error unmarshalling response", err)
		return "", err
	}
	return chatGPTData.Choices[0].Message.Content, nil
}

func JSONEscape(str string) string {
	b, err := json.Marshal(str)
	if err != nil {
		return str
	}
	s := string(b)
	return s[1 : len(s)-1]
}

func main() {
	err := readConfig()
	if err != nil {
		fmt.Println(err.Error())
	}

	start()
}
