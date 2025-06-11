// This module is a simple Discord bot that interacts with OpenAI's API.
// It uses the DiscordGo library for Discord interactions and the OpenAI client for API calls.
module discord-bot-openai

go 1.21 // This specifies the Go version used for this module.

require github.com/bwmarrin/discordgo v0.23.0 // This is the DiscordGo library version used.
require github.com/sashabaranov/go-openai v1.3.0 // This is the OpenAI client library version used.