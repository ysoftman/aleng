// 공통 함수들
package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

// BannerCmdText : banner command text
const BannerCmdText = "English Banner, pre-banner (up) / next-banner (down)"

// SearchCmdText : search command text
const SearchCmdText = "search word in dic.daum.net and banners (enter)"

// SearchResultCmdText : search word / history command text
const SearchResultCmdText = "Search Result, pre-history (left) / next-history (right)"

// QuitCmdText : quit command text
const QuitCmdText = "quit (ctrl+c)"

// BannerRefreshSec : banner refresh interval(seconds)
const BannerRefreshSec = 10

var remainRefreshSec int

// MaxHistoryLimit : Max word history size
const MaxHistoryLimit = 10

var done = make(chan struct{})

var banners []string

// WordData : word data
type WordData struct {
	word      string
	pronounce string
	meanings  string
}

var wordHistory []WordData
var curBannerIndex int
var curWordHistoryIndex int

// ClearScreen : clear the screen
func ClearScreen() {
	cmdName := "clear"
	cmdArg1 := ""
	cmdArg2 := ""
	if runtime.GOOS == "windows" {
		cmdName = "cmd"
		cmdArg1 = "/c"
		cmdArg2 = "cls"
	}
	cmd := exec.Command(cmdName, cmdArg1, cmdArg2)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// GetNextColorString : get next color value(string)
func GetNextColorString(i int, str string) string {
	n := i % 6
	switch n {
	case 0:
		yellow := color.New(color.FgYellow).SprintFunc()
		return yellow(str)
	case 1:
		green := color.New(color.FgGreen).SprintFunc()
		return green(str)
	case 2:
		red := color.New(color.FgRed).SprintFunc()
		return red(str)
	case 3:
		blue := color.New(color.FgBlue).SprintFunc()
		return blue(str)
	case 4:
		magenta := color.New(color.FgMagenta).SprintFunc()
		return magenta(str)
	case 5:
		cyan := color.New(color.FgCyan).SprintFunc()
		return cyan(str)
	default:
		white := color.New(color.FgWhite).SprintFunc()
		return white(str)
	}
}

// ReadBannerFile : read the banner file
func ReadBannerFile() {
	rand.Seed(time.Now().UnixNano())
	eng, _ := ioutil.ReadFile("banner.txt")
	banners = strings.Split(string(eng), "--")
	curBannerIndex = rand.Intn(len(banners))
}

// ReadHistoryFile : read the history file
func ReadHistoryFile() {
	wordHistory = nil
	curWordHistoryIndex = -1
	history, err := ioutil.ReadFile("history.txt")
	if err != nil {
		return
	}
	spWord := strings.Split(string(history), "--\n")
	for i := 0; i < len(spWord); i++ {
		curWord := (strings.Split(spWord[i], "\n"))
		wd := WordData{}
		if len(curWord) > 0 {
			wd.word = curWord[0]
		}
		if len(curWord) > 1 {
			wd.pronounce = curWord[1]
		}
		if len(curWord) > 2 {
			wd.meanings = strings.Join(curWord[2:], "\n")
		}
		wordHistory = append(wordHistory, wd)
	}

	// limit max history size
	if len(wordHistory) > MaxHistoryLimit {
		wordHistory = wordHistory[:MaxHistoryLimit]
	}
}

// WordData2String : make one line string message from word data.
func WordData2String(wd []WordData) string {
	out := ""
	wc := len(wd)
	for i := 0; i < wc; i++ {
		out += wd[i].word + "\n" + wd[i].pronounce + "\n" + wd[i].meanings + "\n"
		if wc > 1 && i < wc-1 {
			out += "--\n"
		}
	}
	return out
}

// GetNextBannerIndex : get next banner index
func GetNextBannerIndex() int {
	curBannerIndex++
	if curBannerIndex >= len(banners) {
		curBannerIndex = 0
	}
	return curBannerIndex
}

// GetPreBannerIndex : get previous banner index
func GetPreBannerIndex() int {
	curBannerIndex--
	if curBannerIndex < 0 {
		curBannerIndex = len(banners) - 1
	}
	return curBannerIndex
}

// GetNextWordHistoryIndex : get next word history index
func GetNextWordHistoryIndex() int {
	curWordHistoryIndex++
	if curWordHistoryIndex >= len(wordHistory) {
		curWordHistoryIndex = 0
	}
	return curWordHistoryIndex
}

// GetPreWordHistoryIndex : get previous word history index
func GetPreWordHistoryIndex() int {
	curWordHistoryIndex--
	if curWordHistoryIndex < 0 {
		curWordHistoryIndex = len(wordHistory) - 1
	}
	return curWordHistoryIndex
}

// GetPreBanner : get previous banner
func GetPreBanner() []string {
	if len(banners) > 0 {
		return strings.Split(strings.TrimPrefix(banners[GetPreBannerIndex()], "\n"), "\n")
	}
	return nil
}

// FindBanner : find banner including keyword
func FindBanner(keyword string) []string {
	keyword = strings.ToLower(keyword)
	for i := range banners {
		// https://github.com/google/re2/wiki/Syntax
		// \b at ascii word boundary
		if matched, _ := regexp.MatchString("\\b"+keyword+"\\b", strings.ToLower(banners[i])); matched {
			// reset remainRefreshSec
			remainRefreshSec = BannerRefreshSec
			return strings.Split(strings.TrimPrefix(banners[i], "\n"), "\n")
		}
		// if strings.Contains(banners[i], keyword) {
		// 	return strings.Split(strings.TrimPrefix(banners[i], "\n"), "\n")
		// }
	}
	return nil
}

// GetNextBanner : get next banner
func GetNextBanner() []string {
	if len(banners) > 0 {
		return strings.Split(strings.TrimPrefix(banners[GetNextBannerIndex()], "\n"), "\n")
	}
	return nil
}

// GetPreWord : get previous word
func GetPreWord() (int, *WordData) {
	if len(wordHistory) > 0 {
		idx := GetPreWordHistoryIndex()
		return idx, &wordHistory[idx]
	}
	return 0, nil
}

// GetNextWord : get next word
func GetNextWord() (int, *WordData) {
	if len(wordHistory) > 0 {
		idx := GetNextWordHistoryIndex()
		return idx, &wordHistory[idx]
	}
	return 0, nil
}

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

	selector = "#mArticle > div.search_cont > div:nth-child(2) > div.search_box > div:nth-child(1) > div.search_word > strong > a"
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
		addWord := WordData{strings.TrimSpace(resultWord.word), strings.TrimSpace(resultWord.pronounce), strings.TrimSpace(resultWord.meanings)}

		// push-front
		wordHistory = append([]WordData{addWord}, wordHistory...)

		// save only MaxHistoryLimit
		buffer := []byte(WordData2String(wordHistory[:MaxHistoryLimit]))
		ioutil.WriteFile("history.txt", buffer, 0644)
	}
	return resultWord.word, resultWord.pronounce, resultWord.meanings
}
