# robotstxt
Package robotstxt implements the Robots Exclusion Protocol, https://en.wikipedia.org/wiki/Robots_exclusion_standard, with a simple API.

[![Go Report Card](https://goreportcard.com/badge/github.com/itmayziii/robotstxt)](https://goreportcard.com/report/github.com/itmayziii/robotstxt)
[![](https://godoc.org/github.com/itmayziii/robotstxt?status.svg)](https://godoc.org/github.com/itmazyiii/robotstxt)

Link to the GoDocs -> [here](https://godoc.org/github.com/itmayziii/robotstxt).


## Basic Example
```go
    package main

    import (
    	"fmt"
    	"github.com/itmayziii/robotstxt"
    )

    func main () {
    	robotsTxt, _ := robotstxt.New("", `
# Robots.txt test file
# 06/04/2018
    # Indented comments are allowed
        
User-agent : *
Crawl-delay: 5
Disallow: /cms/
Disallow: /pricing/frontend
Disallow: /pricing/admin/ # SPA application built into the site
Disallow : *?s=lightbox
Disallow: /se/en$
Disallow:*/retail/*/frontend/*
        
Allow: /be/fr_fr/retail/fr/
        
# Multiple groups with all access
User-agent: AdsBot-Google
User-agent: AdsBot-Bing
Allow: /
        
# Multiple sitemaps
Sitemap: https://www.dumpsters.com/sitemap.xml
Sitemap: https://www.dumpsters.com/sitemap-launch-index.xml
`)
        canCrawl, err := robotsTxt.CanCrawl("googlebot", "/cms/pages")
        fmt.Println(canCrawl)
        fmt.Println(err)
        // Output:
        // false
        // <nil>
    }
```

## Specification

A large portion of how this package handles the specification comes from https://developers.google.com/search/reference/robots_txt.
In fact this package tests against all of the examples listed at
https://developers.google.com/search/reference/robots_txt#url-matching-based-on-path-values plus many more.

Important Notes From the Spec

1. User Agents are case insensitive so "googlebot" and "Googlebot" are the same thing.

2. Directive "Allow" and "Disallow" values are case sensitive so "/pricing" and "/Pricing" are not the same thing.

3. The entire file must be valid UTF-8 encoded, this package will return an error if that is not the case.

4. The most specific user agent wins.

5. Allow and disallow directives also respect the one that is most specific and in the event of a tie the allow directive will win.

6. Directives listed in the robots.txt file apply only to a host, protocol, and port number,
https://developers.google.com/search/reference/robots_txt#file-location--range-of-validity. This package validates the host, protocol,
and port number every time it is asked if a robot "CanCrawl" a path and the path contains the host, protocol, and port.
```go
 robotsTxt := robotstxt.New("https://www.dumpsters.com", `
     User-agent: *
     Disallow: "/wiki/"
 `)
 robotsTxt.CanCrawl("googlebot", "/products/") // True
 robotsTxt.CanCrawl("googlebot", "https://www.dumpsters.com/products/") // True
 robotsTxt.CanCrawl("googlebot", "http://www.dumpsters.com/products/") // False - the URL did not match the URL provided when "robotsTxt" was created
```
