package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html/charset"

	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/chardet"
)

func main() {
	url := "https://beam.apache.org/get-started/downloads/#releases"
	searchWord := "change"

	urls := getUrlsByWord(url, "Release notes")

	// それぞれのurlをvisitする
	for _, v := range urls {
		visitUrl(v, searchWord)
	}
}

func getDocumentByUrl(url string) *goquery.Document {
	// Getリクエストでレスポンス取得
	res, _ := http.Get(url)
	defer res.Body.Close()

	// Body内を読み取り
	buffer, _ := ioutil.ReadAll(res.Body)

	// 文字コードを判定
	detector := chardet.NewTextDetector()
	detectResult, _ := detector.DetectBest(buffer)
	// fmt.Println(detectResult.Charset)
	// => UTF-8

	// 文字コードの変換
	bufferReader := bytes.NewReader(buffer)
	reader, _ := charset.NewReaderLabel(detectResult.Charset, bufferReader)

	// HTMLをパース
	document, _ := goquery.NewDocumentFromReader(reader)

	return document
}

// textにwordを含むaタグのurlを抽出する
func getUrlsByWord(url string, word string) []string {
	urls := []string{}

	document := getDocumentByUrl(url)

	// urlを抜き出し
	document.Find("a").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		if title == word {
			val, exists := s.Attr("href")
			if exists {
				urls = append(urls, val)
			}
		}
	})

	return urls
}

func visitUrl(url string, searchWord string) {
	document := getDocumentByUrl(url)

	document.Find("#main").Find("ul").Find("li").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		isContain := strings.Contains(title, searchWord)
		if isContain {
			fmt.Println(isContain, title)
		}
	})
}
