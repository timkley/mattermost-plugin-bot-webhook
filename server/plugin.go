package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

const pluginID = "bot-webhook-plugin"

// Configuration maps to the lowercase keys defined in plugin.json settings_schema.
// LoadPluginConfiguration uses case-insensitive JSON unmarshaling, so the Go
// field names do not need to match the key casing exactly.
type Configuration struct {
	BotUsername string `json:"botusername"`
	BotUserID   string `json:"botuserid"`
	WebhookURL  string `json:"webhookurl"`
	BearerToken string `json:"bearertoken"`

	BotUsername2 string `json:"botusername2"`
	BotUserID2   string `json:"botuserid2"`
	WebhookURL2  string `json:"webhookurl2"`
	BearerToken2 string `json:"bearertoken2"`

	BotUsername3 string `json:"botusername3"`
	BotUserID3   string `json:"botuserid3"`
	WebhookURL3  string `json:"webhookurl3"`
	BearerToken3 string `json:"bearertoken3"`

	BotUsername4 string `json:"botusername4"`
	BotUserID4   string `json:"botuserid4"`
	WebhookURL4  string `json:"webhookurl4"`
	BearerToken4 string `json:"bearertoken4"`

	BotUsername5 string `json:"botusername5"`
	BotUserID5   string `json:"botuserid5"`
	WebhookURL5  string `json:"webhookurl5"`
	BearerToken5 string `json:"bearertoken5"`

	BotUsername6 string `json:"botusername6"`
	BotUserID6   string `json:"botuserid6"`
	WebhookURL6  string `json:"webhookurl6"`
	BearerToken6 string `json:"bearertoken6"`

	BotUsername7 string `json:"botusername7"`
	BotUserID7   string `json:"botuserid7"`
	WebhookURL7  string `json:"webhookurl7"`
	BearerToken7 string `json:"bearertoken7"`

	BotUsername8 string `json:"botusername8"`
	BotUserID8   string `json:"botuserid8"`
	WebhookURL8  string `json:"webhookurl8"`
	BearerToken8 string `json:"bearertoken8"`

	BotUsername9 string `json:"botusername9"`
	BotUserID9   string `json:"botuserid9"`
	WebhookURL9  string `json:"webhookurl9"`
	BearerToken9 string `json:"bearertoken9"`

	BotUsername10 string `json:"botusername10"`
	BotUserID10   string `json:"botuserid10"`
	WebhookURL10  string `json:"webhookurl10"`
	BearerToken10 string `json:"bearertoken10"`
}

type botMapping struct {
	BotUserID   string
	WebhookURL  string
	BearerToken string
}

// BotWebhookPlugin implements the interface expected by the Mattermost server.
type BotWebhookPlugin struct {
	plugin.MattermostPlugin
	configuration *Configuration
}

// activeBots returns the list of configured bot mappings with non-empty user IDs.
// If a bot has a username but no user ID, it attempts to resolve the username
// via the API as a fallback.
func (p *BotWebhookPlugin) activeBots() []botMapping {
	cfg := p.configuration
	all := []struct {
		username    string
		userID      string
		webhookURL  string
		bearerToken string
	}{
		{cfg.BotUsername, cfg.BotUserID, cfg.WebhookURL, cfg.BearerToken},
		{cfg.BotUsername2, cfg.BotUserID2, cfg.WebhookURL2, cfg.BearerToken2},
		{cfg.BotUsername3, cfg.BotUserID3, cfg.WebhookURL3, cfg.BearerToken3},
		{cfg.BotUsername4, cfg.BotUserID4, cfg.WebhookURL4, cfg.BearerToken4},
		{cfg.BotUsername5, cfg.BotUserID5, cfg.WebhookURL5, cfg.BearerToken5},
		{cfg.BotUsername6, cfg.BotUserID6, cfg.WebhookURL6, cfg.BearerToken6},
		{cfg.BotUsername7, cfg.BotUserID7, cfg.WebhookURL7, cfg.BearerToken7},
		{cfg.BotUsername8, cfg.BotUserID8, cfg.WebhookURL8, cfg.BearerToken8},
		{cfg.BotUsername9, cfg.BotUserID9, cfg.WebhookURL9, cfg.BearerToken9},
		{cfg.BotUsername10, cfg.BotUserID10, cfg.WebhookURL10, cfg.BearerToken10},
	}

	active := make([]botMapping, 0, len(all))
	for _, b := range all {
		id := b.userID
		if id == "" && b.username != "" {
			u := strings.TrimPrefix(b.username, "@")
			if user, appErr := p.API.GetUserByUsername(u); appErr == nil {
				id = user.Id
			}
		}
		if id != "" {
			active = append(active, botMapping{id, b.webhookURL, b.bearerToken})
		}
	}
	return active
}

// getString reads a string value from the settings map by its lowercase key.
// All keys in plugin.json are lowercase, and Mattermost submits form values
// with lowercase keys — so no fallback or case conversion is needed.
func getString(settings map[string]interface{}, key string) string {
	if v, ok := settings[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// ConfigurationWillBeSaved intercepts the config save, resolves @usernames to
// user IDs, and injects the resolved IDs into the config before it is persisted.
// If the username is empty or cannot be resolved, the User ID is cleared.
func (p *BotWebhookPlugin) ConfigurationWillBeSaved(newCfg *model.Config) (*model.Config, error) {
	if newCfg.PluginSettings.Plugins == nil {
		return nil, nil
	}
	settings, ok := newCfg.PluginSettings.Plugins[pluginID]
	if !ok {
		return nil, nil
	}

	// Remove legacy CamelCase keys written by versions ≤0.5.2.
	// After v0.5.3, canonical keys are all lowercase. If both casings exist in
	// the stored config, Go's case-insensitive JSON unmarshal picks one
	// non-deterministically — the old CamelCase value can "win" and load a stale ID.
	for _, k := range []string{
		"BotUsername", "BotUserID", "WebhookURL", "BearerToken",
		"BotUsername2", "BotUserID2", "WebhookURL2", "BearerToken2",
		"BotUsername3", "BotUserID3", "WebhookURL3", "BearerToken3",
		"BotUsername4", "BotUserID4", "WebhookURL4", "BearerToken4",
		"BotUsername5", "BotUserID5", "WebhookURL5", "BearerToken5",
		"BotUsername6", "BotUserID6", "WebhookURL6", "BearerToken6",
		"BotUsername7", "BotUserID7", "WebhookURL7", "BearerToken7",
		"BotUsername8", "BotUserID8", "WebhookURL8", "BearerToken8",
		"BotUsername9", "BotUserID9", "WebhookURL9", "BearerToken9",
		"BotUsername10", "BotUserID10", "WebhookURL10", "BearerToken10",
	} {
		delete(settings, k)
	}

	resolve := func(usernameKey, userIDKey string) {
		// Always clear the ID first — prevents stale values when username changes.
		settings[userIDKey] = ""

		username := getString(settings, usernameKey)
		if username == "" {
			return
		}
		u := strings.TrimPrefix(username, "@")
		user, appErr := p.API.GetUserByUsername(u)
		if appErr != nil {
			p.API.LogError("[BotWebhook] Failed to resolve username", "username", u, "error", appErr.Error())
			return
		}
		p.API.LogInfo("[BotWebhook] Resolved username", "username", u, "userID", user.Id)
		settings[userIDKey] = user.Id
	}

	resolve("botusername", "botuserid")
	resolve("botusername2", "botuserid2")
	resolve("botusername3", "botuserid3")
	resolve("botusername4", "botuserid4")
	resolve("botusername5", "botuserid5")
	resolve("botusername6", "botuserid6")
	resolve("botusername7", "botuserid7")
	resolve("botusername8", "botuserid8")
	resolve("botusername9", "botuserid9")
	resolve("botusername10", "botuserid10")

	newCfg.PluginSettings.Plugins[pluginID] = settings
	return newCfg, nil
}

func (p *BotWebhookPlugin) OnConfigurationChange() error {
	var configuration Configuration
	if err := p.API.LoadPluginConfiguration(&configuration); err != nil {
		p.API.LogError("[BotWebhook] Failed to load configuration", "error", err.Error())
		return err
	}
	p.configuration = &configuration
	return nil
}

func (p *BotWebhookPlugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	channel, err := p.API.GetChannel(post.ChannelId)
	if err != nil {
		p.API.LogError("[BotWebhook] Failed to get channel", "error", err.Error())
		return
	}

	for _, bot := range p.activeBots() {
		if post.UserId == bot.BotUserID {
			continue
		}
		if !strings.Contains(channel.Name, bot.BotUserID) {
			continue
		}

		p.API.LogInfo("[BotWebhook] Message to bot detected", "channel", channel.Name, "user", post.UserId, "bot", bot.BotUserID)

		jsonPayload, err := json.Marshal(post)
		if err != nil {
			p.API.LogError("[BotWebhook] Failed to marshal JSON payload", "error", err.Error())
			continue
		}

		req, err := http.NewRequest("POST", bot.WebhookURL, bytes.NewBuffer(jsonPayload))
		if err != nil {
			p.API.LogError("[BotWebhook] Failed to create HTTP request", "error", err.Error())
			continue
		}
		req.Header.Set("Authorization", "Bearer "+bot.BearerToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			p.API.LogError("[BotWebhook] Failed to make HTTP request", "error", err.Error())
			continue
		}
		resp.Body.Close()
		p.API.LogInfo("[BotWebhook] Webhook dispatched", "bot", bot.BotUserID, "url", bot.WebhookURL)
		break
	}
}

func (p *BotWebhookPlugin) OnActivate() error {
	return p.OnConfigurationChange()
}

func main() {
	plugin.ClientMain(&BotWebhookPlugin{})
}
