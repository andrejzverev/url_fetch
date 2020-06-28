package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

func parseUrls(s []string) map[string]int {
	// Init map
	urls := make(map[string]int)
	var url string
	// Make array as string and split it
	ss := strings.Split(strings.Join(s, " "), ",")
	for _, url = range ss {
		urls[strings.TrimSpace(url)] = 0
	}
	return urls
}

func parseSite(site string, word string, ct int, wg *sync.WaitGroup, result map[string]int) {

	defer wg.Done()
	myword := []byte(word)

	client := http.Client{
		Timeout: time.Duration(ct) * time.Second,
	}

	resp, err := client.Get(site)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	count := bytes.Count(body, myword)
	result[site] = count
}

func main() {

	var wg sync.WaitGroup
	var urls map[string]int
	var url string

	word := flag.String("word", "", "string to find")
	ct := flag.Int("ct", 10, "Connection Timeout")
	flag.Parse()

	urls = parseUrls(flag.Args())
	wg.Add(len(urls))
	for url = range urls {
		go parseSite(url, *word, *ct, &wg, urls)

	}
	wg.Wait()

	var count int
	for url, count = range urls {
		fmt.Printf("%s - contains word(%s) %d of times\n", url, *word, count)
	}
}
