package models

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestReplaceUrlsWithLinks(t *testing.T) {
	messageBody := `
    Check out https://www.yahoo.com for more info,
    and also https://duckduckgo.com!
    `
	expectedBody := `
    Check out <a href="https://www.yahoo.com" target="_top">https://www.yahoo.com</a> for more info,
    and also <a href="https://duckduckgo.com" target="_top">https://duckduckgo.com</a>!
    `
	urls := []string{"https://www.yahoo.com", "https://duckduckgo.com"}
	message := &Message{Body: messageBody}

	result := replaceUrlsWithLinks(message, urls)
	if result != expectedBody {
		t.Errorf("TestReplaceUrlsWithLinks: got (%s), want (%s)", result, expectedBody)
	}
}

func TestReplaceUserMentionsWithLinks(t *testing.T) {
	room := makeFakeRoom()
	room.Id = "TestReplaceUserMentionsWithLinks-1"
	room.TeamId = "TestReplaceUserMentionsWithLinks-team1"
	_ = insertFakeRoom(room, t)

	user := makeFakeUserWithPrefixAndId("TestReplaceUserMentionsWithLinks", 1)
	user.TeamId = room.TeamId
	err := Db.Insert(user)
	if err != nil {
		t.Fatal("TestReplaceUserMentionsWithLinks error:", err)
	}

	messageBody := fmt.Sprintf("Have you seen @%s around?", user.Username)
	expectedBody := fmt.Sprintf(
		`Have you seen <a href="%s" target="_top">@%s</a> around?`,
		user.ProfileUrl,
		user.Username)

	message := &Message{
		RoomId: room.Id,
		Body:   messageBody,
	}

	result := replaceUserMentionsWithLinks(message, []string{user.Username})
	if result != expectedBody {
		t.Errorf("TestReplaceUserMentionsWithLinks: got (%s), want (%s)", result, expectedBody)
	}
}

func TestReplaceRoomMentionsWithLinks(t *testing.T) {
	room := makeFakeRoom()
	room.Id = "TestReplaceRoomMentionsWithLinks-1"
	room.TeamId = "TestReplaceRoomMentionsWithLinks-team1"
	room.Slug = "TestReplaceRoomMentionsWithLinks-slug"
	room.Topic = "TestReplaceRoomMentionsWithLinks-topic"
	_ = insertFakeRoom(room, t)

	message := &Message{
		RoomId: room.Id,
		Body:   fmt.Sprintf("There's a lot going on at #%s", room.Slug),
	}
	expectedBody := fmt.Sprintf(
		`There's a lot going on at <a href="#/rooms/%s" target="_top" title="%s">#%s</a>`,
		room.Slug,
		room.Topic,
		room.Slug)
	result := replaceRoomMentionsWithLinks(message, []string{room.Slug})
	if result != expectedBody {
		t.Errorf("TestReplaceRoomMentionsWithLinks: got (%s), want (%s)", result, expectedBody)
	}
}

func TestRegisterUnread(t *testing.T) {
	roomId := "TestRegisterUnread-room"
	userId := "TestRegisterUnread-user"
	userId2 := "TestRegisterUnread-user2"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.Method != "POST" {
			t.Errorf("TestRegisterUnread: got (%s), want (%s)", r.Method, "POST")
		}
		if r.RequestURI != "/articles" {
			t.Errorf("TestRegisterUnread: got (%s), want (%s)", r.RequestURI, "/articles")
		}
		if r.Header.Get("Authorization") != "Basic eHl6Og==" {
			t.Errorf("TestRegisterUnread: got (%s), want (%s)", r.Header.Get("Authorization"), "Basic eHl6Og==")
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal("TestRegisterUnread error:", err)
		}
		data := string(body)
		expected := fmt.Sprintf(`{"key":"%s","recipients":["%s"]}`, roomId, userId2)
		if data != expected {
			t.Errorf("TestRegisterUnread: got (%s), want (%s)", data, expected)
		}
	}))
	defer server.Close()
	err := os.Setenv("RR_URL", server.URL)
	if err != nil {
		t.Fatal("TestRegisterUnread error:", err)
	}
	err = os.Setenv("RR_PRIVATE_KEY", "xyz")
	if err != nil {
		t.Fatal("TestRegisterUnread error:", err)
	}
	membership := makeFakeRoomMembership()
	membership.Id = "TestRegisterUnread-1"
	membership.RoomId = roomId
	membership.UserId = userId2
	_ = insertFakeRoomMembership(membership, t)

	err = registerUnread(roomId, userId)
	if err != nil {
		t.Fatal("TestRegisterUnread error:", err)
	}
}

func TestParseMessage(t *testing.T) {
	teamId := "TestParseMessage-team1"
	user := makeFakeUser()
	user.Id = "TestParseMessage-1"
	user.Username = "TestParseMessage-user1"
	user.ProfileUrl = "www.profileurl.com"
	user.AvatarUrl = "www.avatarurl.com"
	user.TeamId = teamId
	err := Db.Insert(user)
	if err != nil {
		t.Fatal("TestParseMessage error:", err)
	}
	room := makeFakeRoom()
	room.Id = "TestParseMessage-1"
	room.Slug = "TestParseMessage-slug1"
	room.TeamId = teamId
	room.Topic = "TestParseMessage-topic"
	_ = insertFakeRoom(room, t)

	url := "http://www.randomurl.com"

	body := fmt.Sprintf(
		`This message has user @%s, room #%s, and url %s mentions.
# header 1
## header 2`,
		user.Username,
		room.Slug,
		url)
	expectedBody := `<p>This message has user <a href="www.profileurl.com" rel="nofollow">@TestParseMessage-user1</a>, room <a href="#/rooms/TestParseMessage-slug1" title="TestParseMessage-topic" rel="nofollow">#TestParseMessage-slug1</a>, and url <a href="http://www.randomurl.com" rel="nofollow">http://www.randomurl.com</a> mentions.</p>

<h1>header 1</h1>

<h2>header 2</h2>
`

	m := &Message{
		RoomId: room.Id,
		UserId: user.Id,
		Body:   body,
	}
	result := ParseMessage(m)
	if result != expectedBody {
		t.Errorf("TestParseMessage: got (%s), want (%s)", result, expectedBody)
	}
}
