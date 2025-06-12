package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	openai "github.com/sashabaranov/go-gpt3"
)

func main() {
	// Create a new session using the DISCORD_TOKEN environment variable from Railway
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		fmt.Printf("Error while starting bot: %s", err)
		return
	}

	// Add the message handler
	dg.AddHandler(messageCreate)

	// Add the Guild Messages intent
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Connect to the gateway
	err = dg.Open()
	if err != nil {
		fmt.Printf("Error while connecting to gateway: %s", err)
		return
	}

	// Wait until Ctrl+C or another signal is received
	fmt.Println("The bot is now running. Press Ctrl+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Close the Discord session
	dg.Close()
}

func aiChat(x string) {
	client := openai.NewClient(os.Getenv("OPENAI_TOKEN"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: x,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
	return resp.Choices[0].Message.Content

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Don't proceed if the message author is a bot
	if m.Author.Bot {
		return
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, aiChat(m.Content))
		return
	}

	if m.Content == "hello" {
		s.ChannelMessageSend(m.ChannelID, aiChat(m.Content))
		return
	}
}
