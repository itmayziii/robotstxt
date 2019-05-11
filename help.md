# Helpful Information
This file is more for my own reference than anything else.

**Coveralls Coverage**

`export COVERALLS_TOKEN=<token-value>`
`go test -cover -coverprofile=coverage.out && $GOPATH/bin/goveralls -coverprofile=coverage.out -repotoken=$COVERALLS_TOKEN`

**[Go Report Card](https://goreportcard.com/report/github.com/itmayziii/robotstxt)**
This just needs refreshed on every deploy.

**Creating a Release / Tag**
`git tag <tagname>, i.e. v1.0.0 && git push origin --tags`