package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/rcrowley/mergician/html"
)

func TestFeed(t *testing.T) {
	f := &Feed{
		Author: "Author Name",
		Path:   "index.atom.xml",
		Title:  "Site Name",
		URL:    "http://example.com",

		Entries: []Entry{
			{
				Date: "2024-12-03 22:28:00",
				Path: "newest.html",

				Node: must2(html.ParseString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8"/>
<title>Newest Article Title — Site Name</title>
</head>
<body>
<header><h1>Site Name</h1></header>
<article class="body">
<time datetime="2024-12-03 22:28:00">2024-12-03 22:28:00</time>
<h1>Newest Article Title</h1>
<p>Newest article body.</p>
</article>
</body>
</html>
`)),
			},
			{
				Date: "1970-01-01 00:00:00",
				Path: "oldest.html",
				Node: must2(html.ParseString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8"/>
<title>Oldest Article Title — Site Name</title>
</head>
<body>
<header><h1>Site Name</h1></header>
<article class="body">
<time datetime="1970-01-01 00:00:00">1970-01-01 00:00:00</time>
<h1>Oldest Article Title</h1>
<p>Oldest article body.</p>
</article>
</body>
</html>
`)),
			},
		},

		t: time.Now(),
	}
	stdout := &bytes.Buffer{}
	if err := f.Render(stdout); err != nil {
		t.Fatal(err)
	}
	actual := stdout.String()
	expected := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
	<author>
		<name>Author Name</name>
	</author>
	<id>http://example.com/</id>
	<link href="http://example.com/" rel="alternate"></link>
	<link href="http://example.com/index.atom.xml" rel="self"></link>
	<title>Site Name</title>
	<updated>%s</updated>
	<entry>
		<id>http://example.com/newest.html</id>
		<link href="http://example.com/newest.html" rel="alternate"></link>
		<title>Newest Article Title</title>
		<updated>2024-12-03T22:28:00Z</updated>
		<content type="html">&lt;article class=&#34;body&#34;&gt;
&lt;time datetime=&#34;2024-12-03 22:28:00&#34;&gt;2024-12-03 22:28:00&lt;/time&gt;
&lt;h1&gt;Newest Article Title&lt;/h1&gt;
&lt;p&gt;Newest article body.&lt;/p&gt;
&lt;/article&gt;</content>
	</entry>
	<entry>
		<id>http://example.com/oldest.html</id>
		<link href="http://example.com/oldest.html" rel="alternate"></link>
		<title>Oldest Article Title</title>
		<updated>1970-01-01T00:00:00Z</updated>
		<content type="html">&lt;article class=&#34;body&#34;&gt;
&lt;time datetime=&#34;1970-01-01 00:00:00&#34;&gt;1970-01-01 00:00:00&lt;/time&gt;
&lt;h1&gt;Oldest Article Title&lt;/h1&gt;
&lt;p&gt;Oldest article body.&lt;/p&gt;
&lt;/article&gt;</content>
	</entry>
</feed>
`, f.t.Format(time.RFC3339))
	if actual != expected {
		t.Fatalf("actual: %s != expected: %s", actual, expected)
	}
}

func TestFeedAddIndexHTML(t *testing.T) {
	f := &Feed{}
	f.Add("1970-01-01", "test/index.html", nil)
	if len(f.Entries) != 1 || f.Entries[0].Path != "test/" {
		t.Fatalf("%+v", f.Entries)
	}
}
