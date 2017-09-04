// 공통 함수들
package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

const SEARCH_WORD_TEXT = "dic.daum.net search (enter)"
const QUIT_WORD_TEXT = "quit (ctrl + c)"
const BANNER_REFRESH_SEC = 5

var done = make(chan struct{})

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

func SearchEngWord(word string) (string, string) {
	// using http.Get() in NewDocument
	query := "http://dic.daum.net/search.do?q=" + word
	doc, err := goquery.NewDocument(query)
	if err != nil {
		log.Fatal(err)
	}

	meanings := ""
	pronounce := ""

	selector := "#mArticle div.search_cont div.card_word.\\23 word.\\23 eng div.search_box.\\23 box div  div.wrap_listen span:nth-child(1) .txt_pronounce"
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		pronounce += s.Text() + "  "
	})
	pronounce += "\n"

	// copy selector string using chrome dev tool
	// #mArticle > div.search_cont > div.card_word.\23 word.\23 eng > div.search_box.\23 box > div > ul > li:nth-child(1) > span.txt_search
	selector = "#mArticle div.search_cont div.card_word.\\23 word.\\23 eng .search_box.\\23 box div ul.list_search span.txt_search"

	cnt := 1
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		// meanings += s.Find("txt_search").Text()
		meanings += strconv.Itoa(cnt) + ". " + s.Text() + "\n"
		cnt++
	})

	return meanings, pronounce
}
