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
	"time"

	"github.com/fatih/color"
)

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

func main() {
	eng, _ := ioutil.ReadFile("eng.dic")
	dic := strings.Split(string(eng), "--")
	for {
		for i := 1; i < len(dic); i++ {
			clearScreen()
			inner := strings.Split(string(dic[i]), "\n")
			for j := 1; j < len(inner); j++ {
				fmt.Println(getNextColorString(j-1, inner[j]))
			}
			time.Sleep(5 * time.Second)
		}
	}
}
