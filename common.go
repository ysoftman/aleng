// 공통 함수들
package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

const BANNER_CMD_TEXT = "pre-banner (up) / next-banner (down)"
const HISTORY_CMD_TEXT = "pre-history (left) / next-history (right)"
const SEARCH_CMD_TEXT = "dic.daum.net search (enter)"
const QUIT_CMD_TEXT = "quit (ctrl + c)"
const BANNER_REFRESH_SEC = 10

var done = make(chan struct{})

var engDic []string

type WordData struct {
	word      string
	pronounce string
	meanings  string
}

var wordHistory []WordData
var curBannerIndex int
var curWordHistoryIndex int

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

func ReadDicFile() {
	curBannerIndex = -1
	eng, _ := ioutil.ReadFile("eng.dic")
	engDic = strings.Split(string(eng), "--")
}

func ReadHistroyFile() {
	wordHistory = nil
	curBannerIndex = -1
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

func GetNextBannerIndex() int {
	curBannerIndex++
	if curBannerIndex >= len(engDic) {
		curBannerIndex = 0
	}
	return curBannerIndex
}

func GetPreBannerIndex() int {
	curBannerIndex--
	if curBannerIndex < 0 {
		curBannerIndex = len(engDic) - 1
	}
	return curBannerIndex
}

func GetNextWordHistoryIndex() int {
	curWordHistoryIndex++
	if curWordHistoryIndex >= len(wordHistory) {
		curWordHistoryIndex = 0
	}
	return curWordHistoryIndex
}

func GetPreWordHistoryIndex() int {
	curWordHistoryIndex--
	if curWordHistoryIndex < 0 {
		curWordHistoryIndex = len(wordHistory) - 1
	}
	return curWordHistoryIndex
}

func GetPreBanner() []string {
	if len(engDic) > 0 {
		return strings.Split(strings.TrimPrefix(engDic[GetPreBannerIndex()], "\n"), "\n")
	}
	return nil
}

func GetNextBanner() []string {
	if len(engDic) > 0 {
		return strings.Split(strings.TrimPrefix(engDic[GetNextBannerIndex()], "\n"), "\n")
	}
	return nil
}

func GetPreWord() *WordData {
	if len(wordHistory) > 0 {
		return &wordHistory[GetPreWordHistoryIndex()]
	}
	return nil
}

func GetNextWord() *WordData {
	if len(wordHistory) > 0 {
		return &wordHistory[GetNextWordHistoryIndex()]
	}
	return nil
}

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
	meanings_one_line := ""
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		// meanings += s.Find("txt_search").Text()
		meanings += strconv.Itoa(cnt) + ". " + s.Text() + "\n"
		meanings_one_line += strconv.Itoa(cnt) + ". " + s.Text() + "   "
		cnt++
	})

	// save the word to history.txt
	if len(meanings) > 0 {
		addWord := WordData{strings.TrimSpace(word), strings.TrimSpace(pronounce), strings.TrimSpace(meanings_one_line)}

		// pop-back
		_, wordHistory = wordHistory[len(wordHistory)-1], wordHistory[:len(wordHistory)-1]
		// push-front
		wordHistory = append([]WordData{addWord}, wordHistory...)

		buffer := []byte(WordData2String())
		ioutil.WriteFile("history.txt", buffer, 0644)
	}

	pronounce += "\n"
	return word, meanings, pronounce
}
