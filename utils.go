package robotstxt

import (
	"errors"
	"io"
	"log"
	"regexp"
	"strings"
)

func safeClose(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Println(err)
	}
}

func urlMatchLength(url string, paths []string) (int, error) {
	if paths == nil || len(paths) == 0 {
		return 0, nil
	}

	matchLength := 0
	for _, path := range paths {
		// Handle the wildcards.
		if strings.Contains(path, "*") || strings.Contains(path, "$") {
			expression := strings.Replace(path, "*", "(.*)", -1)
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

func findMatchingRobot(robotName string, robots map[string]robot) (robot, bool) {
	// User agents are case insensitive.
	// https://developers.google.com/search/reference/robots_txt#order-of-precedence-for-user-agents
	robotName = strings.ToLower(robotName)

	robotNames := keys(robots)
	matchedRobotName := ""
	for _, name := range robotNames {
		if strings.HasPrefix(robotName, strings.ToLower(name)) && len(name) >= len(matchedRobotName) {
			matchedRobotName = name
		}
	}

	if matchedRobotName != "" {
		return robots[matchedRobotName], true
	}

	// If we made it this far then there is no matching robot, let's check for the wildcard.
	allUserAgents, exists := robots["*"]
	if !exists {
		return robot{}, false
	}

	return allUserAgents, true
}

func keys(robots map[string]robot) []string {
	keys := make([]string, len(robots))

	i := 0
	for k := range robots {
		keys[i] = k
		i++
	}

	return keys
}
