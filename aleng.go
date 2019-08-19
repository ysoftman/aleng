// title : always-english
// author :  ysoftman
// desc : 터미널창에서 영어 공부하기(영어 배너, 단어 찾기)

package main

func main() {
	// for debuging.
	// fmt.Println(SearchEngWord("love"))
	// os.Exit(0)

	ReadBannerFile()
	ReadHistoryFile()
	StartGocui()
	// StartTermBoxGo()
}
