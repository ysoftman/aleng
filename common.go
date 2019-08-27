// 공통 함수들
package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

// BannerCmdText : banner command text
const BannerCmdText = "English Banner, pre-banner (up) / next-banner (down)"

// HistoryCmdText : history command text
const HistoryCmdText = "Word History, pre-history (left) / next-history (right)"

// SearchCmdText : search command text
const SearchCmdText = "dic.daum.net search (enter)"

// QuitCmdText : quit command text
const QuitCmdText = "quit (ctrl + c)"

// BannerRefreshSec : banner refresh interval(seconds)
const BannerRefreshSec = 10

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
		addWord := WordData{curWord[0], curWord[1], curWord[2]}
		wordHistory = append(wordHistory, addWord)
	}

	// limit max history size
	if len(wordHistory) > 10 {
		wordHistory = wordHistory[:10]
	}
}

// WordData2String : make one line string message from word data.
func WordData2String() string {
	out := ""
	wc := len(wordHistory)
	for i := 0; i < wc; i++ {
		out += wordHistory[i].word + "\n" + wordHistory[i].pronounce + "\n" + wordHistory[i].meanings + "\n"
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

// GetNextBanner : get next banner
func GetNextBanner() []string {
	if len(banners) > 0 {
		return strings.Split(strings.TrimPrefix(banners[GetNextBannerIndex()], "\n"), "\n")
	}
	return nil
}

// GetPreWord : get previous word
func GetPreWord() *WordData {
	if len(wordHistory) > 0 {
		return &wordHistory[GetPreWordHistoryIndex()]
	}
	return nil
}

// GetNextWord : get next word
func GetNextWord() *WordData {
	if len(wordHistory) > 0 {
		return &wordHistory[GetNextWordHistoryIndex()]
	}
	return nil
}

// SearchEngWord : search english word through dic.daum.net
func SearchEngWord(word string) (string, string, string) {
	// using http.Get() in NewDocument
	query := "https://dic.daum.net/search.do?q=" + word
	doc, err := goquery.NewDocument(query)
	if err != nil {
		log.Fatal(err)
	}

	pronounce := ""
	meanings := ""

	childIndex := "2"
	// 관련 단어 영역의 존재에 따라 영역 인덱스가 달라진다.
	if doc.Find("#relatedQuery").Length() > 0 {
		childIndex = "3"
	}

	// # : 23(hex)
	// copy selector string using chrome dev tool
	selector := "#mArticle div.search_cont div:nth-child(" + childIndex + ") div:nth-child(2) div div.search_cleanword strong a span.txt_emph1"
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		word = s.Text()
	})

	selector = "#mArticle div.search_cont div:nth-child(" + childIndex + ") div:nth-child(2) div  div.wrap_listen span:nth-child(1) span.txt_pronounce"
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		pronounce += s.Text() + "  "
	})

	// selector = "#mArticle div.search_cont div:nth-child(3) div:nth-child(2) div ul li:nth-child(1) span.txt_search"
	selector = "#mArticle div.search_cont div:nth-child(" + childIndex + ") div:nth-child(2) div ul li"
	cnt := 1
	meaningsOneLine := ""
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		meanings += strconv.Itoa(cnt) + ". " + s.Find("span.txt_search").Text() + "\n"
		meaningsOneLine += strconv.Itoa(cnt) + ". " + s.Find("span.txt_search").Text() + "   "
		cnt++
	})

	// save the word to history.txt
	if len(meanings) > 0 {
		addWord := WordData{strings.TrimSpace(word), strings.TrimSpace(pronounce), strings.TrimSpace(meaningsOneLine)}

		// push-front
		wordHistory = append([]WordData{addWord}, wordHistory...)

		buffer := []byte(WordData2String())
		ioutil.WriteFile("history.txt", buffer, 0644)
	}
	return word, pronounce, meanings
}
