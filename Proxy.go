package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/TwinProduction/go-color"
)

const CACHE_DURATION = 20000000000 // 20 seconds in nanoseconds

type website struct {
	headers     map[string]string
	body        []byte
	timeFetched time.Time
}

var mutex = &sync.RWMutex{}
var twoDots = regexp.MustCompile("\\.")
var blacklist = map[string]bool{} // Blacklist is a map with the URL's as the key and a boolean as the value
var cache = map[string]*website{} // Cache is a map of the URL's as the keys and then the values are the website structs storing the information about the websites
var webTimes = make(map[string]time.Duration, 0)
var cachetimes = make(map[string]time.Duration, 0)

func newSite(res *http.Response, siteResponse []byte) *website {
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

// Blacklisted
// Checks if URL is on blacklist including all subdomains on that URL
// Arguments: URL of Site
// Returns: Boolean
func blacklisted(site string) bool {
	dots := twoDots.FindAllStringIndex(site, -1)
	if len(dots) > 1 {
		sub := dots[len(dots)-2]
		site = site[sub[0]+1:]
	}
	port := strings.Index(site, ":")
	if port > -1 {
		site = site[:port]
	}
	_, blocked := blacklist[site]
	return blocked
}

func cached(site string) bool {
	mutex.RLock()
	website, x := cache[site]
	mutex.RUnlock()
	if x && website != nil && int64(time.Since(website.timeFetched)) < CACHE_DURATION {
		return true
	} else {
		mutex.Lock()
		delete(cache, site)
		mutex.Unlock()
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
		input = input[:len(input)-2]

		if strings.Contains(input, "/add") {
			site := input[5:]
			add2Blacklist(site)
		} else if strings.Contains(input, "/rmv") {
			site := input[5:]
			RmvFromBlacklist(site)
		} else if strings.Contains(input, "/view") {
			fmt.Println(color.Ize(color.Purple, "Blacklist:"))
			for i := range blacklist {
				println(color.Ize(color.Purple, fmt.Sprintf("| %s", i)))
			}
		} else if strings.Contains(input, "/cash") {
			savedTimeURL := make([]int64, 0)

			for i, y := range cachetimes {
				sv := int64(webTimes[i]/time.Millisecond) - int64(y/time.Millisecond)
				savedTimeURL = append(savedTimeURL, int64(sv))
			}
			if len(cachetimes) > 0 {
				average := int64(0)

				for _, y := range savedTimeURL {
					average = average + int64(y)
				}

				average = average / int64(len(savedTimeURL))
				fmt.Printf(color.Cyan, "Average time saved from caching: %dms", average)
				fmt.Println(color.Reset)
			} else if len(cache) == 0 {
				fmt.Println(color.Ize(color.Yellow, "Cache is Empty"))
			}
		}
	}
}

func HTTPHandler(writer http.ResponseWriter, request *http.Request) {
	client := &http.Client{}
	res, e := client.Do(request)

	if e != nil {
		log.Panic(e)
	}

	for i, y := range res.Header {
		for _, z := range y {
			writer.Header().Set(i, z)
		}
	}

	bodyBytes, e := ioutil.ReadAll(res.Body)
	if e != nil {
		log.Panic(e)
	}

	io.WriteString(writer, string(bodyBytes))
	mutex.Lock()
	cache[res.Request.URL.String()] = newSite(res, bodyBytes)
	mutex.Unlock()
	request.Body.Close()
	res.Body.Close()
}

func copyTCP(client *net.TCPConn, conn *net.TCPConn) {
	io.Copy(client, conn)
	client.Close()
	conn.Close()
}

func HTTPSHandler(writer http.ResponseWriter, request *http.Request) {
	time := time.Second * 10
	dest, e := net.DialTimeout("tcp", request.Host, time)

	if e != nil {
		http.Error(writer, e.Error(), http.StatusServiceUnavailable)
		log.Println(e)
		return
	}

	writer.WriteHeader(http.StatusOK)

	hijack, t := writer.(http.Hijacker)
	if !t {
		http.Error(writer, "Hijacking is not supported", http.StatusInternalServerError)
		log.Println(color.Ize(color.Red, "Hijacking is not supported"))
		return
	}
	client, _, e := hijack.Hijack()
	if e != nil {
		http.Error(writer, e.Error(), http.StatusServiceUnavailable)
	}

	destTCP, dOK := dest.(*net.TCPConn)
	clientTCP, cOK := client.(*net.TCPConn)

	if dOK && cOK {
		go copyTCP(destTCP, clientTCP)
		go copyTCP(clientTCP, destTCP)
	}

}

func mainHandler(writer http.ResponseWriter, request *http.Request) {
	request.RequestURI = ""
	url := request.URL.String()
	host := request.Host

	if !blacklisted(host) {
		cached := cached(url)
		if http.MethodConnect == request.Method {
			HTTPSHandler(writer, request)
		} else {
			if cached {
				mutex.RLock()
				site := cache[url]
				mutex.RUnlock()
				for i, y := range site.headers {
					writer.Header().Set(i, y)
				}
				io.WriteString(writer, string(site.body))
			} else {
				HTTPHandler(writer, request)
			}
		}
	} else {
		log.Println(color.Ize(color.Red, "This Site is BLOCKED!"))
		writer.WriteHeader(http.StatusForbidden)
	}
}

func cacheCleaner() {
	for 1 < 2 {
		mutex.RLock()
		for i, y := range cache {
			if int64(time.Since(y.timeFetched)) > CACHE_DURATION {
				mutex.Lock()
				delete(cache, i)
				mutex.Unlock()
			}
		}
		mutex.RUnlock()
	}
}

func main() {
	go userInput()
	go cacheCleaner()

	handler := http.HandlerFunc(mainHandler)
	http.ListenAndServe(":8080", handler)
}
