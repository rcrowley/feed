package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/rcrowley/mergician/files"
	"github.com/rcrowley/mergician/html"
	"golang.org/x/net/html/atom"
)

func Main(args []string, stdin io.Reader, stdout io.Writer) {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	author := flags.String("a", "", "author's name")
	output := flags.String("o", "-", "write to this file instead of standard output")
	title := flags.String("t", "", "feed title")
	url := flags.String("u", "", "site URL with scheme and domain")
	verbose := flags.Bool("v", false, "verbose mode")
	exclude := files.NewStringSliceFlag(flags, "x", "subdirectory of <docroot> to exclude (may be repeated)")
	flags.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: feed -a <author> [-o <output>] -t <title> -u <url> [-v] [-x <exclude>[...]] [<docroot>[...]]
  -a <author>   author's name
  -o <output>   write to this file instead of standard output
  -t <title>    feed title
  -u <url>      site URL with scheme and domain
  -v            verbose mode
  -x <exclude>  subdirectory of <docroot> to exclude (may be repeated)
  <docroot>     document root directory to scan (defaults to the current working directory)

Synopsis: feed scans each <docroot> (or the current working directory) for <article> elements or other containers (including <body>) with class="feed" that contain a <time class="feed"> element, sorts them by those <time> elements, and constructs an Atom feed containing the most recent 10 articles.
`)
	}
	flags.Parse(args[1:])
	if *author == "" && *output == "" || *title == "" || *url == "" {
		flags.Usage()
		os.Exit(1)
	}

	var docroots []string
	if flags.NArg() == 0 {
		docroots = []string{"."}
	} else {
		docroots = flags.Args()
	}
	lists := must2(files.AllHTML(docroots, *exclude))

	feed := &Feed{
		Author: *author,
		Title:  *title,
		URL:    *url,
	}
	if *output != "-" {
		feed.Path = *output
	}

	var wg sync.WaitGroup
	for _, list := range lists {
		for _, path := range list.QualifiedPaths() {
			wg.Add(1)
			go func() {
				defer wg.Done()
				n := must2(files.Parse(path))
				if t := html.Find(n, html.All(
					html.IsAtom(atom.Time),
					html.HasAttr("class", "feed"),
				)); t != nil {
					must(feed.Add(html.Attr(t, "datetime"), path, n))
				}
			}()
		}
	}
	wg.Wait()

	for _, entry := range feed.Entries {
		if *verbose {
			fmt.Printf(
				"frag %s %s # %s\n", // "frag %q %q # %s\n",
				"<h1>", entry.Path, entry.Date,
			)
			if entry.Content.DataAtom == atom.Article {
				fmt.Printf(
					"frag %s %s # %s\n", // "frag %q %q # %s\n",
					"<article>", entry.Path, entry.Date,
				)
			} else {
				fmt.Printf(
					"frag '<%s class=\"feed\">' %s # %s\n", // "frag '<%s class=\"feed\">' %q # %s\n",
					entry.Content.DataAtom, entry.Path, entry.Date,
				)
			}
		}
	}

	var w io.Writer
	if *output == "-" {
		w = stdout
	} else {
		if *verbose {
			fmt.Printf("# wrote Atom feed to %s\n", *output)
		}
		f := must2(os.Create(*output))
		defer f.Close()
		w = f
	}
	must(feed.Render(w))

}

func init() {
	log.SetFlags(0)
}

func main() {
	Main(os.Args, os.Stdin, os.Stdout)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func must2[T any](v T, err error) T {
	must(err)
	return v
}
