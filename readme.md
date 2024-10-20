# Discord Webhook Library

A simple, flexible Go library for creating and sending Discord webhooks with rich embeds. This library supports
customizable content, username, avatar URL, and multiple embeds with fields, authors, footers, images, and colors.

![Discord Webhook](https://i.imgur.com/dOSEShX.png)
## Features

- 🌈 Customizable embed messages with rich content
- 💬 Support for multiple fields, footers, authors, images, and colors
- 🎨 Color support via Hex, RGB, or integers
- 📅 ISO8601 timestamp validation

## Installation

To install the library, use:

```
go get github.com/dozerokz/discord-webhook-go
```

## Setup

You must first configure a Webhook on a Discord server before you can use this package. Instructions can be found
on [Discord's support website](https://support.discord.com/hc/en-us/articles/228383668).

You can read more about webhooks structure [here](https://discord.com/developers/docs/resources/webhook).

## Quick Start

Here's a simple example to get you started:

```
import discordWebhook "github.com/dozerokz/discord-webhook-go"

...

// Create a webhook
webhook, err := discordWebhook.CreateWebhook("Hello, Discord!", "Bot", "SOME IMAGE URL") // replace with actual image url (string)
if err != nil {
    log.Fatal(err)
}

// Create an embed
embed, err := discordWebhook.CreateEmbed("Title", "Description", "https://example.com", "#ff5733")
if err != nil {
    log.Fatal(err)
}

// Add the embed to the webhook
discordWebhook.AddEmbed(embed)

// Send the webhook
err = discordWebhook.SendWebhook("YOUR_DISCORD_WEBHOOK_URL") // replace with your actual webhook url (string)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Webhook sent successfully!")
```

For more detailed examples, check out the [examples](examples) folder.

## License

This project is open-source. You can use, modify, and distribute it under the [MIT License](LICENSE).
