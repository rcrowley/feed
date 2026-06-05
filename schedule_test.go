package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// TestMainSchedule checks that, at a fixed reference time, the feed includes
// the published post but not the not-yet-published one (a future
// <time class="feed">) or the draft.
func TestMainSchedule(t *testing.T) {
	stdout := &bytes.Buffer{}
	Main([]string{
		"feed",
		"-a", "Author Name",
		"-t", "Site Name",
		"-u", "http://example.com",
		"-n", "2026-06-04 00:00:00",
		"testdata",
	}, os.Stdin, stdout)
	out := stdout.String()

	if !strings.Contains(out, "<title>Live Post</title>") {
		t.Errorf("expected Live Post in the feed, got:\n%s", out)
	}
	for _, title := range []string{"Future Post", "Draft Post"} {
		if strings.Contains(out, "<title>"+title+"</title>") {
			t.Errorf("did not expect %s in the feed, got:\n%s", title, out)
		}
	}
}
