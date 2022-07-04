// fortune
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os/exec"
	"strings"
)

// FortuneCmdText : fortune command text
const FortuneCmdText = "Fortune, pre (ctrl+i) / next (ctrl+o)"

// FortuneRefreshSec : refresh fortune interval(seconds)
const FortuneRefreshSec = 10

var fortuneData []string
var curFortuneIndex int

// ReadFortuneData : read fortune(command) data
func ReadFortuneData() {
	execOutput := ""
	if out, err := exec.Command("fortune", "-f").CombinedOutput(); err != nil {
		fmt.Println(err)
		return
	} else {
		execOutput = string(out)
	}
	// expected format
	// 	100.00% /usr/local/Cellar/fortune/9708/share/games/fortunes
	//  7.65% computers
	//  1.01% riddles
	//  4.36% men-women
	//  1.92% literature
	//  1.13% love
	//  0.22% magic
	//  0.79% linuxcookie
	//  1.56% drugs
	//  0.38% pets
	//  3.44% art
	//  1.51% law
	//  0.40% goedel
	//  1.52% education
	//  1.22% ethnic
	//  4.66% science
	//  0.07% ascii-art
	//  4.82% miscellaneous
	//  1.10% sports
	//  4.10% zippy
	//  5.18% politics
	//  1.69% startrek
	//  3.01% wisdom
	//  0.40% news
	//  4.72% work
	//  0.54% medicine
	//  9.22% people
	//  1.48% food
	//  1.47% humorists
	//  3.72% platitudes
	//  8.54% cookie
	//  5.38% songs-poems
	//  8.28% definitions
	//  1.12% kids
	//  3.24% fortunes
	//  0.14% translate-me

	temp := strings.Split(string(execOutput), "\n")
	path := ""
	fortuneFiles := ""
	for i := 0; i < len(temp); i++ {
		w := strings.Split(strings.TrimLeft(temp[i], " "), " ")
		if len(w) != 2 {
			continue
		}
		if i == 0 {
			path = w[1]
		} else {
			db := path + "/" + w[1]
			// fmt.Println(db)
			dbBytes, err := ioutil.ReadFile(db)
			if err != nil {
				log.Fatalf("Can't read file:(%v)  err=> %v", db, err.Error())
				return
			}
			fortuneFiles += string(dbBytes)
		}
	}
	// fmt.Println(fortuneFiles)
	fortuneData = strings.Split(fortuneFiles, "%")
	// for i := 0; i < len(fortuneData); i++ {
	// 	fmt.Printf("[%v]\n%v\n", i, fortuneData[i])
	// }
	curFortuneIndex = rand.Intn(len(fortuneData))
}

// GetFortuneLen : get number of fortune data
func GetFortuneLen() int {
	return len(fortuneData)
}

// GetCurFortuneIndex : get current fortune index
func GetCurFortuneIndex() int {
	return curFortuneIndex
}

// GetNextFortuneIndex : get next fortune index
func GetNextFortuneIndex() int {
	curFortuneIndex++
	if curFortuneIndex >= len(fortuneData) {
		curFortuneIndex = 0
	}
	return curFortuneIndex
}

// GetPreFortuneIndex : get previous fortune index
func GetPreFortuneIndex() int {
	curFortuneIndex--
	if curFortuneIndex < 0 {
		curFortuneIndex = len(fortuneData) - 1
	}
	return curFortuneIndex
}

// GetPreFortune : get previous fortune
func GetPreFortune() []string {
	if len(fortuneData) > 0 {
		return strings.Split(strings.TrimPrefix(fortuneData[GetPreFortuneIndex()], "\n"), "\n")
	}
	return nil
}

// GetNextFortune : get next fortune
func GetNextFortune() []string {
	if len(fortuneData) > 0 {
		return strings.Split(strings.TrimPrefix(fortuneData[GetNextFortuneIndex()], "\n"), "\n")
	}
	return nil
}
