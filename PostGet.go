package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
)

var urls = []string{
"http://webcode.me",
"https://example.com",
"http://httpbin.org",
"https://www.perl.org",
"https://www.php.net",
"https://www.python.org",
"https://code.visualstudio.com",
"https://clojure.org",
}

func post(u string) {
	data := url.Values{
		"name":			{"123"},
		"occupation":   {"456"},
	}
	resp, err := http.PostForm(u, data)
	if err != nil {
		log.Fatalf("Can't post\n")
	}
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	fmt.Println(res["form"])

}

func getReqAsync(urls []string) {
	var wg sync.WaitGroup
	for _, url := range urls{
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			getReq(url)
		}(url)
	}
	wg.Wait()
}

func getReq(url string) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Can't get info\n")
	}
	defer response.Body.Close()
	pageText, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Can't read from get requests body\n")
	}
	f, err := os.Create(getTitle(url, string(pageText)) + ".html")
	if err != nil {
		log.Fatalf("Can't create file\n")
	}
	defer f.Close()
	_, err = f.Write(pageText)
	if err != nil {
		log.Fatalf("Can't write to file\n")
	}
}

func getTitle(url string, pageText string) string {
	re := regexp.MustCompile("<title>(.)*</title>")
	parts := re.FindStringSubmatch(pageText)
	if len(parts) != 0 {
		return strings.Trim(parts[0], "<title></title>")
	}
	return url
}
