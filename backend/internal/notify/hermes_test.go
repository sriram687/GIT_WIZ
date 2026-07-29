package notify

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sriram687/GIT_WIZ/backend/internal/domain"
)

func TestHermesDispatch(t *testing.T) {
	var mu sync.Mutex
	receivedSlack := ""
	receivedDiscord := ""

	// Set up mock Slack and Discord servers
	slackServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err == nil {
			if text, ok := payload["text"].(string); ok {
				receivedSlack = text
			}
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer slackServer.Close()

	discordServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err == nil {
			if content, ok := payload["content"].(string); ok {
				receivedDiscord = content
			}
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer discordServer.Close()

	hub := NewHub()
	hermes := NewHermes(hub, nil)
	hermes.slackURL = slackServer.URL
	hermes.discordURL = discordServer.URL

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go hermes.Run(ctx)

	// Wait briefly for subscription to register
	time.Sleep(50 * time.Millisecond)

	rec := domain.NotificationRecord{
		ID:        "ntf_test",
		SessionID: "session_123",
		ProjectID: "project_abc",
		PRURL:     "https://github.com/sriram687/GIT_WIZ/pull/1",
		Type:      domain.NotificationNeedsInput,
		Title:     "Need help",
		Body:      "Review this conflict",
		CreatedAt: time.Now(),
	}

	err := hub.Publish(ctx, rec)
	if err != nil {
		t.Fatalf("failed to publish to hub: %v", err)
	}

	// Wait for async dispatch to complete
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if !strings.Contains(receivedSlack, "[Chronos Alert]") || !strings.Contains(receivedSlack, "Need help") {
		t.Errorf("Slack did not receive correct message, got: %q", receivedSlack)
	}
	if !strings.Contains(receivedDiscord, "[Chronos Alert]") || !strings.Contains(receivedDiscord, "Need help") {
		t.Errorf("Discord did not receive correct message, got: %q", receivedDiscord)
	}
}
