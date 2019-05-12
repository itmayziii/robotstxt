package robotstxt

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestNewFromUrl(t *testing.T) {
	old := httpGet
	defer func() { httpGet = old }()
	httpGet = func(url string) (*http.Response, error) {
		return &http.Response{Body: ioutil.NopCloser(strings.NewReader(fakeHTML()))}, nil
	}

	ch := make(chan ProtocolResult)
	go NewFromURL("https://www.dumpsters.com", ch)
	protocol := <-ch

	assert.Nil(t, protocol.Error)
}

func TestNew_utf8_validation(t *testing.T) {
	old := validateUTF8
	defer func() { validateUTF8 = old }()
	validateUTF8 = func(s string) bool {
		return false
	}
	_, err := New("https://www.dumpsters.com", `
# This file is mocked to have invalid UTF8
`)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid encoding detected on line 1, all characters must be UTF-8 encoded", err.Error())
}

func fakeHTML() string {
	return `
<html><head></head><body><pre style="word-wrap: break-word; white-space: pre-wrap;"># Robots.txt for dumpsters.com
# 06/04/2018

User-agent: *
Disallow: /bdso/
Disallow: /bdso3/
Disallow: /bdso4/
Disallow: /cgi-bin/
Disallow: /dev/
Disallow: /help/
Disallow: /phones/
Disallow: /page-content/
Disallow: /bd2/
Disallow: /old-site/
Disallow: /xfers/
Disallow: /pdf/
Disallow: /demo/
Disallow: /wiki/
Disallow: /calc/
Disallow: /lp/
Disallow: /*?PageSpeed=noscript
Disallow: /*?Modpagespeed=noscript
Disallow: /*?pagespeed=noscript
Disallow: /*?mod*
Disallow: /*?Page*
Disallow: /*?page*
Disallow: /*?Mod*
Disallow: /simplified-rentals*
Disallow: /flexible-solutions*
Disallow: /superior-service*
Disallow: /cleveland-dumpster-rental* 

User-agent: AdsBot-Google-Mobile
Allow: /

User-agent: AdsBot-Google
Allow: /

Sitemap: https://www.dumpsters.com/sitemap.xml</pre></body></html>
`
}
