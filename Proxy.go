package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const CACHE_DURATION = 5

var twoDots = regexp.MustCompile("\\.")
var blacklist = map[string]bool{} // Blacklist is a map with the URL's as the key and a boolean as the value
var cache = map[string]*website{} // Cache is a map of the URL's as the keys and then the values are the website structs storing the information about the websites
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
	for k, i := range res.Header {
		for _, y := range i {
			site.headers[k] = y
		}
	}
	return &site
}

func add2Blacklist(site string) {
	_, blocked := blacklist[site]
	if !blocked {
		blacklist[site] = true
		fmt.Printf("Added %s to Blacklist\n", site)
	} else {
		fmt.Printf("Sites already blocked lad")
	}
}

func RmvFromBlacklist(site string) {
	_, blocked := blacklist[site]
	if !blocked {
		fmt.Println("Site is not blocked lad")
	} else {
		delete(blacklist, site)
		fmt.Printf("%s has been removed from the blacklist", site)
	}
}

func isBlocked(site string) bool {
	dots := twoDots.FindAllStringIndex(site, -1)
	if len(dots) > 1 {
		subIndex := dots[len(dots)-2]
		site = site[subIndex[0]+1:]
	}
	port := strings.Index(site, ":")
	if port > -1 {
		site = site[:port]
	}
	_, blocked := blacklist[site]
	return blocked
}
