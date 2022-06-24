// 공통 함수들
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

// BannerCmdText : banner command text
const BannerCmdText = "English Banner, pre (ctrl+u) / next (ctrl+d)"

// SearchCmdText : search command text
const SearchCmdText = "search word (enter) / pre (up) / next (down)"

// SearchResultCmdText : search word / history command text
const SearchResultCmdText = "Search Result, pre (left or ctrl+h) / next (right or ctrl+l)"

// SortCmdText : sort by time or frequency
const SortCmdText = "Sort by time(ctrl+t) / frequency(ctrl+f)"

// QuitCmdText : quit command text
const QuitCmdText = "quit (ctrl+c)"

// NoResult : no result text
const NoResult = "-- NO RESULT --"

// BannerRefreshSec : banner refresh interval(seconds)
const BannerRefreshSec = 10

var remainRefreshSec int

var historyFile string = "aleng_history.txt"

var usr *user.User

// MaxWordHistoryLimit : Max word history size
const MaxWordHistoryLimit = 5000
const historyPerPage = 10

// SortType : sort type
type SortType int

// enum sortType
const (
	SortByTime SortType = 1 + iota
	SortBySearchFrequency
)

var done = make(chan struct{})

var banners []string

// WordData : word data
type WordData struct {
	word      string
	pronounce string
	meanings  string
}

// WordHistoryData : word history data
type WordHistoryData struct {
	wd              WordData
	date            time.Time
	searchFrequency int
}

var wordHistory []WordHistoryData
var curHistoryIndex int
var curHistoryPageIndex int
var curBannerIndex int

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

// ReadBannerRawData : read the banner file
func ReadBannerRawData() {
	rand.Seed(time.Now().UnixNano())
	banners = strings.Split(string(bannerRawData), "--")
	curBannerIndex = rand.Intn(len(banners))
}

// ReadHistoryFile : read the history file
func ReadHistoryFile() {
	wordHistory = nil
	curHistoryPageIndex = 0
	var err error
	usr, err = user.Current()
	if err != nil {
		log.Fatal("Can't find user current() err=>", err.Error())
	}
	historyFile = usr.HomeDir + "/" + historyFile

	if _, err = os.Stat(historyFile); os.IsNotExist(err) {
		_, err = os.Create(historyFile)
		if err != nil {
			log.Fatalf("Can't create file:(%v)  err=> %v", historyFile, err.Error())
		}
		fmt.Printf("%v file created.\n", historyFile)
	}
	history, err := ioutil.ReadFile(historyFile)
	if err != nil {
		log.Fatalf("Can't read file:(%v)  err=> %v", historyFile, err.Error())
		return
	}

	spWord := strings.Split(string(history), "--\n")
	for i := 0; i < len(spWord); i++ {
		curWord := (strings.Split(spWord[i], "\n"))
		whd := WordHistoryData{}
		if len(curWord) < 5 {
			continue
		}
		whd.date, _ = time.Parse(time.RFC3339, curWord[0])
		whd.searchFrequency, _ = strconv.Atoi(curWord[1])
		whd.wd.word = curWord[2]
		whd.wd.pronounce = curWord[3]
		whd.wd.meanings = strings.Join(curWord[4:], "\n")
		wordHistory = append(wordHistory, whd)
	}

	// limit max history size
	if len(wordHistory) > MaxWordHistoryLimit {
		wordHistory = wordHistory[:MaxWordHistoryLimit]
	}
}

// SortWordHistoryData : sort word history data by date or frequency...
func SortWordHistoryData(whd []WordHistoryData, sortType SortType) {
	if sortType == SortByTime {
		sort.Slice(whd, func(a, b int) bool {
			return whd[a].date.Unix() > whd[b].date.Unix()
		})
		return
	}
	if sortType == SortBySearchFrequency {
		sort.Slice(whd, func(a, b int) bool {
			return whd[a].searchFrequency > whd[b].searchFrequency
		})
		return
	}
}

// WordHistoryData2String : make one line string message from word data.
func WordHistoryData2String(whd []WordHistoryData) string {
	out := ""
	wc := len(whd)
	for i := 0; i < wc; i++ {
		out += strings.TrimSpace(whd[i].date.Format(time.RFC3339)+"\n"+strconv.Itoa(whd[i].searchFrequency)+"\n"+whd[i].wd.word+"\n"+whd[i].wd.pronounce+"\n"+whd[i].wd.meanings) + "\n"
		if wc > 1 && i < wc-1 {
			out += "--\n"
		}
	}
	return out
}

// GetBannerLen : get number of banners
func GetBannerLen() int {
	return len(banners)
}

// GetCurBannerIndex : get current banner index
func GetCurBannerIndex() int {
	return curBannerIndex
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

// GetNextHistoryPageIndex : get next word history index
func GetNextHistoryPageIndex() int {
	curHistoryPageIndex = curHistoryPageIndex + historyPerPage
	if curHistoryPageIndex >= len(wordHistory) {
		curHistoryPageIndex = 0
	}
	return curHistoryPageIndex
}

// GetPreHistoryPagaeIndex : get previous word history index
func GetPreHistoryPagaeIndex() int {
	curHistoryPageIndex = curHistoryPageIndex - historyPerPage
	if curHistoryPageIndex < 0 {
		curHistoryPageIndex = ((len(wordHistory) - 1) / historyPerPage) * historyPerPage
	}
	return curHistoryPageIndex
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
	var foundBanners []string
	for i := range banners {
		// if strings.Contains(banners[i], keyword) {
		// 	return strings.Split(strings.TrimPrefix(banners[i], "\n"), "\n")
		// }
		// https://github.com/google/re2/wiki/Syntax
		// \b at ascii word boundary
		if matched, _ := regexp.MatchString("\\b"+keyword+"\\b", strings.ToLower(banners[i])); matched {
			// reset remainRefreshSec
			remainRefreshSec = BannerRefreshSec
			fb := strings.Split(strings.TrimPrefix(banners[i], "\n"), "\n")
			for j := range fb {
				foundBanners = append(foundBanners, fb[j])
				curBannerIndex = i
			}
		}
	}
	return foundBanners
}

// GetNextBanner : get next banner
func GetNextBanner() []string {
	if len(banners) > 0 {
		return strings.Split(strings.TrimPrefix(banners[GetNextBannerIndex()], "\n"), "\n")
	}
	return nil
}

func getWordHistoryInPage(startIdx int) (int, []WordHistoryData) {
	whd := []WordHistoryData{}
	if startIdx < 0 {
		return 0, whd
	}
	whLen := len(wordHistory)
	if len(wordHistory) > 0 {
		for i := startIdx; i < startIdx+historyPerPage; i++ {
			if i >= whLen {
				break
			}
			whd = append(whd, wordHistory[i])
		}
	}
	return startIdx, whd
}

// GetWordHistoryInPage : get words in a specific page
func GetWordHistoryInPage(startIdx int) (int, []WordHistoryData) {
	return getWordHistoryInPage(startIdx)
}

// GetPreWordHistoryInPage : get previous words in a page
func GetPreWordHistoryInPage() (int, []WordHistoryData) {
	return getWordHistoryInPage(GetPreHistoryPagaeIndex())
}

// GetNextWordsInPage : get next words in a page
func GetNextWordsInPage() (int, []WordHistoryData) {
	return getWordHistoryInPage(GetNextHistoryPageIndex())
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
