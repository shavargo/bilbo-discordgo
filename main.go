package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
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
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Close the Discord session
	dg.Close()
}

// messageCreate is called whenever a new message is created in a channel the bot has access to
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Don't proceed if the message author is a bot
	if m.Author.Bot {
		return
	}
	// Initialize a variable to hold the content of the message
	var content string

	// Check if the message contains "/chat" or "!chat"
	if strings.Contains(m.Content, "!chat") {
		content = strings.TrimPrefix(m.Content, "!chat")   // Remove the command prefix
		s.ChannelMessageSend(m.ChannelID, aiChat(content)) // Call the aiChat function to get a response from the AI
		return
	}

	// Check if the message contains "/chat" or "!chat"
	if strings.Contains(m.Content, "!image") {
		content = strings.TrimPrefix(m.Content, "!image")   // Remove the command prefix
		s.ChannelMessageSend(m.ChannelID, aiImage(content)) // Call the aiImage function to generate an image
		return
	}

}

func aiImage(x string) string {
	// Create a new OpenAI client using the OPENAI_TOKEN environment variable from Railway
	client := openai.NewClient(os.Getenv("OPENAI_TOKEN"))
	resp, err := client.CreateImage(context.Background(), openai.ImageRequest{
		Prompt: x,
		N:      1,
		Size:   openai.CreateImageSize256x256,
	})
	// Handle any errors that occur during the request
	if err != nil {
		fmt.Printf("Image generation error: %v\n", err)
		return "Sorry, there was an error generating the image."
	}

	fmt.Println(resp.Data[0].URL)
	return resp.Data[0].URL
}

// aiChat sends a message to the OpenAI API and returns the response
// It uses the OpenAI client to create a chat completion request with the provided message.
func aiChat(x string) string {
	// Create a new OpenAI client using the OPENAI_TOKEN environment variable from Railway
	client := openai.NewClient(os.Getenv("OPENAI_TOKEN"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: x,
				},
			},
		},
	)
	// Handle any errors that occur during the request
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "Sorry, there was an error processing your request."
	}

	fmt.Println(resp.Choices[0].Message.Content)
	return resp.Choices[0].Message.Content

}
