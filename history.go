// history
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"sort"
	"strconv"
	"strings"
	"time"
)

// SortHistoryCmdText : sort history by time or frequency
const SortHistoryCmdText = "Sort by time(ctrl+t) / frequency(ctrl+f)"

// MaxWordHistoryLimit : Max word history size
const MaxWordHistoryLimit = 5000
const historyPerPage = 10

var historyFile string = "aleng_history.txt"
var wordHistory []WordHistoryData
var curHistoryIndex int
var curHistoryPageIndex int

// SortType : sort type
type SortType int

// enum sortType
const (
	SortByTime SortType = 1 + iota
	SortBySearchFrequency
)

// ReadHistoryFile : read the history file
func ReadHistoryFile() {
	wordHistory = nil
	curHistoryPageIndex = 0
	usr, err := user.Current()
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

// GetNextWordHistoryInPage : get next words in a page
func GetNextWordHistoryInPage() (int, []WordHistoryData) {
	return getWordHistoryInPage(GetNextHistoryPageIndex())
}
