package service

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/randyhg/test-log-scanner/model"
	"github.com/randyhg/test-log-scanner/util"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

func ScanGzFiles(directoryURL string) error {
	resp, err := http.Get(directoryURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// parse html from http GET request
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	// get all the .gz file urls from html nodes
	var fileList []string
	var extractLinks func(*html.Node)
	extractLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" && strings.HasSuffix(attr.Val, ".log.gz") {
					fileList = append(fileList, directoryURL+attr.Val)
				} // else if attr.Key == "href" && strings.HasSuffix(attr.Val, ".log") {
				//	fileList = append(fileList, directoryURL+attr.Val)
				//}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractLinks(c)
		}
	}
	extractLinks(doc)

	for _, fileName := range fileList {
		//is file has been scanned mechanism
		//if scanned := IsScanned(fileName); scanned {
		//	continue
		//}

		if err := GzippedLogFileReader(fileName); err != nil {
			fmt.Println("Error scanning gz file:", err)
			return err
		}

	}
	return nil
}

func GzippedLogFileReader(logURL string) error {
	resp, err := http.Get(logURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var stackFlag = false
	var stackTraces []string
	var currentTrace strings.Builder
	var traceHeadMessage string

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "error") && !strings.Contains(line, "goroutine") && !stackFlag {
			message := model.TimestampRegex.ReplaceAllString(line, "")
			match := model.TimestampRegex.FindStringSubmatch(line)
			if len(match) > 0 {
				timestamp := match[0]
				hash := model.Sha256(message)
				logError := model.LogErrors{
					Message:  message,
					FailedAt: timestamp,
					Hash:     &hash,
				}

				// store data to DB or another stores
				if err = util.Master().Create(&logError).Error; err != nil {
					return err
				}
			}
		} else if strings.Contains(line, "goroutine") && !stackFlag {
			traceHeadMessage = line
			stackFlag = true
			currentTrace.WriteString(line + "\n")
		} else if stackFlag && !model.TimestampRegex.MatchString(line) {
			currentTrace.WriteString(line + "\n")
		} else if stackFlag && model.TimestampRegex.MatchString(line) {
			stackFlag = false
			stackTraces = append(stackTraces, currentTrace.String())
			index := len(stackTraces)
			stackTrace := stackTraces[index-1]
			currentTrace.Reset()

			message := model.TimestampRegex.ReplaceAllString(traceHeadMessage, "")
			match := model.TimestampRegex.FindStringSubmatch(line)
			if len(match) > 0 {
				timestamp := match[0]
				hash := model.Sha256(message)
				logError := model.LogErrors{
					Message:    message,
					StackTrace: &stackTrace,
					FailedAt:   timestamp,
					Hash:       &hash,
				}

				// store data to DB or another stores
				if err = util.Master().Create(&logError).Error; err != nil {
					return err
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Println(logURL, "successfully scanned")
	return nil
}
