// Package webhook is a Go library for creating and sending Discord webhooks.
// It simplifies the process of building Discord webhook payloads, including rich embed messages,
// with support for customizable content, usernames, avatars, and multiple embeds.
// The library provides a set of flexible, easy-to-use functions to handle complex message structures,
// such as embeds with fields, authors, footers, images, and colors.

package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	minDecimalRGBValue = 0
	maxDecimalRGBValue = 16777215
	minRGBValue        = 0
	maxRGBValue        = 255
)

// Webhook represents the structure for sending a message via Discord webhooks.
// It can include optional content, username, avatar URL, and an array of rich embed objects.
type Webhook struct {
	Content   string  `json:"content,omitempty"`
	Username  string  `json:"username,omitempty"`
	AvatarURL string  `json:"avatar_url,omitempty"`
	Embeds    []Embed `json:"embeds,omitempty"`
}

// Embed represents a rich embed object for Discord
type Embed struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	URL         string    `json:"url,omitempty"`
	Color       int       `json:"color,omitempty"`
	Timestamp   string    `json:"timestamp,omitempty"`
	Footer      Footer    `json:"footer,omitempty"`
	Image       Image     `json:"image,omitempty"`
	Thumbnail   Thumbnail `json:"thumbnail,omitempty"`
	Author      Author    `json:"author,omitempty"`
	Fields      []Field   `json:"fields,omitempty"`
}

// Footer represents the footer section of an embed
type Footer struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// Image represents an image in an embed
type Image struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

// Author represents the author section of an embed
type Author struct {
	Name         string `json:"name,omitempty"`
	URL          string `json:"url,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// Thumbnail represents a thumbnail image in an embed
type Thumbnail struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

// Field represents a key-value pair in the embed's fields section
type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// RGB is a struct representing an RGB color
type RGB struct {
	R, G, B int
}

// colorEmbedTypes is an interface that defines valid color types (string, int, or RGB)
type colorEmbedTypes interface {
	string | int | RGB
}

// AddEmbed adds an embed to the webhook
func (w *Webhook) AddEmbed(embed Embed) {
	w.Embeds = append(w.Embeds, embed)
}

// AddField adds a field to the embed
func (e *Embed) AddField(field Field) {
	e.Fields = append(e.Fields, field)
}

// AddFields adds multiple fields to the embed
func (e *Embed) AddFields(fields []Field) {
	e.Fields = append(e.Fields, fields...)
}

// SetFooter sets the footer for the embed
func (e *Embed) SetFooter(footer Footer) {
	e.Footer = footer
}

// SetImage sets the image for the embed
func (e *Embed) SetImage(image Image) {
	e.Image = image
}

// SetThumbnail sets the thumbnail for the embed
func (e *Embed) SetThumbnail(thumbnail Thumbnail) {
	e.Thumbnail = thumbnail
}

// SetAuthor sets the author for the embed
func (e *Embed) SetAuthor(author Author) {
	e.Author = author
}

// SetTimestamp validates and sets an ISO8601 timestamp for the embed
func (e *Embed) SetTimestamp(timestamp string) error {
	if isValidISO8601(timestamp) {
		e.Timestamp = timestamp
		return nil
	}

	return fmt.Errorf("timestamp is not ISO8601 timestamp")
}

// CreateWebhook creates a new Webhook with the specified content, username, and avatar URL
func CreateWebhook(content, username, avatarURL string) (Webhook, error) {
	if len(content) > 2000 {
		return Webhook{}, fmt.Errorf("the length of the content cannot exceed 2000 characters (your length: %d)", len(content))
	}
	return Webhook{
		Content:   content,
		Username:  username,
		AvatarURL: avatarURL,
	}, nil
}

// CreateEmbed creates an embed with title, description, URL, and color (supports string, int, or RGB)
func CreateEmbed[T colorEmbedTypes](title, description, url string, color T) (Embed, error) {
	var colorValue int
	var err error

	switch v := any(color).(type) {
	case string:
		colorValue, err = hexToColorInt(v)
	case int:
		if v > maxDecimalRGBValue || v < minDecimalRGBValue {
			err = fmt.Errorf("you can only use numbers from 0 to 16777215")
		}
		colorValue = v
	case RGB:
		colorValue, err = rgbToInt(v)
	}

	if err != nil {
		return Embed{}, err
	}

	return Embed{
		Title:       title,
		Description: description,
		URL:         url,
		Color:       colorValue,
	}, nil
}

// CreateFooter creates a footer object for an embed
func CreateFooter(text string, iconURL string, proxyIconURL string) Footer {
	footer := Footer{
		Text:         text,
		IconURL:      iconURL,
		ProxyIconURL: proxyIconURL,
	}
	return footer
}

// CreateImage creates an image object for an embed
func CreateImage(url string, proxyURL string, height int, width int) Image {
	image := Image{
		URL:      url,
		ProxyURL: proxyURL,
		Height:   height,
		Width:    width,
	}
	return image
}

// CreateThumbnail creates a thumbnail object for an embed
func CreateThumbnail(url string, proxyURL string, height int, width int) Thumbnail {
	thumbnail := Thumbnail{
		URL:      url,
		ProxyURL: proxyURL,
		Height:   height,
		Width:    width,
	}
	return thumbnail
}

// CreateAuthor creates an author object for an embed
func CreateAuthor(name string, url string, iconURL string, proxyIconURL string) Author {
	author := Author{
		Name:         name,
		URL:          url,
		IconURL:      iconURL,
		ProxyIconURL: proxyIconURL,
	}
	return author
}

// CreateField creates a field object for an embed
func CreateField(name string, value string, inline bool) Field {
	field := Field{
		Name:   name,
		Value:  value,
		Inline: inline,
	}
	return field
}

// rgbToInt converts an RGB struct to an integer representation
func rgbToInt(rgb RGB) (int, error) {
	var err error
	if rgb.R > maxRGBValue || rgb.R < minRGBValue ||
		rgb.G > maxRGBValue || rgb.G < minRGBValue ||
		rgb.B > maxRGBValue || rgb.B < minRGBValue {
		err = fmt.Errorf("rgb color can only be from 0 to 255")
	}
	return (rgb.R << 16) | (rgb.G << 8) | rgb.B, err
}

// hexToColorInt converts a hex string to an integer representation
func hexToColorInt(hex string) (int, error) {
	if len(hex) == 7 && hex[0] == '#' {
		hex = hex[1:]
	} else if len(hex) != 6 {
		return 0, fmt.Errorf("invalid hex color format")
	}

	// Parse the hex string to an integer
	decimalValue, err := strconv.ParseInt(hex, 16, 0)
	if err != nil {
		return 0, fmt.Errorf("error parsing hex to int: %v", err)
	}

	// Return the integer value directly
	return int(decimalValue), nil
}

// isValidISO8601 checks if the string is a valid ISO8601 timestamp
func isValidISO8601(timestamp string) bool {
	_, err := time.Parse(time.RFC3339, timestamp)
	return err == nil
}

// SendWebhook sends the webhook payload to the specified Discord Webhook URL
func SendWebhook(webhookUrl string, webhookPayload Webhook) error {

	jsonData, err := json.Marshal(webhookPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %v", err)
	}

	resp, err := http.Post(webhookUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to post to Discord: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("discord webhook returned status %d", resp.StatusCode)
	}

	return nil
}
