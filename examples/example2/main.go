// Example 2: More complicated webhook with author, thumbnail, image, footer and several fields.
package main

import (
	discordWebhook "github.com/dozerokz/discord-webhook-go" // Import the webhook library
	"log"
	"time"
)

func main() {
	// Discord webhook URL from your Discord server settings
	webhookURL := "YOUR_WEBHOOK_URL"

	// URL of an avatar image for the bot (optional)
	imageURL := "https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Aqua.png"

	// Create a basic webhook message with content, bot username, and avatar
	webhookPayload, err := discordWebhook.CreateWebhook("Hello Discord From Go!", "Bot's username", imageURL)
	if err != nil {
		log.Fatalf("Error creating webhook payload: %v", err)
	}

	// Create an embed with a title, description, URL, and color
	embed, err := discordWebhook.CreateEmbed("Embed Title", "Embed description text", "https://example.com", discordWebhook.RGB{R: 111, G: 222, B: 77})
	if err != nil {
		log.Fatalf("Error creating embed: %v", err)
	}

	// Create embed components: thumbnail, footer, image, author, and fields
	thumbnail := discordWebhook.CreateThumbnail(imageURL, "", 300, 300)
	footer := discordWebhook.CreateFooter("Footer Text", imageURL, "")
	image := discordWebhook.CreateImage(imageURL, "", 200, 200)
	author := discordWebhook.CreateAuthor("@Dozerokz", "", imageURL, "")

	// Create fields and add them to the embed
	field1 := discordWebhook.CreateField("Field Name 1", "Value 1", true)
	field2 := discordWebhook.CreateField("Field Name 2", "Value 2", true)
	field3 := discordWebhook.CreateField("Field Name 3", "Value 3", true)
	field4 := discordWebhook.CreateField("Field Name 4", "Value 4", false)
	fields := []discordWebhook.Field{field1, field2, field3}

	// Set the embed's thumbnail, footer, image, author, and fields
	embed.SetThumbnail(thumbnail)
	embed.SetFooter(footer)
	embed.SetImage(image)
	embed.SetAuthor(author)
	embed.AddFields(fields) // Add multiple fields
	embed.AddField(field4)  // Add an extra field

	// Set a timestamp for the embed using the current time
	err = embed.SetTimestamp(time.Now().Format(time.RFC3339))
	if err != nil {
		log.Fatalf("Error setting timestamp: %v", err)
	}

	// Add the embed to the webhook payload
	webhookPayload.AddEmbed(embed)

	// Send the webhook with the payload to the specified URL
	err = discordWebhook.SendWebhook(webhookURL, webhookPayload)
	if err != nil {
		log.Fatalf("Error sending webhook: %v", err)
	}
}
