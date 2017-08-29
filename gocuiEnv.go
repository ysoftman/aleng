// gocui 패키지를 이용한 환경 구성
// gocui 내부적으로 termbox-go 사용
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"

	// "github.com/jroimartin/gocui" // 한글(utf8) 출력에 문제가 있음
	"github.com/ysoftman/gocui"
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("english_banner", 0, 0, maxX-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "English Banner"
		// fmt.Fprintln(v, GetNextColorString(0, "english banner"))
	}
	if v, err := g.SetView("search", 0, maxY/2+1, maxX-1, maxY/2+3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Highlight = true
		v.Frame = true
		v.Title = SEARCH_WORD_TEXT + ", " + QUIT_WORD_TEXT
		// fmt.Fprintln(v, GetNextColorString(2, QUIT_WORD_TEXT))
		// fmt.Fprintln(v, GetNextColorString(1, SEARCH_WORD_TEXT))
		v.SetCursor(0, 0)
		g.SetCurrentView("search")
	}

	if v, err := g.SetView("searchResult", 0, maxY/2+5, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Title = "Search Result"
		// fmt.Fprintln(v, GetNextColorString(2, QUIT_WORD_TEXT))
		// fmt.Fprintln(v, GetNextColorString(1, SEARCH_WORD_TEXT))
		v.SetCursor(0, 0)
	}

	return nil
}

func setViewTextAndCursor(v *gocui.View, s string, x, y int) {
	v.SetCursor(x, y)
	fmt.Fprintln(v, s)
}

func searchAction(g *gocui.Gui, v *gocui.View) error {

	g.Update(func(g *gocui.Gui) error {
		searchView, _ := g.View("search")
		word := strings.TrimSpace(searchView.Buffer())
		meanings, _ := SearchEngWord(word)
		searchView.Clear()
		searchView.SetCursor(0, 0)

		searchResultView, _ := g.View("searchResult")
		searchResultView.Clear()
		setViewTextAndCursor(searchResultView, GetNextColorString(0, word), 0, 1)
		setViewTextAndCursor(searchResultView, GetNextColorString(0, meanings), 0, 2)
		return nil
	})
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	close(done)
	return gocui.ErrQuit
}

func StartGocui() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.SetManagerFunc(layout)
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, searchAction); err != nil {
		log.Panicln(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		eng, _ := ioutil.ReadFile("eng.dic")
		dic := strings.Split(string(eng), "--")
		index := 0
		for {
			select {
			case <-done:
				return

			case <-time.After(3 * time.Second):
				// ClearScreen()
				g.Update(func(g *gocui.Gui) error {
					bannerView, _ := g.View("english_banner")
					bannerView.Clear()
					inner := strings.Split(string(dic[index]), "\n")
					for j := 1; j < len(inner); j++ {
						// fmt.Println(GetNextColorString(j-1, inner[j]))
						setViewTextAndCursor(bannerView, GetNextColorString(j-1, inner[j]), 0, 0)
					}
					return nil
				})
				index++
				if index >= len(dic) {
					index = 0
				}
			}
		}
	}()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	wg.Wait()
}
