// common
package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/fatih/color"
)

// QuitCmdText : quit command text
const QuitCmdText = "quit (ctrl+c)"

// NoResult : no result text
const NoResult = "-- NO RESULT --"

var remainRefreshSec int

var done = make(chan struct{})

// WordData : word data
type WordData struct {
	word      string
	pronounce string
	meanings  string
}

// WordHistoryData : word history data
type WordHistoryData struct {
	wd              WordData
	date            time.Time
	searchFrequency int
}

// ClearScreen : clear the screen
func ClearScreen() {
	cmdName := "clear"
	cmdArg1 := ""
	cmdArg2 := ""
	if runtime.GOOS == "windows" {
		cmdName = "cmd"
		cmdArg1 = "/c"
		cmdArg2 = "cls"
	}
	cmd := exec.Command(cmdName, cmdArg1, cmdArg2)
	// link command stdout to os stdout (clear target)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// GetNextColorString : get next color value(string)
func GetNextColorString(i int, str string) string {
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
