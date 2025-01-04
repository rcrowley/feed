package main

import (
	"bytes"
	"os"
	"regexp"
	"testing"
)

func TestMain(t *testing.T) {
	stdout := &bytes.Buffer{}
	Main([]string{"feed", "-a", "Author Name", "-t", "Site Name", "-u", "http://example.com"}, os.Stdin, stdout)
	actual := stdout.String()
	pattern := `<\?xml version="1.0" encoding="UTF-8"\?>
<feed xmlns="http://www.w3.org/2005/Atom">
	<author>
		<name>Author Name</name>
	</author>
	<id>http://example.com/</id>
	<link href="http://example.com/" rel="alternate"></link>
	<title>Site Name</title>
	<updated>[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}[-+][0-9]{2}:[0-9]{2}</updated>
</feed>
`
	if matched, err := regexp.MatchString(pattern, actual); err != nil {
		t.Fatal(err)
	} else if !matched {
		t.Fatal(actual)
	}
}
