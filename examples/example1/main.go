// Example 1: Simplest webhook that sends a message to a Discord channel
package main

import (
	discordWebhook "github.com/dozerokz/discord-webhook-go"
	"log"
)

func main() {
	// Discord webhook URL obtained from your Discord server settings
	webhookURL := "YOUR_WEBHOOK_URL"

	// URL of the image you want to display as the avatar for the bot's message (can be empty string)
	imageURL := "https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Aqua.png"

	// Create a new webhook payload with a message, bot's username, and an avatar image
	webhookPayload, err := discordWebhook.CreateWebhook("Hello Discord From Go!",
		"Bot's username", imageURL) // (username also can be empty string)
	if err != nil {
		// Log and exit if there’s an error during payload creation
		log.Fatalf("Error creating webhook payload: %v", err)
	}

	// Send the created webhook payload to the specified Discord webhook URL
	err = discordWebhook.SendWebhook(webhookURL, webhookPayload)
	if err != nil {
		// Log and exit if there’s an error while sending the webhook
		log.Fatalf("Error sending webhook: %v", err)
	}
}
