// title : always-english
// author :  ysoftman
// desc : 터미널창에서 영어 공부하기(영어 배너, 단어 찾기)

package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	ReadExampleRawData()
	ReadHistoryFile()
	ReadFortuneData()
	if len(os.Args) > 1 {
		word, pronounce, meanings := SearchEngWord(os.Args[1])
		fmt.Println(GetNextColorString(0, word))
		fmt.Println(GetNextColorString(1, pronounce))
		fmt.Println(GetNextColorString(2, meanings))
		os.Exit(0)
	}
	fixTerminfoEnv()
	StartGocui()
}

// fixTerminfoEnv termbox-go 는 TERMINFO 가 설정돼 있으면 그 디렉터리만 찾고 시스템 terminfo 로
// fallback 하지 않는다(ncurses 와 다름). Ghostty 처럼 TERMINFO 를 자체 경로로 지정하는 터미널 안에서
// tmux 를 쓰면 TERM=tmux-256color 를 못 찾아 panic 나므로, 거기에 현재 TERM 이 없으면 TERMINFO 를 뺀다.
func fixTerminfoEnv() {
	dir, term := os.Getenv("TERMINFO"), os.Getenv("TERM")
	if dir == "" || term == "" {
		return
	}
	// 일반 유닉스(t/) 와 darwin(74/) 디렉터리 구조 둘 다 확인
	for _, sub := range []string{term[:1], fmt.Sprintf("%x", term[0])} {
		if _, err := os.Stat(filepath.Join(dir, sub, term)); err == nil {
			return
		}
	}
	_ = os.Unsetenv("TERMINFO")
}
