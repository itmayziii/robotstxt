package robotstxt_test

import (
	"fmt"
	"github.com/itmayziii/robotstxt/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Test cases derived from https://developers.google.com/search/reference/robots_txt#url-matching-based-on-path-values.
func Test_examples_mentioned_in_google_spec(t *testing.T) {
	// Matches the root and any lower level URL.
	t.Run("/", func(t *testing.T) {
		robotsTxt, err := robotstxt.New("https://www.dumpsters.com", strings.NewReader(`
User-Agent: *
Disallow: /
`))
		assert.Nil(t, err)

		testRobot(t, "Bingbot", robotsTxt, []testUrl{
			{url: "/", crawlable: false, hasError: false},
			{url: "", crawlable: false, hasError: false},
			{url: "/anything", crawlable: false, hasError: false},
			{url: "/anything?test=1", crawlable: false, hasError: false},
			{url: "/anything/else", crawlable: false, hasError: false},
		})
	})

	// Equivalent to /. The trailing wildcard is ignored
	t.Run("/*", func(t *testing.T) {
		robotsTxt, err := robotstxt.New("https://www.dumpsters.com", strings.NewReader(`
User-Agent: *
Disallow: /*
`))
		assert.Nil(t, err)

		testRobot(t, "Bingbot", robotsTxt, []testUrl{
			{url: "/", crawlable: false, hasError: false},
			{url: "", crawlable: false, hasError: false},
			{url: "/anything", crawlable: false, hasError: false},
			{url: "/anything?test=1", crawlable: false, hasError: false},
			{url: "/anything/else", crawlable: false, hasError: false},
		})
	})

	t.Run("/fish", func(t *testing.T) {
		robotsTxt, err := robotstxt.New("https://www.dumpsters.com", strings.NewReader(`
User-Agent: *
Disallow: /fish
`))
		assert.Nil(t, err)

		testRobot(t, "Bingbot", robotsTxt, []testUrl{
			{url: "/fish", crawlable: false, hasError: false},
			{url: "/fish.html", crawlable: false, hasError: false},
			{url: "/fish/salmon.html", crawlable: false, hasError: false},
			{url: "/fishheads", crawlable: false, hasError: false},
			{url: "/fishheads/yummy.html", crawlable: false, hasError: false},
			{url: "/fish.php?id=anything", crawlable: false, hasError: false},
			{url: "/Fish.asp", crawlable: true, hasError: false},
			{url: "/catfish", crawlable: true, hasError: false},
			{url: "/?id=fish", crawlable: true, hasError: false},
		})
	})

	// Equivalent to /fish. The trailing wildcard is ignored.
	t.Run("/fish*", func(t *testing.T) {
		robotsTxt, err := robotstxt.New("https://www.dumpsters.com", strings.NewReader(`
User-Agent: *
Disallow: /fish*
`))
		assert.Nil(t, err)

		testRobot(t, "Bingbot", robotsTxt, []testUrl{
			{url: "/fish", crawlable: false, hasError: false},
			{url: "/fish.html", crawlable: false, hasError: false},
			{url: "/fish/salmon.html", crawlable: false, hasError: false},
			{url: "/fishheads", crawlable: false, hasError: false},
			{url: "/fishheads/yummy.html", crawlable: false, hasError: false},
			{url: "/fish.php?id=anything", crawlable: false, hasError: false},
			{url: "/Fish.asp", crawlable: true, hasError: false},
			{url: "/catfish", crawlable: true, hasError: false},
			{url: "/?id=fish", crawlable: true, hasError: false},
		})
	})

	// Equivalent to /fish. The trailing wildcard is ignored.
	t.Run("/fish/", func(t *testing.T) {
		robotsTxt, err := robotstxt.New("https://www.dumpsters.com", strings.NewReader(`
User-Agent: *
Disallow: /fish/
`))
		assert.Nil(t, err)

		testRobot(t, "Bingbot", robotsTxt, []testUrl{
			{url: "/fish/", crawlable: false, hasError: false},
			{url: "/fish/?id=anything", crawlable: false, hasError: false},
			{url: "/fish/salmon.htm", crawlable: false, hasError: false},
			{url: "/fish", crawlable: true, hasError: false},
			{url: "/fish.html", crawlable: true, hasError: false},
			{url: "/Fish/Salmon.asp", crawlable: true, hasError: false},
		})
	})

	t.Run("/*.php", func(t *testing.T) {
		robotsTxt, err := robotstxt.New("https://www.dumpsters.com", strings.NewReader(`
User-Agent: *
Disallow: /*.php
`))
		assert.Nil(t, err)

		testRobot(t, "Bingbot", robotsTxt, []testUrl{
			{url: "/filename.php", crawlable: false, hasError: false},
			{url: "/folder/filename.php", crawlable: false, hasError: false},
			{url: "/folder/filename.php?parameters", crawlable: false, hasError: false},
			{url: "/folder/any.php.file.html", crawlable: false, hasError: false},
			{url: "/filename.php/", crawlable: false, hasError: false},
			{url: "/", crawlable: true, hasError: false},
			{url: "/windows.PHP", crawlable: true, hasError: false},
		})
	})

	t.Run("/*.php$", func(t *testing.T) {
		robotsTxt, err := robotstxt.New("https://www.dumpsters.com", strings.NewReader(`
User-Agent: *
Disallow: /*.php$
`))
		assert.Nil(t, err)

		testRobot(t, "Bingbot", robotsTxt, []testUrl{
			{url: "/filename.php", crawlable: false, hasError: false},
			{url: "/folder/filename.php", crawlable: false, hasError: false},
			{url: "/folder/filename.php?parameters", crawlable: true, hasError: false},
			{url: "/filename.php/", crawlable: true, hasError: false},
			{url: "/filename.php5", crawlable: true, hasError: false},
			{url: "/windows.PHP", crawlable: true, hasError: false},
		})
	})

	t.Run("/fish*.php", func(t *testing.T) {
		robotsTxt, err := robotstxt.New("https://www.dumpsters.com", strings.NewReader(`
User-Agent: *
Disallow: /fish*.php
`))
		assert.Nil(t, err)

		testRobot(t, "Bingbot", robotsTxt, []testUrl{
			{url: "/fish.php", crawlable: false, hasError: false},
			{url: "/fishheads/catfish.php?parameters", crawlable: false, hasError: false},
			{url: "/Fish.PHP", crawlable: true, hasError: false},
		})
	})
}

func TestNew(t *testing.T) {
	_, err := robotstxt.New("https://www.dumpsters.com", getExampleRobotsTxt())
	assert.Nil(t, err)
}

func TestNew_fails_if_no_scheme_is_provided(t *testing.T) {
	_, err := robotstxt.New("www.dumpsters.com", getExampleRobotsTxt())
	assert.NotNil(t, err)
}

func TestNew_fails_if_no_host_is_provided(t *testing.T) {
	_, err := robotstxt.New("https://", getExampleRobotsTxt())
	assert.NotNil(t, err)
}

func TestNewFromFile(t *testing.T) {
	filePath, err := filepath.Abs("./robots.txt")
	assert.Nil(t, err)

	_, err = robotstxt.NewFromFile("https://www.dumpsters.com", filePath)
	assert.Nil(t, err)
}

func TestRobotsTxt_CanCrawl(t *testing.T) {
	robotsTxt, err := robotstxt.New("https://www.dumpsters.com", getExampleRobotsTxt())
	assert.Nil(t, err)

	testRobot(t, "googlebot", robotsTxt, []testUrl{
		{url: "/cms/", crawlable: false, hasError: false},
		{url: "/cms", crawlable: true, hasError: false},
		{url: "/cms/pages", crawlable: false, hasError: false},
		{url: "/cms/pages?products=123", crawlable: false, hasError: false},

		{url: "/pricing/frontend", crawlable: false, hasError: false},
		{url: "/pricing/frontend-app", crawlable: false, hasError: false},
		{url: "/pricing/frontend/product", crawlable: false, hasError: false},

		{url: "/pricing/admin/product", crawlable: false, hasError: false},
		{url: "/pricing/admin", crawlable: true, hasError: false},

		{url: "/pricing?s=lightbox", crawlable: false, hasError: false},
		{url: "/pricing?s=lightbox&cart=full", crawlable: false, hasError: false},
		{url: "/pricing?cart=full&s=lightbox", crawlable: false, hasError: false},

		{url: "/se/en", crawlable: false, hasError: false},
		{url: "/se/en/", crawlable: true, hasError: false},
		{url: "/se", crawlable: true, hasError: false},
		{url: "/se/en/fr", crawlable: true, hasError: false},

		{url: "/retail/online/frontend/", crawlable: false, hasError: false},
		{url: "/store/retail/online/frontend/", crawlable: false, hasError: false},
		{url: "/retail/online/frontend/pages?page=2", crawlable: false, hasError: false},
		{url: "/online/frontend/", crawlable: true, hasError: false},
	})

	testRobot(t, "AdsBot-Google", robotsTxt, []testUrl{
		{url: "/cms/", crawlable: true, hasError: false},
		{url: "/cms/", crawlable: true, hasError: false},
		{url: "/cms", crawlable: true, hasError: false},
		{url: "/cms/pages", crawlable: true, hasError: false},
		{url: "/cms/pages?products=123", crawlable: true, hasError: false},

		{url: "/pricing/frontend", crawlable: true, hasError: false},
		{url: "/pricing/frontend-app", crawlable: true, hasError: false},
		{url: "/pricing/frontend/product", crawlable: true, hasError: false},

		{url: "/pricing/admin/product", crawlable: true, hasError: false},
		{url: "/pricing/admin", crawlable: true, hasError: false},

		{url: "/pricing?s=lightbox", crawlable: true, hasError: false},
		{url: "/pricing?s=lightbox&cart=full", crawlable: true, hasError: false},
		{url: "/pricing?cart=full&s=lightbox", crawlable: true, hasError: false},

		{url: "/se/en", crawlable: true, hasError: false},
		{url: "/se/en/", crawlable: true, hasError: false},
		{url: "/se", crawlable: true, hasError: false},
		{url: "/se/en/fr", crawlable: true, hasError: false},

		{url: "/retail/online/frontend/", crawlable: true, hasError: false},
		{url: "/store/retail/online/frontend/", crawlable: true, hasError: false},
		{url: "/retail/online/frontend/pages?page=2", crawlable: true, hasError: false},
		{url: "/online/frontend/", crawlable: true, hasError: false},
	})

	assert.Equal(t, 5*time.Second, robotsTxt.CrawlDelay("googlebot"))
	assert.Equal(t, 0*time.Second, robotsTxt.CrawlDelay("adsbot-google"))
	assert.Equal(t, []string{"https://www.dumpsters.com/sitemap.xml", "https://www.dumpsters.com/sitemap-launch-index.xml"}, robotsTxt.Sitemaps())
}

func TestRobotsTxt_CanCrawl_fails_if_robot_url_and_given_url_have_different_ports(t *testing.T) {
	robotsTxt, err := robotstxt.New("https://www.dumpsters.com:4343", getExampleRobotsTxt())
	assert.Nil(t, err)

	testRobot(t, "googlebot", robotsTxt, []testUrl{
		{url: "https://www.dumpsters.com:4000/cms", crawlable: true, hasError: true},
	})
}

func TestRobotsTxt_CanCrawl_fails_if_robot_url_and_given_url_have_different_schemes(t *testing.T) {
	robotsTxt, err := robotstxt.New("https://www.dumpsters.com:4000", getExampleRobotsTxt())
	assert.Nil(t, err)

	testRobot(t, "googlebot", robotsTxt, []testUrl{
		{url: "http://www.dumpsters.com:4000/cms", crawlable: true, hasError: true},
	})
}

func TestRobotsTxt_CrawlDelay(t *testing.T) {
	robotsTxt, err := robotstxt.New("https://www.dumpsters.com", getExampleRobotsTxt())
	assert.Nil(t, err)
	assert.Equal(t, 5*time.Second, robotsTxt.CrawlDelay("googlebot"))
}

func TestRobotsTxt_Sitemaps(t *testing.T) {
	robotsTxt, err := robotstxt.New("https://www.dumpsters.com", getExampleRobotsTxt())
	assert.Nil(t, err)
	assert.Equal(t, []string{"https://www.dumpsters.com/sitemap.xml", "https://www.dumpsters.com/sitemap-launch-index.xml"}, robotsTxt.Sitemaps())
}

func TestRobotsTxt_URL(t *testing.T) {
	robotsTxt, err := robotstxt.New("https://www.dumpsters.com", getExampleRobotsTxt())
	assert.Nil(t, err)
	assert.Equal(t, "https://www.dumpsters.com:443", robotsTxt.URL())
}

func TestRobotsTxt_URL_specifying_port(t *testing.T) {
	robotsTxt, err := robotstxt.New("https://www.dumpsters.com:4000", getExampleRobotsTxt())
	assert.Nil(t, err)
	assert.Equal(t, "https://www.dumpsters.com:4000", robotsTxt.URL())
}

/*
 *********************************************** START BENCHMARKS ***********************************************
 */

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = robotstxt.New("https://www.dumpsters.com", getExampleRobotsTxt())
	}
}

func BenchmarkRobotsTxt_CanCrawl(b *testing.B) {
	robotsTxt, _ := robotstxt.New("https://www.dumpsters.com", getExampleRobotsTxt())
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = robotsTxt.CanCrawl("Bingbot", "/cms/")
	}
}

func BenchmarkRobotsTxt_CanCrawl_multiple_times(b *testing.B) {
	robotsTxt, _ := robotstxt.New("https://www.dumpsters.com", getExampleRobotsTxt())
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = robotsTxt.CanCrawl("Bingbot", "/cms/")
		_, _ = robotsTxt.CanCrawl("Bingbot", "/cms")
		_, _ = robotsTxt.CanCrawl("Bingbot", "/cms/pages")
		_, _ = robotsTxt.CanCrawl("Bingbot", "/cms/pages?products=123")
		_, _ = robotsTxt.CanCrawl("Bingbot", "/pricing/frontend")
		_, _ = robotsTxt.CanCrawl("Bingbot", "/pricing/frontend-app")
		_, _ = robotsTxt.CanCrawl("Bingbot", "/pricing/frontend/product")
		_, _ = robotsTxt.CanCrawl("Bingbot", "/pricing/admin/product")
		_, _ = robotsTxt.CanCrawl("Bingbot", "/pricing/admin")
		_, _ = robotsTxt.CanCrawl("Bingbot", "/pricing?s=lightbox")
		_, _ = robotsTxt.CanCrawl("Bingbot", "/pricing?s-lightbox&cart=full")
		_, _ = robotsTxt.CanCrawl("Bingbot", "/pricing?cart=full&s=lightbox")
		_, _ = robotsTxt.CanCrawl("Bingbot", "/se/en")
		_, _ = robotsTxt.CanCrawl("Bingbot", "/se/en/")
		_, _ = robotsTxt.CanCrawl("Bingbot", "/se")
		_, _ = robotsTxt.CanCrawl("Bingbot", "/se/en/fr")
		_, _ = robotsTxt.CanCrawl("Bingbot", "/retail/online/frontend")
		_, _ = robotsTxt.CanCrawl("Bingbot", "/store/retail/online/frontend")
		_, _ = robotsTxt.CanCrawl("Bingbot", "/retail/online/frontend/pages?page=2")
		_, _ = robotsTxt.CanCrawl("Bingbot", "/online/frontend/")
	}
}

/*
 *********************************************** END BENCHMARKS ***********************************************
 */

/*
 *********************************************** START EXAMPLES ***********************************************
 */

func ExampleNew() {
	robotsTxt, _ := robotstxt.New("https://www.dumpsters.com", strings.NewReader(`
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
	fmt.Println(robotsTxt.Sitemaps())
	fmt.Println(robotsTxt.URL())
	fmt.Println(robotsTxt.CrawlDelay("googlebot"))
	// Output:
	// false
	// <nil>
	// [https://www.dumpsters.com/sitemap.xml https://www.dumpsters.com/sitemap-launch-index.xml]
	// https://www.dumpsters.com:443
	// 5s
}

func ExampleNewFromFile() {
	filePath, err := filepath.Abs("./robots.txt")
	fmt.Println(err)

	robotsTxt, err := robotstxt.NewFromFile("https://www.dumpsters.com", filePath)
	fmt.Println(err)

	canCrawl, err := robotsTxt.CanCrawl("googlebot", "/cms/pages")
	fmt.Println(canCrawl)
	fmt.Println(err)
	fmt.Println(robotsTxt.Sitemaps())
	fmt.Println(robotsTxt.URL())
	fmt.Println(robotsTxt.CrawlDelay("googlebot"))
	// Output:
	// <nil>
	// <nil>
	// false
	// <nil>
	// [https://www.dumpsters.com/sitemap.xml https://www.dumpsters.com/sitemap-launch-index.xml]
	// https://www.dumpsters.com:443
	// 5s
}

func ExampleNewFromURL() {
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

/*
 *********************************************** END EXAMPLES ***********************************************
 */

type testUrl struct {
	url       string
	crawlable bool
	hasError  bool
}

// I know it's bad to write code for tests, but testing each thing was painful and this will be consistent / less human error prone
func testRobot(t *testing.T, robotName string, robotsTxt *robotstxt.RobotsTxt, testUrls []testUrl) {
	for _, test := range testUrls {
		canCrawl, err := robotsTxt.CanCrawl(robotName, test.url)
		hasError := err != nil
		assert.Equal(t, test, testUrl{test.url, canCrawl, hasError})
	}
}

func getExampleRobotsTxt() io.Reader {
	return strings.NewReader(`
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

# Some odd cases are added below
user-agent test # Invalid line without a colon
: # Just a colon
`)
}
