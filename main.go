package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

func parseSite(site string, word string, ct int, finished chan bool, result map[string]int) {

	myword := []byte(word)
	client := http.Client{Timeout: time.Duration(ct) * time.Second}
	resp, err := client.Get(site)
	if err != nil {
		defer resp.Body.Close()
	}
	body, _ := ioutil.ReadAll(resp.Body)
	count := bytes.Count(body, myword)
	result[site] = count
	//fmt.Printf("Site %s is finished\n", site)
	finished <- true
}

func main() {

	var urls map[string]int
	var url string
	finished := make(chan bool)

	word := flag.String("word", "", "string to find")
	ct := flag.Int("ct", 5, "Connection Timeout")
	flag.Parse()

	urls = parseUrls(flag.Args())
	for url = range urls {
		//fmt.Printf("%s\n", url)
		go parseSite(url, *word, *ct, finished, urls)

	}
	<-finished
	var count int
	for url, count = range urls {
		fmt.Printf("%s - %d\n", url, count)
	}
	//mt.Printf("%s - %d - %d\n", *word, *ct, *rt)
	//fmt.Printf("%s\n", urls)
}
