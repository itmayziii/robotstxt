// Test cases derived from https://developers.google.com/search/reference/robots_txt#url-matching-based-on-path-values
package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRobot_CanCrawl_allows_all_if_nothing_is_disallowed(t *testing.T) {
	robot := Robot{}
	testRobot(t, robot, []testUrl{
		{url: "/", crawlable: true, hasError: false},                             // root path
		{url: "/pricing", crawlable: true, hasError: false},                      // path
		{url: "/contact-us", crawlable: true, hasError: false},                   // path, mixed spelling
		{url: "/pricing?id=123", crawlable: true, hasError: false},               // query params
		{url: "/pricing/product", crawlable: true, hasError: false},              // nested path
		{url: "/contact/more-information.php", crawlable: true, hasError: false}, // nested path, mixed spelling
		{url: "/pricing/product/sale", crawlable: true, hasError: false},         // deeply nested path
		{url: "/pricing.html", crawlable: true, hasError: false},                 // file extension
		{url: "/pricing.php?id=123", crawlable: true, hasError: false},           // file extension with query param
		{url: "pricing/test", crawlable: true, hasError: false},                  // relative path
	})
}

func TestRobot_CanCrawl_allows_nothing_if_root_is_disallowed(t *testing.T) {
	robot := Robot{disallowed: []string{"/"}}
	testRobot(t, robot, []testUrl{
		{url: "/", crawlable: false, hasError: false},                             // root path
		{url: "/pricing", crawlable: false, hasError: false},                      // path
		{url: "/contact-us", crawlable: false, hasError: false},                   // path, mixed spelling
		{url: "/pricing?id=123", crawlable: false, hasError: false},               // query params
		{url: "/pricing/product", crawlable: false, hasError: false},              // nested path
		{url: "/contact/more-information.php", crawlable: false, hasError: false}, // nested path, mixed spelling
		{url: "/pricing/product/sale", crawlable: false, hasError: false},         // deeply nested path
		{url: "/pricing.html", crawlable: false, hasError: false},                 // file extension
		{url: "/pricing.php?id=123", crawlable: false, hasError: false},           // file extension with query param
		{url: "pricing/test", crawlable: false, hasError: false},                  // relative path
	})
}

func TestRobot_CanCrawl_allows_nothing_if_root_is_nothing_with_wildcard(t *testing.T) {
	robot := Robot{disallowed: []string{"/*"}}
	testRobot(t, robot, []testUrl{
		{url: "/", crawlable: false, hasError: false},                             // root path
		{url: "/pricing", crawlable: false, hasError: false},                      // path
		{url: "/contact-us", crawlable: false, hasError: false},                   // path, mixed spelling
		{url: "/pricing?id=123", crawlable: false, hasError: false},               // query params
		{url: "/pricing/product", crawlable: false, hasError: false},              // nested path
		{url: "/contact/more-information.php", crawlable: false, hasError: false}, // nested path, mixed spelling
		{url: "/pricing/product/sale", crawlable: false, hasError: false},         // deeply nested path
		{url: "/pricing.html", crawlable: false, hasError: false},                 // file extension
		{url: "/pricing.php?id=123", crawlable: false, hasError: false},           // file extension with query param
		{url: "pricing/test", crawlable: false, hasError: false},                  // relative path
	})
}

func TestRobot_CanCrawl_normal_path(t *testing.T) {
	robot := Robot{disallowed: []string{"/pricing"}}
	testRobot(t, robot, []testUrl{
		{url: "/", crawlable: true, hasError: false},                             // root path
		{url: "/pricing", crawlable: false, hasError: false},                     // path
		{url: "/contact-us", crawlable: true, hasError: false},                   // path, mixed spelling
		{url: "/pricing?id=123", crawlable: false, hasError: false},              // query params
		{url: "/pricing/product", crawlable: false, hasError: false},             // nested path
		{url: "/contact/more-information.php", crawlable: true, hasError: false}, // nested path, mixed spelling
		{url: "/pricing/product/sale", crawlable: false, hasError: false},        // deeply nested path
		{url: "/pricing.html", crawlable: false, hasError: false},                // file extension
		{url: "/pricing.php?id=123", crawlable: false, hasError: false},          // file extension with query param
		{url: "pricing/test", crawlable: false, hasError: false},                 // relative path
	})
}

func TestRobot_CanCrawl_short_path(t *testing.T) {
	robot := Robot{disallowed: []string{"/p"}}
	testRobot(t, robot, []testUrl{
		{url: "/", crawlable: true, hasError: false},                             // root path
		{url: "/pricing", crawlable: false, hasError: false},                     // path
		{url: "/contact-us", crawlable: true, hasError: false},                   // path, mixed spelling
		{url: "/pricing?id=123", crawlable: false, hasError: false},              // query params
		{url: "/pricing/product", crawlable: false, hasError: false},             // nested path
		{url: "/contact/more-information.php", crawlable: true, hasError: false}, // nested path, mixed spelling
		{url: "/pricing/product/sale", crawlable: false, hasError: false},        // deeply nested path
		{url: "/pricing.html", crawlable: false, hasError: false},                // file extension
		{url: "/pricing.php?id=123", crawlable: false, hasError: false},          // file extension with query param
		{url: "pricing/test", crawlable: false, hasError: false},                 // relative path
	})
}

func TestRobot_CanCrawl_nested_path(t *testing.T) {
	robot := Robot{disallowed: []string{"/pricing/product"}}
	testRobot(t, robot, []testUrl{
		{url: "/", crawlable: true, hasError: false},                             // root path
		{url: "/pricing", crawlable: true, hasError: false},                      // path
		{url: "/contact-us", crawlable: true, hasError: false},                   // path, mixed spelling
		{url: "/pricing?id=123", crawlable: true, hasError: false},               // query params
		{url: "/pricing/product", crawlable: false, hasError: false},             // nested path
		{url: "/contact/more-information.php", crawlable: true, hasError: false}, // nested path, mixed spelling
		{url: "/pricing/product/sale", crawlable: false, hasError: false},        // deeply nested path
		{url: "/pricing.html", crawlable: true, hasError: false},                 // file extension
		{url: "/pricing.php?id=123", crawlable: true, hasError: false},           // file extension with query param
		{url: "pricing/test", crawlable: true, hasError: false},                  // relative path
	})
}

func TestRobot_CanCrawl_multiple_disallowed(t *testing.T) {
	robot := Robot{disallowed: []string{"/pricing", "/contact-us"}}
	testRobot(t, robot, []testUrl{
		{url: "/", crawlable: true, hasError: false},                             // root path
		{url: "/pricing", crawlable: false, hasError: false},                     // path
		{url: "/contact-us", crawlable: false, hasError: false},                  // path, mixed spelling
		{url: "/pricing?id=123", crawlable: false, hasError: false},              // query params
		{url: "/pricing/product", crawlable: false, hasError: false},             // nested path
		{url: "/contact/more-information.php", crawlable: true, hasError: false}, // nested path, mixed spelling
		{url: "/pricing/product/sale", crawlable: false, hasError: false},        // deeply nested path
		{url: "/pricing.html", crawlable: false, hasError: false},                // file extension
		{url: "/pricing.php?id=123", crawlable: false, hasError: false},          // file extension with query param
		{url: "pricing/test", crawlable: false, hasError: false},                 // relative path
	})
}

func TestRobot_CanCrawl_allowed_overrides_disallowed_when_allowed_has_greater_length(t *testing.T) {
	robot := Robot{disallowed: []string{"/p"}, allowed: []string{"/pric"}}
	testRobot(t, robot, []testUrl{
		{url: "/", crawlable: true, hasError: false},                             // root path
		{url: "/pricing", crawlable: true, hasError: false},                      // path
		{url: "/contact-us", crawlable: true, hasError: false},                   // path, mixed spelling
		{url: "/pricing?id=123", crawlable: true, hasError: false},               // query params
		{url: "/pricing/product", crawlable: true, hasError: false},              // nested path
		{url: "/contact/more-information.php", crawlable: true, hasError: false}, // nested path, mixed spelling
		{url: "/pricing/product/sale", crawlable: true, hasError: false},         // deeply nested path
		{url: "/pricing.html", crawlable: true, hasError: false},                 // file extension
		{url: "/pricing.php?id=123", crawlable: true, hasError: false},           // file extension with query param
		{url: "pricing/test", crawlable: true, hasError: false},                  // relative path
	})
}

func TestRobot_CanCrawl_allowed_overrides_disallowed_when_allowed_has_equal_length(t *testing.T) {
	robot := Robot{disallowed: []string{"/con"}, allowed: []string{"/con"}}
	testRobot(t, robot, []testUrl{
		{url: "/", crawlable: true, hasError: false},                             // root path
		{url: "/pricing", crawlable: true, hasError: false},                      // path
		{url: "/contact-us", crawlable: true, hasError: false},                   // path, mixed spelling
		{url: "/pricing?id=123", crawlable: true, hasError: false},               // query params
		{url: "/pricing/product", crawlable: true, hasError: false},              // nested path
		{url: "/contact/more-information.php", crawlable: true, hasError: false}, // nested path, mixed spelling
		{url: "/pricing/product/sale", crawlable: true, hasError: false},         // deeply nested path
		{url: "/pricing.html", crawlable: true, hasError: false},                 // file extension
		{url: "/pricing.php?id=123", crawlable: true, hasError: false},           // file extension with query param
		{url: "pricing/test", crawlable: true, hasError: false},                  // relative path
	})
}

func TestRobot_CanCrawl_disallowed_overrides_allowed_when_disallowed_has_greater_length(t *testing.T) {
	robot := Robot{disallowed: []string{"/cont"}, allowed: []string{"/con"}}
	testRobot(t, robot, []testUrl{
		{url: "/", crawlable: true, hasError: false},                              // root path
		{url: "/pricing", crawlable: true, hasError: false},                       // path
		{url: "/contact-us", crawlable: false, hasError: false},                   // path, mixed spelling
		{url: "/pricing?id=123", crawlable: true, hasError: false},                // query params
		{url: "/pricing/product", crawlable: true, hasError: false},               // nested path
		{url: "/contact/more-information.php", crawlable: false, hasError: false}, // nested path, mixed spelling
		{url: "/pricing/product/sale", crawlable: true, hasError: false},          // deeply nested path
		{url: "/pricing.html", crawlable: true, hasError: false},                  // file extension
		{url: "/pricing.php?id=123", crawlable: true, hasError: false},            // file extension with query param
		{url: "pricing/test", crawlable: true, hasError: false},                   // relative path
	})
}

type testUrl struct {
	url       string
	crawlable bool
	hasError  bool
}

// I know it's bad to write code for tests, but testing each thing was painful and this will be consistent / less human error prone
func testRobot(t *testing.T, robot Robot, testUrls []testUrl) {
	for _, testUrl := range testUrls {
		canCrawl, err := robot.CanCrawl(testUrl.url)

		if testUrl.crawlable {
			assert.True(t, canCrawl, testUrl)
		} else {
			assert.False(t, canCrawl, testUrl)
		}

		if testUrl.hasError {
			assert.NotNil(t, err, testUrl)
		} else {
			assert.Nil(t, err, testUrl)
		}
	}
}
