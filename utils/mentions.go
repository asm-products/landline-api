package utils

import (
	"regexp"
)

var roomMentionPattern = regexp.MustCompile(
	`(?:^|\W)#((?i)[a-z0-9][a-z0-9-]*)`,
)

var userMentionPattern = regexp.MustCompile(
	`(?:^|\W)@((?i)[a-z0-9][a-z0-9-]*)`,
)

// ParseRoomMentions finds strings prefixed by `#` and returns
// an array of the strings without the prefix
func ParseRoomMentions(body string) []string {
	mentions := roomMentionPattern.FindAllStringSubmatch(body, -1)

	var rooms = make([]string, len(mentions))
	for i, s := range mentions {
		rooms[i] = s[1]
	}
	return rooms
}

// ParseUserMentions finds strings prefixed by `@` and returns
// an array of the strings without the prefix
func ParseUserMentions(body string) []string {
	mentions := userMentionPattern.FindAllStringSubmatch(body, -1)

	var usernames = make([]string, len(mentions))
	for i, s := range mentions {
		usernames[i] = s[1]
	}
	return usernames
}
