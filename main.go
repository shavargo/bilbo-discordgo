package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	openai "github.com/sashabaranov/go-gpt3"
)

var token = os.Getenv("DISCORD_TOKEN")

var aitoken = os.Getenv("OPENAI_TOKEN")

var GuildID = os.Getenv("DISCORD_GUILD_ID")
var RemoveCommands = flag.Bool("remove-commands", false, "Remove commands on shutdown")

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "chat",
			Description: "Chat with Bilbo",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "prompt",
					Description: "Your message to Bilbo",
					Required:    true,
				},
			},
		},
		{
			Name:        "image",
			Description: "Generate an image with Bilbo",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "prompt",
					Description: "Your image prompt to Bilbo",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"chat": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var prompt string
			for _, opt := range i.ApplicationCommandData().Options {
				if opt.Name == "prompt" {
					prompt = opt.StringValue()
					break
				}
			}
			response := aiChat(prompt)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
				},
			})
		},
		"image": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var prompt string
			for _, opt := range i.ApplicationCommandData().Options {
				if opt.Name == "prompt" {
					prompt = opt.StringValue()
					break
				}
			}
			response := aiImage(prompt)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
				},
			})
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, 0)
	guilds, err := s.UserGuilds(100, "", "", false)
	if err != nil {
		log.Fatalf("Cannot fetch guilds: %v", err)
	}
	for _, guild := range guilds {
		// Only register commands if they haven't been registered yet for this guild
		if len(registeredCommands) == 0 {
			for _, v := range commands {
				cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guild.ID, v)
				if err != nil {
					log.Panicf("Cannot create '%v' command in guild '%v': %v", v.Name, guild.ID, err)
				}
				registeredCommands = append(registeredCommands, cmd)
			}
		}
		s.UpdateStatusComplex(discordgo.UpdateStatusData{
			Status: "online",
			Activities: []*discordgo.Activity{
				{
					Name: "Riddles",
					Type: discordgo.ActivityTypeGame,
				},
			},
		})
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
	if *RemoveCommands {
		log.Println("Removing commands...")
		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}

// aiImage generates an image based on the provided prompt using the OpenAI API
func aiImage(x string) string {
	// Create a new OpenAI client using the OPENAI_TOKEN environment variable from Railway
	client := openai.NewClient(aitoken)
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
	client := openai.NewClient(aitoken)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "system",
					Content: "You are Bilbo Baggins from The Lord of the Rings. Respond as Bilbo: witty, kind, and with a touch of Hobbit wisdom. Reference your adventures and Hobbiton when appropriate.",
				},
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
