// termbox-go 패키지를 이용한 환경 구성
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

func printString(x, y int, str string, fgcolor termbox.Attribute) {
	// for index, runeValue := range str {
	// 	termbox.SetCell(x+index, y, runeValue, fgcolor, termbox.ColorDefault)
	// }
	for _, runeValue := range str {
		termbox.SetCell(x, y, runeValue, fgcolor, termbox.ColorDefault)
		w := runewidth.RuneWidth(runeValue)
		if w == 0 || (w == 2 && runewidth.IsAmbiguousWidth(runeValue)) {
			w = 1
		}
		x += w
	}
	termbox.Flush()
}

func StartTermBoxGo() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return

			case <-time.After(BANNER_REFRESH_SEC * time.Second):
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				printString(0, 5, SEARCH_CMD_TEXT+", "+QUIT_CMD_TEXT, termbox.ColorWhite)
				inner := GetNextBanner()
				for j := 0; j < len(inner); j++ {
					fgcolor := termbox.ColorYellow
					if j%2 == 0 {
						fgcolor = termbox.ColorGreen
					} else if j%3 == 0 {
						fgcolor = termbox.ColorRed
					}
					printString(0, j-1, inner[j], fgcolor)
				}
			}
		}
	}()

	termbox.SetCursor(0, 6)

	inputString := ""
	inputXpos := 0
mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlQ, termbox.KeyCtrlC:
				close(done)
				break mainloop
			case termbox.KeyEnter:
				word, meaning, pronounce := SearchEngWord(inputString)
				fmt.Println(word)
				fmt.Println(pronounce)
				fmt.Println(meaning)
				inputString = ""
				inputXpos = 0
			default:
				if ev.Ch != 0 {
					inputString += string(ev.Ch)
					printString(inputXpos, 7, string(ev.Ch), termbox.ColorWhite)
					inputXpos++
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
	wg.Wait()
}
