package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const CACHE_DURATION = 5000000000 // 5 seconds in nanoseconds

type website struct {
	headers     map[string]string
	body        []byte
	timeFetched time.Time
}

var twoDots = regexp.MustCompile("\\.")
var blacklist = map[string]bool{} // Blacklist is a map with the URL's as the key and a boolean as the value
var cache = map[string]*website{} // Cache is a map of the URL's as the keys and then the values are the website structs storing the information about the websites
var webTimes = make(map[string]time.Duration, 0)
var cachetimes = make(map[string]time.Duration, 0)

func add2Cache(res *http.Response, siteResponse []byte) *website {
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
		fmt.Printf("%s\n", site)
		// fmt.Printf("%s Blacklisted\n", site)
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
		fmt.Printf("%s has been removed from the blacklist\n", site)
	}
}

func blacklisted(site string) bool {
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

func cached(site string) bool {
	website, x := cache[site]
	if x && website != nil && int64(time.Since(website.timeFetched)) < CACHE_DURATION {
		return true
	} else {
		delete(cache, site)
		return false
	}
}

func userInput() {
	Scanner := bufio.NewReader(os.Stdin)
	fmt.Println("|--------------------------------|")
	fmt.Println("| Dans Web Proxy Console Bro ;-) |")
	fmt.Println("|--------------------------------|")

	for 1 < 2 {
		fmt.Print(">> ")
		input, _ := Scanner.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)

		fmt.Printf("%s\n", input)

		if strings.Contains(input, "/add") {
			site := input[4:]
			add2Blacklist(site)
		} else if strings.Contains(input, "/rmv") {
			site := input[4:]
			RmvFromBlacklist(site)
		} else if strings.Contains(input, "/view") {
			fmt.Println("Blacklist:")
			for i := range blacklist {
				fmt.Printf("| %s\n", i)
			}
		}
	}
}

func main() {
	userInput()

}
