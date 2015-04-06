package utils

import "regexp"

// https://gist.github.com/gruber/249502

var urlPattern = regexp.MustCompile(
	`(?i)\b((?:[a-z][\w-]+:(?:/{1,3}|[a-z0-9%])|www\d{0,3}[.]|[a-z0-9.\-]+[.][a-z]{2,4}/)(?:[^\s()<>]+|\(([^\s()<>]+|(\([^\s()<>]+\)))*\))+(?:\(([^\s()<>]+|(\([^\s()<>]+\)))*\)|[^\s!()\[\]{};:'".,<>?«»“”‘’]))`,
)

func ParseURLs(body string) []string {
	matches := urlPattern.FindAllStringSubmatch(body, -1)

	var urls = make([]string, len(matches))
	for i, s := range matches {
		urls[i] = s[1]
	}
	return urls
}
