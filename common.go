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
const BannerCmdText = "pre-banner (up) / next-banner (down)"

// HistoryCmdText : history command text
const HistoryCmdText = "pre-history (left) / next-history (right)"

// SearchCmdText : search command text
const SearchCmdText = "dic.daum.net search (enter)"

// QuitCmdText : quit command text
const QuitCmdText = "quit (ctrl + c)"

// BannerRefreshSec : banner refresh interval(seconds)
const BannerRefreshSec = 10

var done = make(chan struct{})

var engDic []string

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

// ReadDicFile : read the dictionary file(.dic)
func ReadDicFile() {
	rand.Seed(time.Now().UnixNano())
	eng, _ := ioutil.ReadFile("eng.dic")
	engDic = strings.Split(string(eng), "--")
	curBannerIndex = rand.Intn(len(engDic))
}

// ReadHistoryFile : read the history file(.txt)
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
	if curBannerIndex >= len(engDic) {
		curBannerIndex = 0
	}
	return curBannerIndex
}

// GetPreBannerIndex : get previous banner index
func GetPreBannerIndex() int {
	curBannerIndex--
	if curBannerIndex < 0 {
		curBannerIndex = len(engDic) - 1
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
	if len(engDic) > 0 {
		return strings.Split(strings.TrimPrefix(engDic[GetPreBannerIndex()], "\n"), "\n")
	}
	return nil
}

// GetNextBanner : get next banner
func GetNextBanner() []string {
	if len(engDic) > 0 {
		return strings.Split(strings.TrimPrefix(engDic[GetNextBannerIndex()], "\n"), "\n")
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
	query := "http://dic.daum.net/search.do?q=" + word
	doc, err := goquery.NewDocument(query)
	if err != nil {
		log.Fatal(err)
	}

	meanings := ""
	pronounce := ""

	selector := "#mArticle div.search_cont div.card_word.\\23 word.\\23 eng div.search_box.\\23 box div div.search_cleanword strong a span"
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		word = s.Text()
	})

	selector = "#mArticle div.search_cont div.card_word.\\23 word.\\23 eng div.search_box.\\23 box div  div.wrap_listen span:nth-child(1) .txt_pronounce"
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		pronounce += s.Text() + "  "
	})

	// copy selector string using chrome dev tool
	// #mArticle > div.search_cont > div.card_word.\23 word.\23 eng > div.search_box.\23 box > div > ul > li:nth-child(1) > span.txt_search
	selector = "#mArticle div.search_cont div.card_word.\\23 word.\\23 eng .search_box.\\23 box div ul.list_search span.txt_search"

	cnt := 1
	meaningsOneLine := ""
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		// meanings += s.Find("txt_search").Text()
		meanings += strconv.Itoa(cnt) + ". " + s.Text() + "\n"
		meaningsOneLine += strconv.Itoa(cnt) + ". " + s.Text() + "   "
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

	pronounce += "\n"
	return word, meanings, pronounce
}
