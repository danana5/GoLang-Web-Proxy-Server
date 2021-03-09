package main

import (
	"fmt"
	"net/http"
	"regexp"
	"time"
)

const cacheTime = 10

var twoDots = regexp.MustCompile("\\.")

var blacklist = map[string]bool{}
var cache = map[string]*website{}

var webTimes = make(map[string]time.Duration, 0)
var cachetimes = make(map[string]time.Duration, 0)

type website struct {
	headers     map[string]string
	body        []byte
	timeFetched time.Time
}

func addSite2Cache(res *http.Response, siteResponse []byte) *website {
	site := website{headers: make(map[string]string, 0), body: siteResponse}
	site.timeFetched = time.Now()

	for k, v := range res.Header {
		for _, vv := range v {
			site.headers[k] = vv
		}
	}
	return &site
}

func addSite2Blacklist(site string) {
	_, exists := blacklist[site]

	if !exists {
		blacklist[site] = true
		fmt.Printf("Added %s to Blacklist\n", site)
	} else {
		fmt.Println("Sites already blocked lad")
	}
}
