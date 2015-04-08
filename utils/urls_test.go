package utils

import "testing"

func TestParseUrls(t *testing.T) {
    cases := []struct {
        in string
        want []string
    }{
        {"Here is a link: duckduckgo.com", []string{}},
        {"Here is a link: https://duckduckgo.com", []string{"https://duckduckgo.com"}},
        {"Sites you can check: www.google.com and www.yahoo.com", []string{"www.google.com", "www.yahoo.com"}},
        {"www.", []string{}},
        {"", []string{}},
    }
    for _, c := range cases {
        got := ParseURLs(c.in)
        if len(got) != len(c.want) {
            t.Errorf("ParseURLs(%q): got %d elements, want %d elements", c.in, len(got), len(c.want))
        } else {
            for i, _ := range c.want {
                if got[i] != c.want[i] {
                    t.Errorf("ParseURLs(%q): got %s, want %s", c.in, got[i], c.want[i])
                }
            }
        }
    }
}
