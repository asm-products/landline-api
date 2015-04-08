package utils

import "testing"

func TestParseRoomMentions(t *testing.T) {
    cases := []struct {
        in string
        want []string
    }{
        {"Hey check out my new room at #general", []string{"general"}},
        {"#general", []string{"general"}},
        {"There are #many #room #mentions in this #sentence", []string{"many", "room", "mentions", "sentence"}},
        {"This sentence has no room mentions", []string{}},
        {"", []string{}},
    }
    for _, c := range cases {
        got := ParseRoomMentions(c.in)
        if len(got) != len(c.want) {
            t.Errorf("ParseRoomMentions(%q): got %d elements, want %d elements", c.in, len(got), len(c.want))
        } else {
            for i, _ := range c.want {
                if got[i] != c.want[i] {
                    t.Errorf("ParseRoomMentions(%q): got %s, want %s", c.in, got[i], c.want[i])
                }
            }
        }
    }
}

func TestParseUserMentions(t *testing.T) {
    cases := []struct {
        in string
        want []string
    }{
        {"Hey @tom did you see this yet?", []string{"tom"}},
        {"@QaZ", []string{"QaZ"}},
        {"@paul, @george, @ringo and @john are coming", []string{"paul", "george", "ringo", "john"}},
        {"This sentence has no user mentions", []string{}},
        {"", []string{}},
    }
    for _, c := range cases {
        got := ParseUserMentions(c.in)
        if len(got) != len(c.want) {
            t.Errorf("ParseUserMentions(%q): got %d elements, want %d elements", c.in, len(got), len(c.want))
        } else {
            for i, _ := range c.want {
                if got[i] != c.want[i] {
                    t.Errorf("ParseUserMentions(%q): got %s, want %s", c.in, got[i], c.want[i])
                }
            }
        }
    }
}
