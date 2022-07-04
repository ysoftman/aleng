// example
package main

import (
	"math/rand"
	"regexp"
	"strings"
)

// ExampleCmdText : example command text
const ExampleCmdText = "Examples, pre (ctrl+u) / next (ctrl+d)"

var examples []string
var curExampleIndex int

// ReadExampleRawData : read the examples file
func ReadExampleRawData() {
	examples = strings.Split(string(exampleRawData), "--")
	curExampleIndex = rand.Intn(len(examples))
}

// GetExamplesLen : get number of examples
func GetExamplesLen() int {
	return len(examples)
}

// GetCurExampleIndex : get current example index
func GetCurExampleIndex() int {
	return curExampleIndex
}

// GetNextExampleIndex : get next example index
func GetNextExampleIndex() int {
	curExampleIndex++
	if curExampleIndex >= len(examples) {
		curExampleIndex = 0
	}
	return curExampleIndex
}

// GetPreExampleIndex : get previous example index
func GetPreExampleIndex() int {
	curExampleIndex--
	if curExampleIndex < 0 {
		curExampleIndex = len(examples) - 1
	}
	return curExampleIndex
}

// GetPreExample : get previous example
func GetPreExample() []string {
	if len(examples) > 0 {
		return strings.Split(strings.TrimPrefix(examples[GetPreExampleIndex()], "\n"), "\n")
	}
	return nil
}

// GetNextExample : get next example
func GetNextExample() []string {
	if len(examples) > 0 {
		return strings.Split(strings.TrimPrefix(examples[GetNextExampleIndex()], "\n"), "\n")
	}
	return nil
}

// FindExample : find example including keyword
func FindExample(keyword string) []string {
	keyword = strings.ToLower(keyword)
	var foundExamples []string
	for i := range examples {
		// if strings.Contains(examples[i], keyword) {
		// 	return strings.Split(strings.TrimPrefix(examples[i], "\n"), "\n")
		// }
		// https://github.com/google/re2/wiki/Syntax
		// \b at ascii word boundary
		if matched, _ := regexp.MatchString("\\b"+keyword+"\\b", strings.ToLower(examples[i])); matched {
			fb := strings.Split(strings.TrimPrefix(examples[i], "\n"), "\n")
			for j := range fb {
				foundExamples = append(foundExamples, fb[j])
				curExampleIndex = i
			}
		}
	}
	return foundExamples
}
