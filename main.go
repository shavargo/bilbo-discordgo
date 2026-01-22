package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	openai "github.com/sashabaranov/go-gpt3"
)

var (
	token   = os.Getenv("DISCORD_TOKEN")
	aiToken = os.Getenv("OPENAI_TOKEN")
	session *discordgo.Session
)

var commands = []*discordgo.ApplicationCommand{
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

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"chat": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		prompt := getOptionValue(i, "prompt")
		response := aiChat(prompt)
		respondToInteraction(s, i, response)
	},
	"image": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		prompt := getOptionValue(i, "prompt")
		response := aiImage(prompt)
		respondToInteraction(s, i, response)
	},
}

func init() {
	var err error
	session, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

// getOptionValue extracts a string option value from an interaction by name.
func getOptionValue(i *discordgo.InteractionCreate, name string) string {
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == name {
			return opt.StringValue()
		}
	}
	return ""
}

// respondToInteraction sends a response message to a Discord interaction.
func respondToInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

func main() {
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	if err := session.Open(); err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer session.Close()

	guilds, err := session.UserGuilds(100, "", "", false)
	if err != nil {
		log.Fatalf("Cannot fetch guilds: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := registerCommands(guilds)

	session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Status: "online",
		Activities: []*discordgo.Activity{
			{
				Name: "Riddles",
				Type: discordgo.ActivityTypeGame,
			},
		},
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Removing commands...")
	unregisterCommands(guilds, registeredCommands)

	log.Println("Gracefully shutting down.")
}

// registerCommands registers all slash commands to each guild the bot is in.
func registerCommands(guilds []*discordgo.UserGuild) []*discordgo.ApplicationCommand {
	registeredCommands := make([]*discordgo.ApplicationCommand, 0)
	for _, guild := range guilds {
		registered, err := session.ApplicationCommandBulkOverwrite(session.State.User.ID, guild.ID, commands)
		if err != nil {
			log.Printf("Cannot register commands in guild '%v': %v", guild.ID, err)
			continue
		}
		registeredCommands = append(registeredCommands, registered...)
	}
	return registeredCommands
}

// unregisterCommands removes all registered commands from each guild.
func unregisterCommands(guilds []*discordgo.UserGuild, registeredCommands []*discordgo.ApplicationCommand) {
	for _, guild := range guilds {
		log.Printf("Removing commands from guild: %v", guild.ID)
		for _, cmd := range registeredCommands {
			if err := session.ApplicationCommandDelete(session.State.User.ID, guild.ID, cmd.ID); err != nil {
				log.Printf("Cannot delete '%v' command from guild '%v': %v", cmd.Name, guild.ID, err)
			}
		}
	}
}

// aiImage generates an image based on the provided prompt using the OpenAI API.
func aiImage(prompt string) string {
	client := openai.NewClient(aiToken)
	resp, err := client.CreateImage(context.Background(), openai.ImageRequest{
		Prompt: prompt,
		N:      1,
		Size:   openai.CreateImageSize256x256,
	})
	if err != nil {
		log.Printf("Image generation error: %v", err)
		return "Sorry, there was an error generating the image."
	}

	return resp.Data[0].URL
}

// aiChat sends a message to the OpenAI API and returns the response.
func aiChat(prompt string) string {
	client := openai.NewClient(aiToken)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: "system",
					Content: `You are Bilbo Baggins, the famous Hobbit of the Shire. Respond in character with wit, warmth, and Hobbit wisdom.

Your history and adventures:
- You live at Bag End in Hobbiton, a cozy hobbit-hole built by your father Bungo Baggins
- Gandalf the Grey recruited you as the burglar for Thorin Oakenshield's company of 13 dwarves
- You found the One Ring in Gollum's cave beneath the Misty Mountains and won the riddle game against him
- You encountered the trolls Tom, Bert, and William, and helped delay them until sunrise turned them to stone
- You escaped goblins in the Misty Mountains and were rescued by eagles
- You stayed with Beorn the skin-changer who could transform into a great bear
- You saved the dwarves from giant spiders in Mirkwood using your sword Sting
- You freed the dwarves from the Wood-elves' dungeon using the Ring's invisibility and escaped in barrels
- You spoke with Smaug the dragon in Erebor, discovering the weakness in his armor
- You witnessed the Battle of Five Armies and Thorin's death
- You possess the mithril coat given to you by Thorin, and your sword Sting that glows blue near orcs
- At your 111th birthday party, you disappeared using the Ring and left the Shire
- You passed the Ring to your nephew Frodo and retired to Rivendell to write your book "There and Back Again"
- You eventually sailed to the Undying Lands with Gandalf, Frodo, and the Elves

Your personality:
- You love comfort: good food, your armchair, pipe-weed, and a warm fire
- You never miss breakfast, second breakfast, elevenses, luncheon, afternoon tea, dinner, or supper
- Your Took side craves adventure while your Baggins side yearns for home
- You write poetry and songs, and enjoy riddles
- You are clever, resourceful, and braver than you appear
- You speak fondly of the Shire and Hobbit customs

Speak as Bilbo would: polite, occasionally flustered, fond of expressions like "Good gracious!" and references to meals and home comforts.`,
				},
				{
					Role:    "user",
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		log.Printf("ChatCompletion error: %v", err)
		return "Sorry, there was an error processing your request."
	}

	return resp.Choices[0].Message.Content
}
