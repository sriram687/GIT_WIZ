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
	defaultDiscordURL  = "https://discord.com/api/webhooks/1531892222520201307/vIaB8So7gE0-pntDAqtvpGS8VFxkUUEYkuOPD2bKCw5I6BRI_FCYL4Hew6Utw0iMSz0u"
	defaultSlackURL    = "https://hooks.slack.com/services/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"
	defaultNotionToken = "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"
	defaultNotionDBID  = "3ac38a82-7e2d-816a-a068-ee32cf08536b"
)

// Hermes (The Communicator) routes notifications/alerts to Slack, Discord, and Notion.
type Hermes struct {
	hub         *Hub
	slackURL    string
	discordURL  string
	notionToken string
	notionDBID  string
	client      *http.Client
	log         *slog.Logger
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

	notionTkn := os.Getenv("HERMES_NOTION_TOKEN")
	if notionTkn == "" {
		notionTkn = defaultNotionToken
	}

	notionDB := os.Getenv("HERMES_NOTION_DATABASE_ID")
	if notionDB == "" {
		notionDB = defaultNotionDBID
	}

	return &Hermes{
		hub:         hub,
		slackURL:    slack,
		discordURL:  discord,
		notionToken: notionTkn,
		notionDBID:  notionDB,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		log: log,
	}
}

// Run subscribes to the notification Hub and processes alerts until the context is cancelled.
func (h *Hermes) Run(ctx context.Context) {
	if h.slackURL == "" && h.discordURL == "" && h.notionToken == "" {
		h.log.Info("Hermes: Webhook and Notion integrations disabled (no credentials configured)")
		return
	}

	h.log.Info("Hermes: Agent is active and watching for notifications", 
		"slack_enabled", h.slackURL != "", 
		"discord_enabled", h.discordURL != "",
		"notion_enabled", h.notionToken != "" && h.notionDBID != "",
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
	// Dispatch concurrently to prevent slow external requests from blocking the listener loop
	go func() {
		if h.slackURL != "" {
			h.sendSlack(rec)
		}
		if h.discordURL != "" {
			h.sendDiscord(rec)
		}
		if h.notionToken != "" && h.notionDBID != "" {
			h.syncNotion(rec)
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

// syncNotion queries the database to see if a page for the session already exists.
// If it does, it updates the page. Otherwise, it creates a new page.
func (h *Hermes) syncNotion(rec domain.NotificationRecord) {
	pageID, err := h.findNotionPage(rec.SessionID)
	if err != nil {
		h.log.Error("Hermes: Failed to query Notion for existing page", "err", err)
		return
	}

	status := mapNotificationTypeToNotionStatus(rec.Type)

	if pageID != "" {
		h.updateNotionPage(pageID, rec, status)
	} else {
		h.createNotionPage(rec, status)
	}
}

func (h *Hermes) findNotionPage(sessionID domain.SessionID) (string, error) {
	url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", h.notionDBID)
	payload := map[string]interface{}{
		"filter": map[string]interface{}{
			"property": "Session",
			"rich_text": map[string]string{
				"equals": string(sessionID),
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+h.notionToken)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("Notion query status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Results []struct {
			ID string `json:"id"`
		} `json:"results"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	if len(result.Results) > 0 {
		return result.Results[0].ID, nil
	}

	return "", nil
}

func (h *Hermes) createNotionPage(rec domain.NotificationRecord, status string) {
	url := "https://api.notion.com/v1/pages"

	properties := map[string]interface{}{
		"Title": map[string]interface{}{
			"title": []map[string]interface{}{
				{
					"text": map[string]string{
						"content": rec.Title,
					},
				},
			},
		},
		"Session": map[string]interface{}{
			"rich_text": []map[string]interface{}{
				{
					"text": map[string]string{
						"content": string(rec.SessionID),
					},
				},
			},
		},
		"Project": map[string]interface{}{
			"rich_text": []map[string]interface{}{
				{
					"text": map[string]string{
						"content": string(rec.ProjectID),
					},
				},
			},
		},
		"Status": map[string]interface{}{
			"status": map[string]string{
				"name": status,
			},
		},
		"Latest Changes": map[string]interface{}{
			"rich_text": []map[string]interface{}{
				{
					"text": map[string]string{
						"content": rec.Body,
					},
				},
			},
		},
	}

	if rec.PRURL != "" {
		properties["PR URL"] = map[string]interface{}{
			"url": rec.PRURL,
		}
	}

	payload := map[string]interface{}{
		"parent": map[string]string{
			"database_id": h.notionDBID,
		},
		"properties": properties,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		h.log.Error("Hermes: Failed to marshal Notion insert payload", "err", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		h.log.Error("Hermes: Failed to create Notion insert request", "err", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+h.notionToken)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		h.log.Error("Hermes: Failed to insert to Notion", "err", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		h.log.Error("Hermes: Notion insert returned non-ok status code", "status", resp.StatusCode, "body", string(respBody))
	} else {
		h.log.Info("Hermes: Created card in Notion database", "session_id", rec.SessionID)
	}
}

func (h *Hermes) updateNotionPage(pageID string, rec domain.NotificationRecord, status string) {
	url := fmt.Sprintf("https://api.notion.com/v1/pages/%s", pageID)

	properties := map[string]interface{}{
		"Status": map[string]interface{}{
			"status": map[string]string{
				"name": status,
			},
		},
		"Latest Changes": map[string]interface{}{
			"rich_text": []map[string]interface{}{
				{
					"text": map[string]string{
						"content": rec.Body,
					},
				},
			},
		},
	}

	if rec.PRURL != "" {
		properties["PR URL"] = map[string]interface{}{
			"url": rec.PRURL,
		}
	}

	payload := map[string]interface{}{
		"properties": properties,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		h.log.Error("Hermes: Failed to marshal Notion update payload", "err", err)
		return
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	if err != nil {
		h.log.Error("Hermes: Failed to create Notion update request", "err", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+h.notionToken)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		h.log.Error("Hermes: Failed to update Notion page", "err", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		h.log.Error("Hermes: Notion update returned non-ok status code", "status", resp.StatusCode, "body", string(respBody))
	} else {
		h.log.Info("Hermes: Updated card status in Notion database", "session_id", rec.SessionID, "status", status)
	}
}

func mapNotificationTypeToNotionStatus(t domain.NotificationType) string {
	switch t {
	case domain.NotificationNeedsInput, domain.NotificationReadyForReview:
		return "Coding"
	case domain.NotificationReadyToMerge:
		return "Ready to Merge"
	case domain.NotificationPRMerged:
		return "Merged"
	case domain.NotificationPRClosedUnmerged:
		return "Closed"
	default:
		return "Coding"
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
