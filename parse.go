package robotstxt

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var validateUTF8 = utf8.ValidString

func parse(url string, reader io.Reader) (RobotsTxt, error) {
	normalizedUrl, err := normalizeUrl(url)
	if err != nil {
		return RobotsTxt{}, err
	}

	robotsTxt := RobotsTxt{}
	robots := make(map[string]robot)
	currentUserAgents := make([]string, 1) // User agents that are part of the same group.
	endUserAgents := false                 // Are we still processing user agents as part of the same group.
	lineScanner := bufio.NewScanner(reader)
	lineScanner.Split(bufio.ScanLines)
	for lineNumber := 1; lineScanner.Scan(); lineNumber++ {
		line := strings.TrimSpace(lineScanner.Text())

		if validateUTF8(line) == false {
			err := errors.New("invalid encoding detected on line " + strconv.Itoa(lineNumber) + ", all characters must be UTF-8 encoded")
			return robotsTxt, err
		}

		// Only process the text before any comment.
		line = strings.Split(line, "#")[0]

		// The entire line is a comment.
		if line == "" {
			continue
		}

		// Check for separator between key and value.
		if !strings.Contains(line, ":") {
			continue
		}

		separateKeyValue := strings.Split(line, ":")
		key := strings.ToLower(strings.TrimSpace(separateKeyValue[0]))
		value := strings.TrimSpace(strings.Join(separateKeyValue[1:], ":"))
		// Another faulty key value pair.
		if key == "" || value == "" {
			continue
		}
		// A value can only be one word, literally ignore anything more than that.
		value = strings.Split(value, " ")[0]

		switch key {
		case "user-agent":
			if endUserAgents {
				currentUserAgents = []string{}
				endUserAgents = false
			}
			currentUserAgents = append(currentUserAgents, value)
			robots[value] = robot{}
			break
		case "allow":
			for _, userAgent := range currentUserAgents {
				robot := robots[userAgent]
				robot.allowed = append(robot.allowed, value)
				robots[userAgent] = robot
			}
			endUserAgents = true
			break
		case "disallow":
			for _, userAgent := range currentUserAgents {
				robot := robots[userAgent]
				robot.disallowed = append(robot.disallowed, value)
				robots[userAgent] = robot
			}
			endUserAgents = true
			break
		case "sitemap":
			robotsTxt.sitemaps = append(robotsTxt.sitemaps, value)
			endUserAgents = true
			break
		case "crawl-delay":
			valueInt, err := strconv.Atoi(value)
			if err != nil {
				return robotsTxt, err
			}
			for _, userAgent := range currentUserAgents {
				robot := robots[userAgent]
				robot.crawlDelay = time.Duration(valueInt) * time.Second
				robots[userAgent] = robot
			}
			endUserAgents = true
			break
		}
	}

	robotsTxt.url = normalizedUrl
	robotsTxt.robots = robots
	return robotsTxt, nil
}
