// gocui 패키지를 이용한 환경 구성
// gocui 내부적으로 termbox-go 사용
package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	// "github.com/jroimartin/gocui" // 한글(utf8) 출력에 문제가 있음
	"github.com/ysoftman/gocui"
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("english_banner", 0, 0, maxX-1, maxY/4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "English Banner, " + BannerCmdText
		// fmt.Fprintln(v, GetNextColorString(0, "english banner"))
	}
	if v, err := g.SetView("word_history", 0, maxY/4+1, maxX-1, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Word History, " + HistoryCmdText
	}
	if v, err := g.SetView("search", 0, maxY/2+1, maxX-1, maxY/2+3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Highlight = true
		v.Frame = true
		v.Title = SearchCmdText + ", " + QuitCmdText
		// fmt.Fprintln(v, GetNextColorString(2, QuitCmdText))
		// fmt.Fprintln(v, GetNextColorString(1, SEARC_CMD_TEXT))
		v.SetCursor(0, 0)
		g.SetCurrentView("search")
	}

	if v, err := g.SetView("searchResult", 0, maxY/2+4, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Title = "Search Result"
		// fmt.Fprintln(v, GetNextColorString(2, QuitCmdText))
		// fmt.Fprintln(v, GetNextColorString(1, SEARC_CMD_TEXT))
		v.SetCursor(0, 0)
	}

	return nil
}

func searchAction(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		searchView, _ := g.View("search")
		word := strings.TrimSpace(searchView.Buffer())
		word, meanings, pronounce := SearchEngWord(word)
		searchView.Clear()
		searchView.SetCursor(0, 0)

		searchResultView, _ := g.View("searchResult")
		searchResultView.Clear()
		fmt.Fprint(searchResultView, GetNextColorString(0, word))
		pronounce = "  " + pronounce
		fmt.Fprint(searchResultView, GetNextColorString(1, pronounce))
		fmt.Fprint(searchResultView, GetNextColorString(4, meanings))
		return nil
	})
	return nil
}

func bannerUp(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		bannerView, _ := g.View("english_banner")
		bannerView.Clear()
		inner := GetPreBanner()
		for j := 0; j < len(inner); j++ {
			fmt.Fprintln(bannerView, GetNextColorString(j, inner[j]))
		}
		return nil
	})
	return nil
}

func bannerDown(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		bannerView, _ := g.View("english_banner")
		bannerView.Clear()
		inner := GetNextBanner()
		for j := 0; j < len(inner); j++ {
			fmt.Fprintln(bannerView, GetNextColorString(j, inner[j]))
		}
		return nil
	})
	return nil
}

func historyPre(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		historyView, _ := g.View("word_history")
		historyView.Clear()
		theWord := GetPreWord()
		if theWord != nil {
			fmt.Fprintln(historyView, GetNextColorString(0, theWord.word))
			fmt.Fprintln(historyView, GetNextColorString(1, theWord.pronounce))
			fmt.Fprintln(historyView, GetNextColorString(2, theWord.meanings))
		}
		return nil
	})
	return nil
}

func historyNext(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		historyView, _ := g.View("word_history")
		historyView.Clear()
		theWord := GetNextWord()
		if theWord != nil {
			fmt.Fprintln(historyView, GetNextColorString(0, theWord.word))
			fmt.Fprintln(historyView, GetNextColorString(1, theWord.pronounce))
			fmt.Fprintln(historyView, GetNextColorString(2, theWord.meanings))
		}
		return nil
	})
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	close(done)
	return gocui.ErrQuit
}

// StartGocui Gocui 구동
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
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, bannerUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, bannerDown); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, historyPre); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, historyNext); err != nil {
		log.Panicln(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return

			case <-time.After(BannerRefreshSec * time.Second):
				// ClearScreen()
				g.Update(func(g *gocui.Gui) error {
					bannerView, _ := g.View("english_banner")
					bannerView.Clear()
					inner := GetNextBanner()
					for j := 0; j < len(inner); j++ {
						fmt.Fprintln(bannerView, GetNextColorString(j, inner[j]))
					}
					return nil
				})
			}
		}
	}()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	wg.Wait()
}
