package robotstxt_test

import (
	"fmt"
	"github.com/itmayziii/robotstxt"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

// Test cases derived from https://developers.google.com/search/reference/robots_txt#url-matching-based-on-path-values.
func Test_examples_mentioned_in_google_spec(t *testing.T) {
	// Matches the root and any lower level URL.
	t.Run("/", func(t *testing.T) {
		robotsTxt, err := robotstxt.New("", `
User-Agent: *
Disallow: /
`)
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
		robotsTxt, err := robotstxt.New("", `
User-Agent: *
Disallow: /*
`)
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
		robotsTxt, err := robotstxt.New("", `
User-Agent: *
Disallow: /fish
`)
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
		robotsTxt, err := robotstxt.New("", `
User-Agent: *
Disallow: /fish*
`)
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
		robotsTxt, err := robotstxt.New("", `
User-Agent: *
Disallow: /fish/
`)
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
		robotsTxt, err := robotstxt.New("", `
User-Agent: *
Disallow: /*.php
`)
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
		robotsTxt, err := robotstxt.New("", `
User-Agent: *
Disallow: /*.php$
`)
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
		robotsTxt, err := robotstxt.New("", `
User-Agent: *
Disallow: /fish*.php
`)
		assert.Nil(t, err)

		testRobot(t, "Bingbot", robotsTxt, []testUrl{
			{url: "/fish.php", crawlable: false, hasError: false},
			{url: "/fishheads/catfish.php?parameters", crawlable: false, hasError: false},
			{url: "/Fish.PHP", crawlable: true, hasError: false},
		})
	})
}

func Test_NewFromFile(t *testing.T) {
	filePath, err := filepath.Abs("./robots.txt")
	assert.Nil(t, err)
	robotsTxt, err := robotstxt.NewFromFile("https://www.dumpsters.com", filePath)
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

	assert.Equal(t, 5, robotsTxt.CrawlDelay("googlebot"))
	assert.Equal(t, 0, robotsTxt.CrawlDelay("adsbot-google"))
	assert.Equal(t, []string{"https://www.dumpsters.com/sitemap.xml", "https://www.dumpsters.com/sitemap-launch-index.xml"}, robotsTxt.Sitemaps())
}

/*
 *********************************************** START BENCHMARKS ***********************************************
 */

func BenchmarkNew_with_parsing(b *testing.B) {
	for i := 0; i < b.N; i++ {
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

func BenchmarkNew_only_methods(b *testing.B) {
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

func ExampleNewFromFile() {
	filePath, err := filepath.Abs("./robots.txt")
	robotsTxt, newFromFileErr := robotstxt.NewFromFile("https://www.dumpsters.com", filePath)
	fmt.Println(newFromFileErr)
	canCrawl, err := robotsTxt.CanCrawl("googlebot", "/cms/pages")
	fmt.Println(canCrawl)
	fmt.Println(err)
	// Output:
	// <nil>
	// false
	// <nil>
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
func testRobot(t *testing.T, robotName string, robotsTxt robotstxt.RobotsExclusionProtocol, testUrls []testUrl) {
	for _, test := range testUrls {
		canCrawl, err := robotsTxt.CanCrawl(robotName, test.url)
		hasError := err != nil
		assert.Equal(t, test, testUrl{test.url, canCrawl, hasError})
	}
}
