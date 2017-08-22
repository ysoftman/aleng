// title : always-english
// author :  ysoftman
// desc : 터미널창에서 영어 문장을 계속 보여줌~ㅋ
// dependency
// go get github.com/fatih/color

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("english_banner", 0, 0, maxX-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, getNextColorString(0, "english banner"))
	}
	if v, err := g.SetView("input", 0, maxY/2+1, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, getNextColorString(1, "search : "))
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
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
	fmt.Fprint(v, s)
}

func main() {

	// g, err := gocui.NewGui(gocui.OutputNormal)
	// if err != nil {
	// 	log.Panicln(err)
	// }
	// defer g.Close()
	// g.SetManagerFunc(layout)
	// if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
	// 	log.Panicln(err)
	// }

	// layout(g)
	// bannerView, _ := g.View("english_banner")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		eng, _ := ioutil.ReadFile("eng.dic")
		dic := strings.Split(string(eng), "--")

		defer wg.Done()
		for {
			for i := 1; i < len(dic); i++ {
				clearScreen()
				// bannerView.Clear()
				inner := strings.Split(string(dic[i]), "\n")
				for j := 1; j < len(inner); j++ {
					fmt.Println(getNextColorString(j-1, inner[j]))
					// setViewTextAndCursor(bannerView, getNextColorString(j-1, inner[j]), 0, j)
				}
				time.Sleep(3 * time.Second)
			}
		}
	}()
	wg.Wait()

	// if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
	// 	log.Panicln(err)
	// }
}
