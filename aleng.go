// title : always-english
// author :  ysoftman
// desc : 터미널창에서 영어 공부하기(영어 배너, 단어 찾기)

package main

import (
	"fmt"
	"os"
)

func main() {
	ReadBannerRawData()
	ReadHistoryFile()
	if len(os.Args) > 1 {
		word, pronounce, meanings := SearchEngWord(os.Args[1])
		fmt.Println(word)
		fmt.Println(pronounce)
		fmt.Println(meanings)
		os.Exit(0)
	}

	StartGocui()
	// StartTermBoxGo()	// 아직 사용하는데 문제 있음
}
