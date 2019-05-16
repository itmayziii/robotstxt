# robotstxt
Package robotstxt implements the Robots Exclusion Protocol, https://en.wikipedia.org/wiki/Robots_exclusion_standard, with a simple API.
This repo also exclusively uses [Go Modules](https://github.com/golang/go/wiki/Modules).

[![Go Report Card](https://goreportcard.com/badge/github.com/itmayziii/robotstxt)](https://goreportcard.com/report/github.com/itmayziii/robotstxt)
[![](https://godoc.org/github.com/itmayziii/robotstxt?status.svg)](https://godoc.org/github.com/itmayziii/robotstxt)
[![Coverage Status](https://coveralls.io/repos/github/itmayziii/robotstxt/badge.svg?branch=master)](https://coveralls.io/github/itmayziii/robotstxt?branch=master)

Link to the GoDocs -> [here](https://godoc.org/github.com/itmayziii/robotstxt).


## Basic Examples

### 1. Creating a robotsTxt with a URL
This is the most common way to use this package since most robots.txt files you will be interested in will be on a server somewhere. This library 
gives you the freedom to specify the `Get` method you want to use to make the HTTP request. This is useful for people that may want to use their 
own `http.Client`. 
```go
package main

import (
    "fmt"
	"github.com/itmayziii/robotstxt/v2"
    "net/http"
)

func main () {
	robotsTxt, err := robotstxt.NewFromURL("https://www.dumpsters.com", http.Get)
	fmt.Println(err)

	canCrawl, err := robotsTxt.CanCrawl("googlebot", "/bdso/pages")
	fmt.Println(canCrawl)
	fmt.Println(err)
	fmt.Println(robotsTxt.Sitemaps())
	fmt.Println(robotsTxt.URL())
	fmt.Println(robotsTxt.CrawlDelay("googlebot"))
	// <nil>
	// false
	// <nil>
	// [https://www.dumpsters.com/sitemap.xml https://www.dumpsters.com/sitemap-launch-index.xml]
	// https://www.dumpsters.com:443
	// 5s
}
```

### 2. Creating a robotsTxt Manually
You likely will not be doing this method as you would need to parse get the robots.txt from the server yourself.
```go
package main

import (
    "fmt"
	"github.com/itmayziii/robotstxt/v2"
    "strings"
)

func main () {
    robotsTxt, _ := robotstxt.New("", strings.NewReader(`
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
`))
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

### Important Notes From the Spec

1. User Agents are case insensitive so "googlebot" and "Googlebot" are the same thing.

2. Directive "Allow" and "Disallow" values are case sensitive so "/pricing" and "/Pricing" are not the same thing.

3. The entire file must be valid UTF-8 encoded, this package will return an error if that is not the case.

4. The most specific user agent wins.

5. Allow and disallow directives also respect the one that is most specific based on length and in the event of a tie the allow directive will win, 
i.e. `disallow: /cms/` loses to `allow: /cms/` and to `allow: /cms*` but not to `allow: /cms`.

6. Directives listed in the robots.txt file apply only to a host, protocol, and port number,
https://developers.google.com/search/reference/robots_txt#file-location--range-of-validity. This package validates the host, protocol,
and port number every time it is asked if a robot "CanCrawl" a path and the path contains the host, protocol, and port.
```go
 robotsTxt := robotstxt.New("https://www.dumpsters.com", strings.NewReader(`
     User-agent: *
     Disallow: "/wiki/"
 `))
 robotsTxt.CanCrawl("googlebot", "/products/") // True
 robotsTxt.CanCrawl("googlebot", "https://www.dumpsters.com/products/") // True
 robotsTxt.CanCrawl("googlebot", "http://www.dumpsters.com/products/") // False - the URL did not match the URL provided when "robotsTxt" was created
```

## Roadmap
* Respect a "noindex" meta tag and HTTP response header as described [here](https://en.wikipedia.org/wiki/Robots_exclusion_standard#Meta_tags_and_headers).
 There a couple of considerations to be taken into account before implementing this:
  * We need to leave the current `CanCrawl` method as is since it is meant to determine whether or not a robot can crawl a page prior to actually 
  loading the page. The "noindex" and meta tag and HTTP response header by nature of where they are located only happen after the crawler has 
  loaded the page.
  * Maybe 2 methods would be needed to implement this. One method that would retrieve the response for the user and hand back an instance of 
  `RobotsExclusionProtocol` as well as the response itself, something like `CanCrawlPage`, which of course would also go through the robots.txt 
  logic before even requesting the page. A separate second method that would take an already retrieved response that the user has goes through the 
  same logic that the first method does. I'm not 100% sure we would need both methods but I can see why some people would want to retrieve the HTTP
  response themselves.
* Potentially support the Host directive as described [here](https://en.wikipedia.org/wiki/Robots_exclusion_standard#Host).  
