package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var slackWebhookForContact string

func InitSlackForContact() {
	slackWebhookForContact = os.Getenv("SLACK_WEBHOOK_URL_CONTACT")
}

// postToSlack は Webhook URL が設定されていればメッセージを POST 送信します
func postToSlack(msg string) (err error) {
	if slackWebhookForContact == "" {
		return fmt.Errorf("slackWebhookForContact is not set")
	}
	payload := map[string]string{"text": msg}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", slackWebhookForContact, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func SendSlack(ctx context.Context, name, email, message string) error {
	msg := "New contact message\n" +
		"Name: " + name + "\n" +
		"Email: " + email + "\n" +
		"Message: " + message
	err := postToSlack(msg)
	return err
}
