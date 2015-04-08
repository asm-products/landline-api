package utils

import (
    "testing"
    "net/http/httptest"
    "net/http"
    "encoding/json"
    "io/ioutil"
)

func TestPostMessageToWebhook(t *testing.T) {
    m := Message{
        Slug: "slug12345",
        UserId: "astro",
        Body: "message body",
    }
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer r.Body.Close()
        if r.Header.Get("Content-Type") != "application/json" {
            t.Errorf("PostMessageToWebhook: got %s, want %s", r.Header.Get("Content-Type"), "application/json")
        }
        if r.Header.Get("Authorization") != "Basic bXlzZWNyZXQ6" {
            t.Errorf("PostMessageToWebhook: got %s, want %s", r.Header.Get("Authorization"), "Basic bXlzZWNyZXQ6")
        }
        if r.Method != "POST" {
            t.Errorf("PostMessageToWebhook: got %s, want %s", r.Method, "POST")
        }

        webhook := MessageWebhookBody{}
        data, err := ioutil.ReadAll(r.Body)
        err = json.Unmarshal(data, &webhook)
        if err != nil {
            t.Error(err)
        }
        if webhook.WebhookType != "message" {
            t.Errorf("PostMessageToWebhook: got %s, want %s", webhook.WebhookType, "message")
        }
        if webhook.Message.Slug != m.Slug {
            t.Errorf("PostMessageToWebhook: got %s, want %s", webhook.Message.Slug, m.Slug)
        }
        if webhook.Message.UserId != m.UserId {
            t.Errorf("PostMessageToWebhook: got %s, want %s", webhook.Message.UserId, m.UserId)
        }
        if webhook.Message.Body != m.Body {
            t.Errorf("PostMessageToWebhook: got %s, want %s", webhook.Message.Body, m.Body)
        }
    }))
    defer ts.Close()
    err := PostMessageToWebhook(ts.URL, "mysecret", m)
    if err != nil {
        t.Error(err)
    }
}