package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {

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
