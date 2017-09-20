// title : always-english
// author :  ysoftman
// desc : 터미널창에서 영어 공부하기(영어 배너, 단어 찾기)

package main

import (
	"fmt"
	"sync"
	"time"
)

func StartBanner() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ClearScreen()
		for {
			select {
			case <-done:
				return

			case <-time.After(BANNER_REFRESH_SEC * time.Second):
				ClearScreen()
				inner := GetNextBannerContent()
				for j := 1; j < len(inner); j++ {
					fmt.Println(GetNextColorString(j-1, inner[j]))
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
		meanings, pronounce := SearchEngWord(word)
		fmt.Println(GetNextColorString(0, pronounce))
		fmt.Println(GetNextColorString(0, meanings))
	}
}

func main() {
	// fmt.Println(SearchEngWord("love"))
	// os.Exit(0)
	ReadDicFile()
	// StartBanner()
	// StartSearchEngWord()
	StartGocui()
	// StartTermBoxGo()
}
