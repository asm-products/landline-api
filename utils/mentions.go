package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
)

var MentionPattern *regexp.Regexp = regexp.MustCompile(
	`(?:^|\W)@((?i)[a-z0-9][a-z0-9-]*)`,
)

func ParseUserMentions(body string) []string {
	mentions := MentionPattern.FindAllStringSubmatch(body, -1)

	var usernames = make([]string, len(mentions))
	for i, s := range mentions {
		usernames[i] = s[1]
	}
	return usernames
}

type MentionWebhookBody struct {
	WebhookType string   `json:"webhook_type"`
	MessageBody string   `json:"message_body"`
	Usernames   []string `json:"usernames"`
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

func buildMentionsRequestBody(body string, usernames []string) (*bytes.Reader, error) {
	requestBody := MentionWebhookBody{
		WebhookType: "mention",
		MessageBody: body,
		Usernames:   usernames,
	}

	r, err := json.Marshal(requestBody)

	return bytes.NewReader(r), err
}
