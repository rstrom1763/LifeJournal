package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// WebhookPayload represents the structure of a Discord webhook message.
type WebhookPayload struct {
	Content string `json:"content"`
}

// SendWebhookMessage sends a simple content message to the specified Discord webhook URL.
func SendWebhookMessage(url, content string) error {
	payload := WebhookPayload{
		Content: content,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("failed to send webhook request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord webhook returned non-success status: %d", resp.StatusCode)
	}

	return nil
}
