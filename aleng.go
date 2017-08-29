// title : always-english
// author :  ysoftman
// desc : 터미널창에서 영어 공부하기(영어 배너, 단어 찾기)

package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"time"
)

func StartBanner() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		eng, _ := ioutil.ReadFile("eng.dic")
		dic := strings.Split(string(eng), "--")
		index := 1
		ClearScreen()
		inner := strings.Split(string(dic[index]), "\n")
		for j := 1; j < len(inner); j++ {
			fmt.Println(GetNextColorString(j-1, inner[j]))
		}

		for {
			select {
			case <-done:
				return

			case <-time.After(3 * time.Second):
				ClearScreen()
				inner := strings.Split(string(dic[index]), "\n")
				for j := 1; j < len(inner); j++ {
					fmt.Println(GetNextColorString(j-1, inner[j]))
				}
				index++
				if index >= len(dic) {
					index = 0
				}
			}
		}
	}()

	wg.Wait()
}

func StartSearchEngWord() {
	for {
		var word string
		fmt.Scanf("%s", &word)
		meanings, _ := SearchEngWord(word)
		fmt.Println(GetNextColorString(0, meanings))
	}
}

func main() {
	// fmt.Println(SearchEngWord("love"))
	// os.Exit(0)
	// StartBanner()
	// StartSearchEngWord()
	StartGocui()
	// StartTermBoxGo()
}
