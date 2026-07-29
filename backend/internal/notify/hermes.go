package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/sriram687/GIT_WIZ/backend/internal/domain"
)

const (
	defaultDiscordURL = "https://discord.com/api/webhooks/1531892222520201307/vIaB8So7gE0-pntDAqtvpGS8VFxkUUEYkuOPD2bKCw5I6BRI_FCYL4Hew6Utw0iMSz0u"
	defaultSlackURL   = "https://hooks.slack.com/services/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"
)

// Hermes (The Communicator) routes notifications/alerts to Slack and Discord.
type Hermes struct {
	hub        *Hub
	slackURL   string
	discordURL string
	client     *http.Client
	log        *slog.Logger
}

// NewHermes constructs a Hermes notifier.
func NewHermes(hub *Hub, log *slog.Logger) *Hermes {
	if log == nil {
		log = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	slack := os.Getenv("HERMES_SLACK_WEBHOOK_URL")
	if slack == "" {
		slack = defaultSlackURL
	}

	discord := os.Getenv("HERMES_DISCORD_WEBHOOK_URL")
	if discord == "" {
		discord = defaultDiscordURL
	}

	return &Hermes{
		hub:        hub,
		slackURL:   slack,
		discordURL: discord,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		log: log,
	}
}

// Run subscribes to the notification Hub and processes alerts until the context is cancelled.
func (h *Hermes) Run(ctx context.Context) {
	if h.slackURL == "" && h.discordURL == "" {
		h.log.Info("Hermes: Webhook integration disabled (no URLs configured)")
		return
	}

	h.log.Info("Hermes: Agent is active and watching for notifications", 
		"slack_enabled", h.slackURL != "", 
		"discord_enabled", h.discordURL != "",
	)

	ch, unsubscribe := h.hub.Subscribe("")
	defer unsubscribe()

	for {
		select {
		case <-ctx.Done():
			h.log.Info("Hermes: Shutting down agent")
			return
		case rec, ok := <-ch:
			if !ok {
				return
			}
			h.dispatch(rec)
		}
	}
}

func (h *Hermes) dispatch(rec domain.NotificationRecord) {
	// Dispatch concurrently to prevent slow webhooks from blocking the listener loop
	go func() {
		if h.slackURL != "" {
			h.sendSlack(rec)
		}
		if h.discordURL != "" {
			h.sendDiscord(rec)
		}
	}()
}

func (h *Hermes) sendSlack(rec domain.NotificationRecord) {
	text := h.formatSlackMessage(rec)
	payload := map[string]interface{}{
		"text": text,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		h.log.Error("Hermes: Failed to marshal Slack payload", "err", err)
		return
	}

	req, err := http.NewRequest("POST", h.slackURL, bytes.NewBuffer(body))
	if err != nil {
		h.log.Error("Hermes: Failed to create Slack request", "err", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		h.log.Error("Hermes: Failed to post to Slack", "err", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		h.log.Error("Hermes: Slack returned non-ok status code", "status", resp.StatusCode, "body", string(respBody))
	} else {
		h.log.Debug("Hermes: Notification sent to Slack", "notification_id", rec.ID)
	}
}

func (h *Hermes) sendDiscord(rec domain.NotificationRecord) {
	content := h.formatDiscordMessage(rec)
	payload := map[string]interface{}{
		"content": content,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		h.log.Error("Hermes: Failed to marshal Discord payload", "err", err)
		return
	}

	req, err := http.NewRequest("POST", h.discordURL, bytes.NewBuffer(body))
	if err != nil {
		h.log.Error("Hermes: Failed to create Discord request", "err", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		h.log.Error("Hermes: Failed to post to Discord", "err", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		h.log.Error("Hermes: Discord returned non-ok status code", "status", resp.StatusCode, "body", string(respBody))
	} else {
		h.log.Debug("Hermes: Notification sent to Discord", "notification_id", rec.ID)
	}
}

func (h *Hermes) formatSlackMessage(rec domain.NotificationRecord) string {
	var emoji string
	switch rec.Type {
	case domain.NotificationNeedsInput, domain.NotificationReadyForReview:
		emoji = "⚠️"
	case domain.NotificationReadyToMerge:
		emoji = "🚀"
	case domain.NotificationPRMerged:
		emoji = "✅"
	case domain.NotificationPRClosedUnmerged:
		emoji = "❌"
	default:
		emoji = "📢"
	}

	var prLink string
	if rec.PRURL != "" {
		prLink = "\n*PR Link:* <" + rec.PRURL + "|" + rec.PRURL + ">"
	}

	return fmt.Sprintf(
		"%s *[Chronos Alert]* %s\n*Message:* %s%s\n*Session:* `%s` | *Project:* `%s`",
		emoji, rec.Title, rec.Body, prLink, rec.SessionID, rec.ProjectID,
	)
}

func (h *Hermes) formatDiscordMessage(rec domain.NotificationRecord) string {
	var emoji string
	switch rec.Type {
	case domain.NotificationNeedsInput, domain.NotificationReadyForReview:
		emoji = "⚠️"
	case domain.NotificationReadyToMerge:
		emoji = "🚀"
	case domain.NotificationPRMerged:
		emoji = "✅"
	case domain.NotificationPRClosedUnmerged:
		emoji = "❌"
	default:
		emoji = "📢"
	}

	var prLink string
	if rec.PRURL != "" {
		prLink = "\n**PR Link:** " + rec.PRURL
	}

	return fmt.Sprintf(
		"%s **[Chronos Alert]** %s\n**Message:** %s%s\n**Session:** `%s` | **Project:** `%s`",
		emoji, rec.Title, rec.Body, prLink, rec.SessionID, rec.ProjectID,
	)
}
