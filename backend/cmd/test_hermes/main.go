package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	discordURL  = "https://discord.com/api/webhooks/1531892222520201307/vIaB8So7gE0-pntDAqtvpGS8VFxkUUEYkuOPD2bKCw5I6BRI_FCYL4Hew6Utw0iMSz0u"
	slackURL    = "https://hooks.slack.com/services/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"
	notionToken = "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"
	notionDBID  = "3ac38a82-7e2d-816a-a068-ee32cf08536b"
)

func main() {
	fmt.Println("🚀 Sending test notifications via Hermes integrations...")

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

	// Notion Test
	sendNotionCard()
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

func sendNotionCard() {
	url := "https://api.notion.com/v1/pages"

	properties := map[string]interface{}{
		"Title": map[string]interface{}{
			"title": []map[string]interface{}{
				{
					"text": map[string]string{
						"content": "🛠️ Setup Notion Integration",
					},
				},
			},
		},
		"Session": map[string]interface{}{
			"rich_text": []map[string]interface{}{
				{
					"text": map[string]string{
						"content": "session_notion_setup",
					},
				},
			},
		},
		"Project": map[string]interface{}{
			"rich_text": []map[string]interface{}{
				{
					"text": map[string]string{
						"content": "git-wiz-repo",
					},
				},
			},
		},
		"Status": map[string]interface{}{
			"status": map[string]string{
				"name": "Ready to Merge",
			},
		},
		"Latest Changes": map[string]interface{}{
			"rich_text": []map[string]interface{}{
				{
					"text": map[string]string{
						"content": "Successfully setup database fields and validated synchronization between GitHub/Git worktrees and Notion task cards.",
					},
				},
			},
		},
		"PR URL": map[string]interface{}{
			"url": "https://github.com/sriram687/GIT_WIZ/pull/1",
		},
	}

	payload := map[string]interface{}{
		"parent": map[string]string{
			"database_id": notionDBID,
		},
		"properties": properties,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("❌ Failed to marshal Notion payload: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("❌ Failed to create Notion request: %v\n", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+notionToken)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Failed to connect to Notion: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Notion API returned error status %d: %s\n", resp.StatusCode, string(respBody))
	} else {
		fmt.Println("✅ Success: Sent test card and status to Notion!")
	}
}
