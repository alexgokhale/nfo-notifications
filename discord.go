package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

func sendDiscordWebhook(url string, webhook Webhook) (*http.Response, error) {
	requestContent, err := json.Marshal(webhook)

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// Webhook represents the content sent to a Discord webhook URL to send a message.
//
// Either Content, Embeds or both must have a value for a message to be sent successfully.
type Webhook struct {
	Content   string  `json:"content,omitempty"`
	Username  string  `json:"username,omitempty"`
	AvatarURL string  `json:"avatar_url,omitempty"`
	Embeds    []Embed `json:"embeds,omitempty"`
}

// Embed represents a rich content object to be rendered in Discord.
type Embed struct {
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	URL         string     `json:"url,omitempty"`
	Color       int        `json:"color,omitempty"`
	Timestamp   *time.Time `json:"timestamp,omitempty"`
	Footer      *Footer    `json:"footer,omitempty"`
	Thumbnail   *Thumbnail `json:"thumbnail,omitempty"`
	Image       *Image     `json:"image,omitempty"`
	Author      *Author    `json:"author,omitempty"`
	Fields      []*Field   `json:"fields,omitempty"`
}

// Footer represents an icon and text to appear at the bottom of an embed.
type Footer struct {
	IconURL string `json:"icon_url,omitempty"`
	Text    string `json:"text"`
}

// Thumbnail represents a small image to be placed in the corner of an embed.
type Thumbnail struct {
	URL string `json:"url"`
}

// Image represents a full-sized image to be placed into an embed.
type Image struct {
	URL string `json:"url"`
}

// Author represents the author of the embed.
type Author struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	IconURL string `json:"icon_url"`
}

// Field represents a name-value pair to be displayed in an embed.
type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}
