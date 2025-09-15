// This module is a simple Discord bot that interacts with OpenAI's API.
// It uses the DiscordGo library for Discord interactions and the OpenAI client for API calls.
module discord-bot-openai

go 1.24.0 // This specifies the Go version used for this module.

toolchain go1.24.2

require github.com/bwmarrin/discordgo v0.29.0 // This is the DiscordGo library version used.

require github.com/sashabaranov/go-gpt3 v1.4.0 // This is the OpenAI client library version used.

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/crypto v0.42.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
)
