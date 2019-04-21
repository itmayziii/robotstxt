/*
Package robotstxt gives a very simple API for determining if a [path] can be crawled by a robot.
Implementation of the robots exclusion protocol https://en.wikipedia.org/wiki/Robots_exclusion_standard.
*/
package robotstxt

import (
	"errors"
	netUrl "net/url"
	"regexp"
	"strings"
)

type Crawler interface {
	CanCrawl(url string) (bool, error)
	CrawlDelay() int
	Sitemaps() []string
}

type Robot struct {
	name       string
	url        string
	disallowed []string
	allowed    []string
	sitemaps   []string
	crawlDelay int
}

func NewRobot(name, url string, disallowed, allowed, sitemaps []string, crawlDelay int) Robot {
	return Robot{name: name, url: url, disallowed: disallowed, allowed: allowed}
}

func (robot Robot) CanCrawl(url string) (bool, error) {
	// Everything is allowed if nothing is disallowed
	if robot.disallowed == nil || len(robot.disallowed) == 0 {
		return true, nil
	}

	// Url provided must be able to be parsed
	parsedUrl, err := netUrl.Parse(url)
	if err != nil {
		return true, err
	}

	// Basically if the url provided is a full URL with a schema then the robot must be told what url it is working with
	// https://developers.google.com/search/reference/robots_txt#file-location--range-of-validity
	if parsedUrl.IsAbs() && robot.url == "" {
		return true, errors.New("absolute URL provided but the robot was not given a URL to validate against")
	}

	//absoluteUrl := parsedUrl.Scheme + "://" + parsedUrl.Host
	//if robot.url != absoluteUrl {
	//	return true, errors.New("absolute URL provided, " + absoluteUrl + ", but it does not match the current robot.url, " + robot.url)
	//}

	// Prepend a leading slash if the url provided does not have one, just one less thing we have to account for later on
	normalizedPath := parsedUrl.RequestURI()
	if !strings.HasPrefix(normalizedPath, "/") {
		normalizedPath = "/" + normalizedPath
	}

	// Allow and disallow directives, the most specific rule based on the length of the [path] entry will trump the less specific (shorter) rule
	// (https://developers.google.com/search/reference/robots_txt#url-matching-based-on-path-values)
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

func urlMatchLength(url string, paths []string) (int, error) {
	if paths == nil || len(paths) == 0 {
		return 0, nil
	}

	matchLength := 0
	for _, path := range paths {
		// Handle wildcards
		if strings.Contains(path, "*") || strings.Contains(path, "$") {
			expression := strings.Replace(path, "*", "(.*)", 1)
			regExp, err := regexp.Compile(expression)
			if err != nil {
				return 0, errors.New("unable to get length of path " + path)
			}
			match := regExp.FindString(url)
			if match == "" {
				continue
			}

			matchLength = len(path)
			break
		}

		if strings.HasPrefix(url, path) {
			matchLength = len(path)
			break
		}
	}

	return matchLength, nil
}
