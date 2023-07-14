package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

type Configuration struct {
	BotUserID   string
	WebhookURL  string
	BearerToken string
}

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type BotWebhookPlugin struct {
	plugin.MattermostPlugin
	configuration *Configuration
}

func (p *BotWebhookPlugin) OnConfigurationChange() error {
	var configuration Configuration
	if err := p.API.LoadPluginConfiguration(&configuration); err != nil {
		p.API.LogError("Failed to load configuration", "error", err.Error())
		return err
	}
	p.configuration = &configuration
	return nil
}

func (p *BotWebhookPlugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	p.API.LogDebug("MessageHasBeenPosted")

	channel, err := p.API.GetChannel(post.ChannelId)

	if err != nil {
		p.API.LogError("Failed to get channel", "error", err.Error())
		return
	}

	if post.UserId == p.configuration.BotUserID {
		return
	}

	if strings.Contains(channel.Name, p.configuration.BotUserID) {
		p.API.LogDebug("Message to bot detected", "channel", channel.Name, "user", post.UserId, "message", post.Message)

		jsonPayload, err := json.Marshal(post)
		if err != nil {
			p.API.LogError("Failed to marshal JSON payload", "error", err.Error())
			return
		}

		req, err := http.NewRequest("POST", p.configuration.WebhookURL, bytes.NewBuffer(jsonPayload))
		req.Header.Set("Authorization", "Bearer "+p.configuration.BearerToken)
		if err != nil {
			p.API.LogError("Failed to create HTTP request", "error", err.Error())
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			p.API.LogError("Failed to make an HTTP request", "error", err.Error())
			return
		}
		defer resp.Body.Close()

		return
	}
}

func (p *BotWebhookPlugin) OnActivate() error {
	return p.OnConfigurationChange()
}

func main() {
	plugin.ClientMain(&BotWebhookPlugin{})
}
