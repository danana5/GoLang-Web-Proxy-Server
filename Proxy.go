package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/TwinProduction/go-color"
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
		fmt.Printf(color.Ize(color.Green, "Blacklisted\n"))
	} else {
		fmt.Println(color.Ize(color.Yellow, "This site is already on the blacklist"))
	}
}

func RmvFromBlacklist(site string) {
	_, blocked := blacklist[site]
	if !blocked {
		fmt.Println(color.Ize(color.Yellow, "Site is not blocked lad"))
	} else {
		delete(blacklist, site)
		fmt.Println(color.Ize(color.Green, "Removed from blacklist"))
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
	fmt.Println(color.Ize(color.Cyan, "|-------------------------------|"))
	fmt.Println(color.Ize(color.Cyan, "| Dans Web Proxy Console Bro ;) |"))
	fmt.Println(color.Ize(color.Cyan, "|-------------------------------|"))

	for 1 < 2 {
		fmt.Print(color.Ize(color.Blue, ">> "))
		input, _ := Scanner.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)

		if strings.Contains(input, "/add") {
			site := input[5:]
			add2Blacklist(site)
		} else if strings.Contains(input, "/rmv") {
			site := input[5:]
			RmvFromBlacklist(site)
		} else if strings.Contains(input, "/view") {
			fmt.Println(color.Ize(color.Bold, "Blacklist:"))
			for i := range blacklist {
				println(color.Ize(color.Purple, fmt.Sprintf("| %s", i)))
			}
		}
	}
}

func HTTPHandler(site http.ResponseWriter, req *http.Request) {
	client := &http.Client{}
	res, e := client.Do(req)

	if e != nil {
		log.Panic(e)
	}

	for i, y := range res.Header {
		for _, z := range y {
			site.Header().Set(i, z)
		}
	}

	bodyBytes, e := ioutil.ReadAll(res.Body)
	if e != nil {
		log.Panic(e)
	}

	io.WriteString(site, string(bodyBytes))
	cache[res.Request.URL.String()] = add2Cache(res, bodyBytes)
	req.Body.Close()
	res.Body.Close()
}

func main() {
	userInput()
}
