// title : always-english
// author :  ysoftman
// desc : 터미널창에서 영어 문장을 계속 보여줌~ㅋ
// dependency
// go get -u github.com/fatih/color
// go get -u github.com/jroimartin/gocui

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
)

const SEARCH_WORD_TEXT = "dic.daum.net search (enter)"
const QUIT_WORD_TEXT = "quit (ctre + c)"

var done = make(chan struct{})

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("english_banner", 0, 0, maxX-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, getNextColorString(0, "english banner"))
	}
	if v, err := g.SetView("search", 0, maxY/2+1, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		fmt.Fprintln(v, getNextColorString(2, QUIT_WORD_TEXT))
		fmt.Fprintln(v, getNextColorString(1, SEARCH_WORD_TEXT))
		v.SetCursor(0, 3)
		g.SetCurrentView("search")
	}

	return nil
}

func getNextColorString(i int, str string) string {
	n := i % 6
	switch n {
	case 0:
		yellow := color.New(color.FgYellow).SprintFunc()
		return yellow(str)
	case 1:
		green := color.New(color.FgGreen).SprintFunc()
		return green(str)
	case 2:
		red := color.New(color.FgRed).SprintFunc()
		return red(str)
	case 3:
		blue := color.New(color.FgBlue).SprintFunc()
		return blue(str)
	case 4:
		magenta := color.New(color.FgMagenta).SprintFunc()
		return magenta(str)
	case 5:
		cyan := color.New(color.FgCyan).SprintFunc()
		return cyan(str)
	default:
		white := color.New(color.FgWhite).SprintFunc()
		return white(str)
	}

}

func clearScreen() {
	cmdName := "clear"
	cmdArg1 := ""
	cmdArg2 := ""
	if runtime.GOOS == "windows" {
		cmdName = "cmd"
		cmdArg1 = "/c"
		cmdArg2 = "cls"
	}
	cmd := exec.Command(cmdName, cmdArg1, cmdArg2)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func setViewTextAndCursor(v *gocui.View, s string, x, y int) {
	v.SetCursor(x, y)
	fmt.Fprintln(v, s)
}

func searchWord(word string) (string, string) {
	// using http.Get() in NewDocument
	query := "http://dic.daum.net/search.do?q=" + word
	doc, err := goquery.NewDocument(query)
	if err != nil {
		log.Fatal(err)
	}

	meanings := ""
	sentence := ""
	// copy selector string using chrome dev tool
	// #mArticle > div.search_cont > div.card_word.\23 word.\23 eng > div.search_box.\23 box > div > ul > li:nth-child(1) > span.txt_search
	selector := "#mArticle div.search_cont div.card_word .search_box div ul.list_search span.txt_search"

	cnt := 0
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		// meanings += s.Find("txt_search").Text()
		if cnt < 3 {
			meanings += s.Text() + "\n"
		}
		cnt++
	})

	return meanings, sentence
}

func searchAction(g *gocui.Gui, v *gocui.View) error {

	g.Update(func(g *gocui.Gui) error {
		searchView, _ := g.View("search")

		// scanner := bufio.NewScanner(os.Stdin)
		// var buf []byte
		// // set scanner buffer
		// scanner.Buffer(buf, 100)
		// // by words
		// scanner.Split(bufio.ScanWords)
		// scanner.Scan()
		// word := string(buf)

		word := strings.TrimSpace(searchView.Buffer())
		meanings, _ := searchWord(word)
		// fmt.Println(meanings)
		// fmt.Println(sentence)

		searchView.Clear()
		setViewTextAndCursor(searchView, getNextColorString(0, QUIT_WORD_TEXT), 0, 0)
		setViewTextAndCursor(searchView, getNextColorString(0, SEARCH_WORD_TEXT), 0, 1)
		setViewTextAndCursor(searchView, getNextColorString(0, "---"), 0, 2)
		setViewTextAndCursor(searchView, getNextColorString(0, word), 0, 3)
		setViewTextAndCursor(searchView, getNextColorString(0, meanings), 0, 4)
		return nil
	})
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	close(done)
	return gocui.ErrQuit
}

func main() {
	// fmt.Println(searchWord("love"))
	// os.Exit(0)

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
				// clearScreen()
				g.Update(func(g *gocui.Gui) error {
					bannerView, _ := g.View("english_banner")
					bannerView.Clear()

					inner := strings.Split(string(dic[index]), "\n")
					for j := 1; j < len(inner); j++ {
						// fmt.Println(getNextColorString(j-1, inner[j]))
						setViewTextAndCursor(bannerView, getNextColorString(j-1, inner[j]), 0, 0)
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
