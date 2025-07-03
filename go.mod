// This module is a simple Discord bot that interacts with OpenAI's API.
// It uses the DiscordGo library for Discord interactions and the OpenAI client for API calls.
module discord-bot-openai

go 1.21 // This specifies the Go version used for this module.

require github.com/bwmarrin/discordgo v0.29.0 // This is the DiscordGo library version used.

require github.com/sashabaranov/go-gpt3 v1.3.0 // This is the OpenAI client library version used.

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
)
