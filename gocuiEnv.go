// gocui 패키지를 이용한 환경 구성
// gocui 내부적으로 termbox-go 사용
package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	// "github.com/jroimartin/gocui" // 한글(utf8) 출력에 문제가 있음
	"github.com/ysoftman/gocui"
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("search", 0, 0, maxX-1, 2); err != nil {
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

	if v, err := g.SetView("searchResult", 0, 3, maxX-1, (maxY/2)+3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Title = SearchResultCmdText
		// fmt.Fprintln(v, GetNextColorString(2, QuitCmdText))
		// fmt.Fprintln(v, GetNextColorString(1, SEARC_CMD_TEXT))
		v.SetCursor(0, 0)
	}
	if v, err := g.SetView("english_banner", 0, ((maxY/2)+3)+1, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = BannerCmdText
		// fmt.Fprintln(v, GetNextColorString(0, "english banner"))
	}
	return nil
}

func refreshBannerTitle(g *gocui.Gui, v *gocui.View, cnt int) error {
	g.Update(func(g *gocui.Gui) error {
		bannerView, _ := g.View("english_banner")
		bannerView.Title = fmt.Sprintf("%s / refresh in %2d sec", BannerCmdText, cnt)
		return nil
	})
	return nil
}

func printBannerResult(g *gocui.Gui, inner []string) {
	bannerView, _ := g.View("english_banner")
	bannerView.Clear()
	if len(inner) == 0 {
		fmt.Fprintln(bannerView, GetNextColorString(0, NoResult))
		return
	}
	str := fmt.Sprintf("banner: %v / %v", GetCurBannerIndex()+1, GetBannerLen())
	fmt.Fprintln(bannerView, GetNextColorString(3, str))
	for j := 0; j < len(inner); j++ {
		fmt.Fprintln(bannerView, GetNextColorString(j, inner[j]))
	}
}

func upBanner(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		inner := GetPreBanner()
		printBannerResult(g, inner)
		return nil
	})
	return nil
}

func downBanner(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		inner := GetNextBanner()
		printBannerResult(g, inner)
		return nil
	})
	return nil
}

func findBanner(g *gocui.Gui, v *gocui.View, keyword string) error {
	g.Update(func(g *gocui.Gui) error {
		inner := FindBanner(keyword)
		printBannerResult(g, inner)
		return nil
	})
	return nil
}

func printSearchWordResult(v *gocui.View, word, pronounce, meanings string, idx int) {
	if len(meanings) == 0 {
		fmt.Fprintln(v, GetNextColorString(0, NoResult))
		return
	}
	if idx >= 0 {
		str := fmt.Sprintf("history: %v / %v "+SortCmdText, idx+1, MaxWordHistoryLimit)
		fmt.Fprintln(v, GetNextColorString(3, str))
	}
	fmt.Fprint(v, GetNextColorString(0, word))
	pronounce = "  " + pronounce + "\n"
	fmt.Fprint(v, GetNextColorString(1, pronounce))
	fmt.Fprint(v, GetNextColorString(2, meanings))
}

func printHistoryWord(v *gocui.View, ts, sf, word, pronounce, meanings string, idx int) {
	fmt.Fprint(v, GetNextColorString(1, ts)+" ")
	fmt.Fprint(v, GetNextColorString(2, sf)+" ")
	fmt.Fprint(v, GetNextColorString(3, word))
	pronounce = "  " + pronounce + " "
	fmt.Fprint(v, GetNextColorString(4, pronounce))
	mlist := strings.Split(meanings, "\n")
	mstr := ""
	for i, v := range mlist {
		if i >= 2 {
			break
		}
		mstr += v + " "
	}
	fmt.Fprint(v, GetNextColorString(2, mstr+"\n"))
}

func printHistoryWords(g *gocui.Gui, idx int, whd []WordHistoryData) error {
	g.Update(func(g *gocui.Gui) error {
		searchResultView, _ := g.View("searchResult")
		searchResultView.Clear()
		str := fmt.Sprintf("history: (%v~%v) / %v "+SortCmdText, idx+1, idx+len(whd), MaxWordHistoryLimit)
		fmt.Fprintln(searchResultView, GetNextColorString(3, str))
		for i := 0; i < len(whd); i++ {
			printHistoryWord(searchResultView,
				whd[i].date.Format(time.RFC3339),
				strconv.Itoa(whd[i].searchFrequency),
				whd[i].wd.word,
				whd[i].wd.pronounce,
				whd[i].wd.meanings,
				idx+i)
		}
		return nil
	})
	return nil
}

func specificHistory(g *gocui.Gui, v *gocui.View, startIdx int) error {
	idx, wl := GetWordHistoryInPage(startIdx)
	return printHistoryWords(g, idx, wl)
}

func previousHistory(g *gocui.Gui, v *gocui.View) error {
	idx, wl := GetPreWordHistoryInPage()
	return printHistoryWords(g, idx, wl)
}

func nextHistory(g *gocui.Gui, v *gocui.View) error {
	idx, wl := GetNextWordsInPage()
	return printHistoryWords(g, idx, wl)
}

func sortbywordHistoryByTime(g *gocui.Gui, v *gocui.View) error {
	SortWordHistoryData(wordHistory, SortByTime)
	return printHistoryWords(g, 0, wordHistory[:historyPerPage])
}

func sortbywordHistoryBySearchFrequency(g *gocui.Gui, v *gocui.View) error {
	SortWordHistoryData(wordHistory, SortBySearchFrequency)
	return printHistoryWords(g, 0, wordHistory[:historyPerPage])
}

func searchWord(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		searchView, _ := g.View("search")
		word := strings.TrimSpace(searchView.Buffer())
		findBanner(g, v, word)
		word, pronounce, meanings := SearchEngWord(word)
		searchView.Clear()
		searchView.SetCursor(0, 0)

		searchResultView, _ := g.View("searchResult")
		searchResultView.Clear()
		printSearchWordResult(searchResultView, word, pronounce, meanings, -1)
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
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, searchWord); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, upBanner); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, downBanner); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlK, gocui.ModNone, upBanner); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlJ, gocui.ModNone, downBanner); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, previousHistory); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, nextHistory); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlH, gocui.ModNone, previousHistory); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlL, gocui.ModNone, nextHistory); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlT, gocui.ModNone, sortbywordHistoryByTime); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlF, gocui.ModNone, sortbywordHistoryBySearchFrequency); err != nil {
		log.Panicln(err)
	}

	downBanner(g, nil)
	specificHistory(g, nil, 0)
	var wg sync.WaitGroup
	wg.Add(1)
	remainRefreshSec = BannerRefreshSec
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			case <-time.After(1 * time.Second):
				if remainRefreshSec == 0 {
					// ClearScreen()
					downBanner(g, nil)
					remainRefreshSec = BannerRefreshSec
				}
				refreshBannerTitle(g, nil, remainRefreshSec)
				remainRefreshSec--
			}
		}
	}()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	wg.Wait()
}
