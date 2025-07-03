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
	// Register the slash command
	cmd := &discordgo.ApplicationCommand{
		Name:        "chat",
		Description: "Chat with the AI",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "prompt",
				Description: "What do you want to say to the AI?",
				Required:    true,
			},
		},
	}
	_, err = dg.ApplicationCommandCreate(os.Getenv("DISCORD_APP_ID"), os.Getenv("DISCORD_GUILD_ID"), cmd)
	if err != nil {
		fmt.Printf("Cannot create slash command: %v\n", err)
		return
	}

	dg.AddHandler(func(
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
	) {
		data := i.ApplicationCommandData()
		switch data.Name {
		case "chat":
			var prompt string
			for _, opt := range data.Options {
				if opt.Name == "prompt" {
					prompt = opt.StringValue()
				}
			}
			response := aiChat(prompt)
			err := s.InteractionRespond(
				i.Interaction,
				&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: response,
					},
				},
			)
			if err != nil {
				fmt.Printf("InteractionRespond error: %v\n", err)
			}
		}
	})

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

// aiImage generates an image based on the provided prompt using the OpenAI API
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
