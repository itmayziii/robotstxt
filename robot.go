package main

import (
	"errors"
	netUrl "net/url"
	"strings"
)

// The directives listed in the robotstxt.txt file apply only to the host, protocol and port number where the file is hosted.

// An optional Unicode BOM (byte order mark) at the beginning of the robotstxt.txt file is ignored.

type Crawler interface {
	CanCrawl(url string) (bool, error)
	Sitemaps() []string
}

type Robot struct {
	name       string
	url        string
	disallowed []string
	allowed    []string
}

func NewRobot(name string) Robot {
	return Robot{name: name}
}

func NewRobotForUrl(name string, url string) Robot {
	return Robot{name: name, url: url}
}

// Will return true by default since robots are allowed to crawl unless specifically told not to
func (robot Robot) CanCrawl(url string) (bool, error) {
	if robot.disallowed == nil || len(robot.disallowed) == 0 { // Everything is allowed if nothing is disallowed
		return true, nil
	}

	parsedUrl, err := netUrl.Parse(url) // Url provided must be able to be parsed
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
	normalizedPath := parsedUrl.Path
	if !strings.HasPrefix(normalizedPath, "/") {
		normalizedPath = "/" + normalizedPath
	}

	// Allow and disallow directives, the most specific rule based on the length of the [path] entry will trump the less specific (shorter) rule
	// (https://developers.google.com/search/reference/robots_txt#url-matching-based-on-path-values)
	disallowedLength := urlMatchLength(normalizedPath, robot.disallowed)
	allowedLength := urlMatchLength(normalizedPath, robot.allowed)
	return disallowedLength == 0 || allowedLength >= disallowedLength, nil
}

func urlMatchLength(url string, paths []string) int {
	if paths == nil || len(paths) == 0 {
		return 0
	}

	matchLength := 0
	for _, path := range paths {
		// Ending with a trailing "*" is the same as not doing it at all | lowercase to make it case insensitive per the spec
		path = strings.ToLower(strings.TrimSuffix(path, "*"))

		//if strings.Contains(path, "*") {
		//
		//}

		if strings.HasPrefix(url, path) {
			matchLength = len(path)
			break
		}
	}

	return matchLength
}
