package robotstxt

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"
)

func parse(url string, reader io.Reader) (RobotsExclusionProtocol, error) {
	robotsTxt := robotsTxt{}
	robots := make(map[string]robot)

	currentUserAgents := make([]string, 0)
	endUserAgents := false
	lineScanner := bufio.NewScanner(reader)
	lineScanner.Split(bufio.ScanLines)
	for lineNumber := 1; lineScanner.Scan(); lineNumber++ {
		// Everything about the a robots.txt is case insensitive so we just make everything lowercase for easier comparisons.
		line := strings.TrimSpace(lineScanner.Text())

		if utf8.ValidString(line) == false {
			err := errors.New("invalid encoding detected on line " + strconv.Itoa(lineNumber) + ", all characters must be UTF-8 encoded")
			log.Println(err)
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
				robot.allowed = append(robots[userAgent].allowed, value)
				robots[userAgent] = robot
			}
			endUserAgents = true
			break
		case "disallow":
			for _, userAgent := range currentUserAgents {
				robot := robots[userAgent]
				robot.disallowed = append(robots[userAgent].disallowed, value)
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
				log.Println("invalid crawl-delay, could not convert " + value + " to an integer")
				return robotsTxt, err
			}
			for _, userAgent := range currentUserAgents {
				robot := robots[userAgent]
				robot.crawlDelay = valueInt
				robots[userAgent] = robot
			}
			endUserAgents = true
			break
		}
	}

	robotsTxt.url = url
	robotsTxt.robots = robots
	return robotsTxt, nil
}
