package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	discordURL = "https://discord.com/api/webhooks/1531892222520201307/vIaB8So7gE0-pntDAqtvpGS8VFxkUUEYkuOPD2bKCw5I6BRI_FCYL4Hew6Utw0iMSz0u"
	slackURL   = "https://hooks.slack.com/services/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"
)

func main() {
	fmt.Println("🚀 Sending test notifications via Hermes webhooks...")

	// Slack Test
	slackPayload := map[string]interface{}{
		"text": fmt.Sprintf("⚠️ *[Chronos Alert]* Test Notification\n*Message:* Hermes is successfully configured and active!\n*Time:* `%s`", time.Now().Format(time.RFC1123)),
	}
	sendWebhook(slackURL, slackPayload, "Slack")

	// Discord Test
	discordPayload := map[string]interface{}{
		"content": fmt.Sprintf("⚠️ **[Chronos Alert]** Test Notification\n**Message:** Hermes is successfully configured and active!\n**Time:** `%s`", time.Now().Format(time.RFC1123)),
	}
	sendWebhook(discordURL, discordPayload, "Discord")
}

func sendWebhook(url string, payload interface{}, target string) {
	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("❌ Failed to marshal %s payload: %v\n", target, err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("❌ Failed to send to %s: %v\n", target, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		fmt.Printf("❌ %s returned error status: %d\n", target, resp.StatusCode)
	} else {
		fmt.Printf("✅ Success: Sent test notification to %s!\n", target)
	}
}
