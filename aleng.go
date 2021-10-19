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
		fmt.Println(GetNextColorString(0, word))
		fmt.Println(GetNextColorString(1, pronounce))
		fmt.Println(GetNextColorString(2, meanings))
		os.Exit(0)
	}

	StartGocui()
}
