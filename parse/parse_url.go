package parse

import (
	"bufio"
	"fmt"
	"nacs/utils"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func ProcessUrlToDiscover(DirectUrls []string) []map[string]interface{} {
	var ForDiscoverResults = make([]map[string]interface{}, 0)
	for _, DirectUrl := range DirectUrls {
		parsedUrl, _ := url.Parse(DirectUrl)
		schema := parsedUrl.Scheme
		host := parsedUrl.Host
		if strings.Contains(host, ":") {
			host = strings.Split(host, ":")[0]
		}
		port := parsedUrl.Port()
		path := parsedUrl.Path
		if schema == "" {
			schema = "http"
			host = DirectUrl
			DirectUrl = "http://" + DirectUrl
			port = "80"
			path = ""
		}
		if port == "" {
			if schema == "http" {
				port = "80"
			} else if schema == "https" {
				port = "443"
			}
		}
		if path == "/" {
			path = ""
			DirectUrl = DirectUrl[:len(DirectUrl)-1]
		}
		portInt, _ := strconv.Atoi(port)
		result := map[string]interface{}{
			"schema": schema,
			"host":   host,
			"port":   portInt,
			"path":   path,
			"url":    DirectUrl,
		}

		ForDiscoverResults = append(ForDiscoverResults, result)
	}
	return ForDiscoverResults
}

func ParseUrl(urls string, filename string) []map[string]interface{} {
	DirectUrls := ParseUrls(urls)
	if filename != "" {
		var fileurl []string
		fileurl, _ = Readurlfile(filename)
		DirectUrls = append(DirectUrls, fileurl...)
	}
	DirectUrls = utils.RemoveDuplicate(DirectUrls)
	return ProcessUrlToDiscover(DirectUrls)
}

// Readurlfile 按行读url
func Readurlfile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Open %s error, %v", filename, err)
		os.Exit(0)
	}
	defer file.Close()
	var content []string
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text != "" {
			content = append(content, text)
		}
	}
	return content, nil
}

func ParseUrls(urls string) (DirectUrls []string) {
	if urls != "" {
		if strings.Contains(urls, ",") {
			DirectUrls = strings.Split(urls, ",")
		} else {
			DirectUrls = append(DirectUrls, urls)
		}
	}
	return
}
