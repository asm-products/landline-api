package utils

import (
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
