package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type MentionWebhookBody struct {
	WebhookType string   `json:"webhook_type"`
	MessageBody string   `json:"message_body"`
	Usernames   []string `json:"usernames"`
}

type Message struct {
	RoomId string `json:"room_id"`
	UserId string `json:"user_id"`
	Body   string `json:"body"`
}

type MessageWebhookBody struct {
	WebhookType string  `json:"webhook_type"`
	Message     Message `json:"message"`
}

func PostMentionsToWebhook(url, secret, body string, usernames []string) error {
	requestBody, err := buildMentionsRequestBody(body, usernames)

	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		url,
		requestBody,
	)

	req.SetBasicAuth(secret, "")

	client := &http.Client{}

	_, err = client.Do(req)

	return err
}

func PostMessageToWebhook(url, secret string, message Message) error {
	requestBody := MessageWebhookBody{
		WebhookType: "message",
		Message:     message,
	}

	r, err := json.Marshal(requestBody)

	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewReader(r),
	)

	req.SetBasicAuth(secret, "")

	client := &http.Client{}

	_, err = client.Do(req)

	return err
}

func buildMentionsRequestBody(body string, usernames []string) (*bytes.Reader, error) {
	requestBody := MentionWebhookBody{
		WebhookType: "mention",
		MessageBody: body,
		Usernames:   usernames,
	}

	r, err := json.Marshal(requestBody)

	return bytes.NewReader(r), err
}
