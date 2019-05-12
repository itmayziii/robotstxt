/*
Package robotstxt implements the Robots Exclusion Protocol, https://en.wikipedia.org/wiki/Robots_exclusion_standard, with a simple API.

Specification

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
 robotsTxt := robotstxt.New("https://www.dumpsters.com", `
     User-agent: *
     Disallow: "/wiki/"
 `)
 robotsTxt.CanCrawl("googlebot", "/products/") // True
 robotsTxt.CanCrawl("googlebot", "https://www.dumpsters.com/products/") // True
 robotsTxt.CanCrawl("googlebot", "http://www.dumpsters.com/products/") // False - the URL did not match the URL provided when "robotsTxt" was created
*/
package robotstxt

import (
	"bytes"
	"errors"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	netUrl "net/url"
	"os"
	"strings"
	"time"
)

var httpGet = http.Get

// RobotsExclusionProtocol exposes all of the things you would want to know about a robots.txt file without giving direct access to the directives
// defined. Directives such as allow and disallow are not important for a robot (user-agent) to know about, they are implementation details,
// instead a robot just needs to know if it is allowed to crawl a given path so this interface provides a "CanCrawl" method as opposed to giving you
// direct access to allow and disallow.
type RobotsExclusionProtocol interface {
	// CanCrawl determines whether or not a given robot (user-agent) is allowed to crawl a URL based on allow and disallow directives in the
	// robots.txt.
	CanCrawl(robotName, url string) (bool, error)
	// Returns the sitemaps that are defined in the robots.txt.
	Sitemaps() []string
	// Getter that returns the URL a particular robots.txt file is associated with, i.e. https://www.dumpsters.com.
	URL() string
	// How long should a robot wait between accessing pages on a site.
	CrawlDelay(robotName string) time.Duration
}

// ProtocolResult is used for concurrent operations such as NewFromFile and NewFromURL.
type ProtocolResult struct {
	Protocol RobotsExclusionProtocol
	Error    error
}

// New creates an implementation of RobotsExclusionProtocol.
func New(url, robotsTxtContent string) (RobotsExclusionProtocol, error) {
	reader := strings.NewReader(robotsTxtContent)
	return parse(url, reader)
}

// NewFromFile creates an implementation of RobotsExclusionProtocol from a local file.
func NewFromFile(url, filePath string, ch chan ProtocolResult) {
	defer close(ch)
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		ch <- ProtocolResult{Protocol: robotsTxt{}, Error: err}
		return
	}
	defer safeClose(file)

	robotsTxt, err := parse(url, file)
	ch <- ProtocolResult{Protocol: robotsTxt, Error: err}
}

/*
NewFromURL retrieves a robots.txt for a given scheme, host, and an optional port number. According to the spec the robots.txt file must always live
at the top level directory, https://developers.google.com/search/reference/robots_txt#file-location--range-of-validity,
so everything that is not the top level is ignored.

The following are examples of only looking at the top level for /robots.txt:
  Given:                                                  Looks for:
  https://www.dumpsters.com/pricing/roll-off-dumpsters -> https://www.dumpsters.com/robots.txt
  https://www.dumpsters.com                            -> https://www.dumpsters.com/robots.txt
  https://www.dumpsters.com/robots.txt                 -> https://www.dumpsters.com/robots.txt
*/
func NewFromURL(url string, ch chan ProtocolResult) {
	defer close(ch)
	parsedUrl, err := netUrl.Parse(url)
	if err != nil {
		log.Println(err)
		ch <- ProtocolResult{Protocol: robotsTxt{}, Error: err}
		return
	}

	normalizedUrl := parsedUrl.Scheme + "://" + parsedUrl.Hostname()
	port := parsedUrl.Port()
	if port != "" {
		normalizedUrl = parsedUrl.Scheme + "://" + parsedUrl.Hostname() + ":" + port
	}

	resp, err := httpGet(normalizedUrl + "/robots.txt")
	if err != nil {
		log.Println(err)
		ch <- ProtocolResult{Protocol: robotsTxt{}, Error: err}
		return
	}
	defer safeClose(resp.Body)

	robotsTxtBody, err := parseRobotsTxtBody(resp.Body)
	if err != nil {
		log.Println(err)
		ch <- ProtocolResult{Protocol: robotsTxt{}, Error: err}
		return
	}

	robotsTxt, err := New(url, robotsTxtBody)
	ch <- ProtocolResult{Protocol: robotsTxt, Error: err}
}

type robotsTxt struct {
	robots   map[string]robot
	sitemaps []string
	url      string
}

type robot struct {
	disallowed []string
	allowed    []string
	crawlDelay time.Duration
}

func (robotsTxt robotsTxt) CanCrawl(robotName, url string) (bool, error) {
	robot, exists := findMatchingRobot(robotName, robotsTxt.robots)
	if !exists {
		return true, nil
	}

	// Everything is allowed if nothing is disallowed.
	if robot.disallowed == nil || len(robot.disallowed) == 0 {
		return true, nil
	}

	// URL provided must be able to be parsed.
	parsedUrl, err := netUrl.Parse(url)
	if err != nil {
		return true, err
	}

	// Basically if the URL provided is a full URL with a schema then the robot URL must match completely.
	// https://developers.google.com/search/reference/robots_txt#file-location--range-of-validity
	if parsedUrl.IsAbs() {
		normalizedUrl, err := normalizeUrl(parsedUrl.String())
		if err != nil {
			log.Println(err)
		}
		if robotsTxt.url != normalizedUrl {
			return true, errors.New("absolute URL provided but the robot URL did not match")
		}
	}

	// Prepend a leading slash if the url provided does not have one, just one less thing we have to account for later on
	normalizedPath := parsedUrl.RequestURI()
	if !strings.HasPrefix(normalizedPath, "/") {
		normalizedPath = "/" + normalizedPath
	}

	// With allow and disallow directives, the most specific rule based on the length of the [path] entry will trump the less specific (shorter) rule.
	// https://developers.google.com/search/reference/robots_txt#url-matching-based-on-path-values
	disallowedLength, err := urlMatchLength(normalizedPath, robot.disallowed)
	if err != nil {
		return true, err
	}
	allowedLength, err := urlMatchLength(normalizedPath, robot.allowed)
	if err != nil {
		return true, err
	}
	return disallowedLength == 0 || allowedLength >= disallowedLength, nil
}

func (robotsTxt robotsTxt) CrawlDelay(robotName string) time.Duration {
	robot, _ := findMatchingRobot(robotName, robotsTxt.robots)
	return robot.crawlDelay
}

func (robotsTxt robotsTxt) Sitemaps() []string {
	return robotsTxt.sitemaps
}

func (robotsTxt robotsTxt) URL() string {
	return robotsTxt.url
}

func parseRobotsTxtBody(readCloser io.ReadCloser) (string, error) {
	node, err := html.Parse(readCloser)
	if err != nil {
		return "", err
	}

	body, err := getBody(node)
	if err != nil {
		return "", err
	}

	bodyString, err := renderNode(body)
	if err != nil {
		return "", err
	}

	bodyTokenizer := html.NewTokenizer(strings.NewReader(bodyString))
	bodyText := ""
Loop:
	for {
		tt := bodyTokenizer.Next()
		switch tt {
		case html.ErrorToken:
			break Loop
		case html.TextToken:
			bodyText += string(bodyTokenizer.Text())
		}

	}

	return bodyText, nil
}

func getBody(doc *html.Node) (*html.Node, error) {
	var body *html.Node
	var parseHtmlForBody func(*html.Node)
	hasMatch := false
	parseHtmlForBody = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "body" {
			body = node
			hasMatch = true
		}
		if !hasMatch {
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				parseHtmlForBody(c)
			}
		}
	}
	parseHtmlForBody(doc)

	if body == nil {
		return nil, errors.New("missing <body> in the node tree")
	}
	return body, nil
}

func renderNode(n *html.Node) (string, error) {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	err := html.Render(w, n)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
