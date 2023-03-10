package main

import (
	"fmt"
	"log"
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

// ReadConfig reads the config file and unmarshals it into the config variable
func ReadConfig() error {
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

func Start() error {
	dg, err := discordgo.New("Bot " + token.Discord)
	if err != nil {
		log.Fatalf("failed to create Discord session: %v", err)
		return err
	}
	// Register the messageCreate func as a callback for MessageCreate events.

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

func main() {
	err := ReadConfig()
	if err != nil {
		fmt.Println(err.Error())
	}

	Start()
}
