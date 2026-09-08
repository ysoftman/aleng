// gocui 패키지를 이용한 환경 구성
// gocui 내부적으로 termbox-go 사용
package main

import (
	"fmt"
	"log"
	"math/rand"
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
		if _, err := g.SetCurrentView("search"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("searchResult", 0, 3, maxX-1, int(float32(maxY)*0.3)+3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Title = SearchResultCmdText
	}
	if v, err := g.SetView("example", 0, int(float32(maxY)*0.3)+3+1, maxX-1, int(float32(maxY)*0.6)+3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = ExampleCmdText
		// fmt.Fprintln(v, GetNextColorString(0, "example"))
	}
	if v, err := g.SetView("fortune", 0, int(float32(maxY)*0.6)+3+1, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = FortuneCmdText
	}

	return nil
}

func printExampleResult(g *gocui.Gui, inner []string) {
	exampleView, _ := g.View("example")
	exampleView.Clear()
	if len(inner) == 0 {
		_, _ = fmt.Fprintln(exampleView, GetNextColorString(0, NoResult))
		return
	}
	str := fmt.Sprintf("example: %v / %v", GetCurExampleIndex()+1, GetExamplesLen())
	_, _ = fmt.Fprintln(exampleView, GetNextColorString(3, str))
	for j := range inner {
		_, _ = fmt.Fprintln(exampleView, GetNextColorString(j, inner[j]))
	}
}

func upExample(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		inner := GetPreExample()
		printExampleResult(g, inner)
		return nil
	})
	return nil
}

func downExample(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		inner := GetNextExample()
		printExampleResult(g, inner)
		return nil
	})
	return nil
}

func findExample(g *gocui.Gui, keyword string) {
	g.Update(func(g *gocui.Gui) error {
		inner := FindExample(keyword)
		printExampleResult(g, inner)
		return nil
	})
}

func refreshFortuneTitle(g *gocui.Gui, cnt int) error {
	g.Update(func(g *gocui.Gui) error {
		fortuneView, _ := g.View("fortune")
		fortuneView.Title = fmt.Sprintf("%s / refresh randomly in %2d sec", FortuneCmdText, cnt)
		return nil
	})
	return nil
}

func printFortuneResult(g *gocui.Gui, inner []string) {
	fortuneView, _ := g.View("fortune")
	fortuneView.Clear()
	if len(inner) == 0 {
		_, _ = fmt.Fprintln(fortuneView, GetNextColorString(0, NoResult))
		return
	}
	str := fmt.Sprintf("fortune: %v / %v", GetCurFortuneIndex()+1, GetFortuneLen())
	_, _ = fmt.Fprintln(fortuneView, GetNextColorString(3, str))
	for j := range inner {
		_, _ = fmt.Fprintln(fortuneView, GetNextColorString(j, inner[j]))
	}
}
func upFortune(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		inner := GetPreFortune()
		printFortuneResult(g, inner)
		return nil
	})
	return nil
}
func downFortune(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		inner := GetNextFortune()
		printFortuneResult(g, inner)
		return nil
	})
	return nil
}
func printSearchWordResult(v *gocui.View, word, pronounce, meanings string, idx int) {
	if len(meanings) == 0 {
		_, _ = fmt.Fprintln(v, GetNextColorString(0, NoResult))
		return
	}
	if idx >= 0 {
		str := fmt.Sprintf("history: %v / %v "+SortHistoryCmdText, idx+1, MaxWordHistoryLimit)
		_, _ = fmt.Fprintln(v, GetNextColorString(3, str))
	}
	_, _ = fmt.Fprint(v, GetNextColorString(0, word))
	pronounce = "  " + pronounce + "\n"
	_, _ = fmt.Fprint(v, GetNextColorString(1, pronounce))
	_, _ = fmt.Fprint(v, GetNextColorString(2, meanings))
}

func printHistoryWord(v *gocui.View, ts, sf, word, pronounce, meanings string) {
	_, _ = fmt.Fprint(v, GetNextColorString(1, ts)+" ")
	_, _ = fmt.Fprint(v, GetNextColorString(2, sf)+" ")
	_, _ = fmt.Fprint(v, GetNextColorString(3, word))
	pronounce = "  " + pronounce + " "
	_, _ = fmt.Fprint(v, GetNextColorString(4, pronounce))
	mlist := strings.Split(meanings, "\n")
	mstr := strings.Join(mlist[:min(2, len(mlist))], " ") + " "
	_, _ = fmt.Fprint(v, GetNextColorString(2, mstr+"\n"))
}

func printHistoryWords(g *gocui.Gui, idx int, whd []WordHistoryData) error {
	g.Update(func(g *gocui.Gui) error {
		searchResultView, _ := g.View("searchResult")
		searchResultView.Clear()
		str := fmt.Sprintf("history: (%v~%v) / %v "+SortHistoryCmdText, idx+1, idx+len(whd), MaxWordHistoryLimit)
		_, _ = fmt.Fprintln(searchResultView, GetNextColorString(3, str))
		for i := range whd {
			printHistoryWord(searchResultView,
				whd[i].date.Format(time.RFC3339),
				strconv.Itoa(whd[i].searchFrequency),
				whd[i].wd.word,
				whd[i].wd.pronounce,
				whd[i].wd.meanings)
		}
		return nil
	})
	return nil
}

func specificHistory(g *gocui.Gui, startIdx int) error {
	idx, wl := GetWordHistoryInPage(startIdx)
	return printHistoryWords(g, idx, wl)
}

func previousHistory(g *gocui.Gui, v *gocui.View) error {
	idx, wl := GetPreWordHistoryInPage()
	return printHistoryWords(g, idx, wl)
}

func nextHistory(g *gocui.Gui, v *gocui.View) error {
	idx, wl := GetNextWordHistoryInPage()
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

func setSearchWordByHitory(g *gocui.Gui) error {
	searchView, _ := g.View("search")
	searchView.Clear()
	if curHistoryIndex < 0 {
		curHistoryIndex = 0
	}
	if curHistoryIndex >= len(wordHistory) {
		curHistoryIndex = len(wordHistory) - 1
	}
	word := wordHistory[curHistoryIndex].wd.word
	if _, err := searchView.Write([]byte(word)); err != nil {
		return err
	}
	_ = searchView.SetCursor(len(word), 0)
	return nil
}
func preSearchWord(g *gocui.Gui, v *gocui.View) error {
	curHistoryIndex--
	return setSearchWordByHitory(g)
}

func nextSearchWord(g *gocui.Gui, v *gocui.View) error {
	curHistoryIndex++
	return setSearchWordByHitory(g)
}

func searchWord(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		searchView, _ := g.View("search")
		word := strings.TrimSpace(searchView.Buffer())
		findExample(g, word)
		word, pronounce, meanings := SearchEngWord(word)
		searchView.Clear()
		_ = searchView.SetCursor(0, 0)

		searchResultView, _ := g.View("searchResult")
		searchResultView.Clear()
		printSearchWordResult(searchResultView, word, pronounce, meanings, -1)
		curHistoryIndex = 0
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
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, preSearchWord); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, nextSearchWord); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlK, gocui.ModNone, preSearchWord); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlJ, gocui.ModNone, nextSearchWord); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlI, gocui.ModNone, upFortune); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlO, gocui.ModNone, downFortune); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlU, gocui.ModNone, upExample); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlD, gocui.ModNone, downExample); err != nil {
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

	_ = downExample(g, nil)
	_ = downFortune(g, nil)
	_ = specificHistory(g, 0)
	var wg sync.WaitGroup
	wg.Add(1)
	remainRefreshSec = FortuneRefreshSec
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			case <-time.After(1 * time.Second):
				if remainRefreshSec == 0 {
					// ClearScreen()
					remainRefreshSec = FortuneRefreshSec
					curFortuneIndex = rand.Intn(len(fortuneData))
					_ = downFortune(g, nil)
				}
				_ = refreshFortuneTitle(g, remainRefreshSec)
				remainRefreshSec--
			}
		}
	}()

	// MainLoop 은 항상 non-nil 에러로 끝난다
	if err := g.MainLoop(); err != gocui.ErrQuit {
		log.Panicln(err)
	}

	wg.Wait()
}
