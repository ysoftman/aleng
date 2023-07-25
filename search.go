// search
package main

import (
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// SearchCmdText : search command text
const SearchCmdText = "search word (enter) / pre (up or ctrl+k) / next (down or ctrl+j)"

// SearchResultCmdText : search word command text
const SearchResultCmdText = "Search Result, pre (left or ctrl+h) / next (right or ctrl+l)"

// SearchEngWord : search english word through dic.daum.net
func SearchEngWord(word string) (string, string, string) {
	// using http.Get() in NewDocument
	urlenc := &url.URL{Path: word}
	query := "https://dic.daum.net/search.do?q=" + urlenc.String()
	doc, err := goquery.NewDocument(query)
	if err != nil {
		log.Fatal(err)
	}

	var resultWord WordData

	childIndex := "2"
	// 관련 단어 영역의 존재에 따라 영역 인덱스가 달라진다.
	if doc.Find("#relatedQuery").Length() > 0 {
		childIndex = "3"
	} else if doc.Find("#mArticle > div.search_cont > div.card_relate").Length() > 0 {
		childIndex = "3"
	}

	// # : 23(hex)
	// copy selector string using chrome dev tool
	selector := "#mArticle > div.search_cont > div:nth-child(" + childIndex + ") > div:nth-child(2) > div > div.search_cleanword > strong > a"
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		resultWord.word = s.Text()
	})

	selector = "#mArticle > div.search_cont > div:nth-child(" + childIndex + ") > div:nth-child(2) > div > div.wrap_listen > span:nth-child(1) > span.txt_pronounce"
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		resultWord.pronounce += s.Text() + "  "
	})

	selector = "#mArticle > div.search_cont > div:nth-child(" + childIndex + ") > div:nth-child(2) > div > ul > li"
	cnt := 1
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		resultWord.meanings += strconv.Itoa(cnt) + ". " + s.Find("span.txt_search").Text() + "\n"
		cnt++
	})

	// save the word to history.txt
	if len(resultWord.meanings) > 0 {
		addWordHistory := WordHistoryData{WordData{strings.TrimSpace(resultWord.word), strings.TrimSpace(resultWord.pronounce), strings.TrimSpace(resultWord.meanings)}, time.Now(), 1}

		// if found, increase search-frequency
		if idx := findWordInHistory(addWordHistory.wd.word); idx >= 0 {
			addWordHistory.searchFrequency = wordHistory[idx].searchFrequency + 1
			// delete found word from word history
			wordHistory = append(wordHistory[:idx], wordHistory[idx+1:]...)
		}
		// push-front
		wordHistory = append([]WordHistoryData{addWordHistory}, wordHistory...)

		// save only MaxHistoryLimit
		if len(wordHistory) > MaxWordHistoryLimit {
			wordHistory = wordHistory[:MaxWordHistoryLimit]
		}
		buffer := []byte(WordHistoryData2String(wordHistory))
		if err := ioutil.WriteFile(historyFile, buffer, 0644); err != nil {
			log.Fatal("error, failed to write history file.")
		}
	}
	return resultWord.word, resultWord.pronounce, resultWord.meanings
}

func findWordInHistory(word string) int {
	for idx, v := range wordHistory {
		if v.wd.word == word {
			return idx
		}
	}
	return -1
}
