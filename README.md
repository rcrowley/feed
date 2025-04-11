Feed
====

Feed scans document root directories for `<article>` elements or other containers (including `<body>`) with `class="feed"` that contain a `<time class="feed">` element, sorts them by those `<time>` elements, and constructs an Atom feed containing the most recent 10 articles.

Installation
------------

```sh
go install github.com/rcrowley/feed@latest
```

Usage
-----

```sh
feed -a <author> [-o <output>] -t <title> -u <url> [-v] [-x <exclude>[...]] [<docroot>[...]]
```

* `-a <author>`: author's name
* `-o <output>`: write to this file instead of standard output
* `-t <title>`: feed title
* `-u <url>`: site URL with scheme and domain
* `-v`: verbose mode
* `-x <exclude>`: subdirectory of <docroot> to exclude (may be repeated)
* `<docroot>`: document root directory to scan (defaults to the current working directory; may be repeated)

See also
--------

Feed is part of the [Mergician](https://github.com/rcrowley/mergician) suite of tools that manipulate HTML documents:

* [Critaique](https://github.com/rcrowley/critaique): Prompt an LLM to make suggestions for improving your writing, ask follow-up questions, etc.
* [Deadlinks](https://github.com/rcrowley/deadlinks): Scan a document root directory for dead links
* [Electrostatic](https://github.com/rcrowley/electrostatic): Mergician-powered, pure-HTML CMS
* [Frag](https://github.com/rcrowley/frag): Extract fragments of HTML documents
* [Sitesearch](https://github.com/rcrowley/sitesearch): Index a document root directory and serve queries to it in AWS Lambda
