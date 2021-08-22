package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strings"
	"sync"
)

type SharedMaps struct {
	x string
	y int
}

var mu sync.Mutex

func countWords(arr []string, c chan int, sharedMap map[string]int) {
	var localMap map[string]int
	localMap = make(map[string]int)
	for _, i := range arr {
		_, ok := localMap[i]
		if ok {
			localMap[i] += 1
		} else {
			localMap[i] = 1
		}
	}
	mu.Lock()
	for key, value := range localMap {
		_, ok := sharedMap[key]
		if ok {
			sharedMap[key] += value
		} else {
			sharedMap[key] = value
		}
	}
	mu.Unlock()
	c <- 1
}

func reducer(arr []string, sharedMap map[string]int, c chan int) {

	j := 0
	for i := range c {
		j += i
		if j == 5 {
			close(c)
		}
	}

	sorted := []SharedMaps{}
	for k, v := range sharedMap {
		sorted = append(sorted, SharedMaps{x: k, y: v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].y == sorted[j].y {
			return sorted[i].x < sorted[j].x
		} else {
			return sorted[i].y > sorted[j].y

		}
	})

	f, _ := os.Create("WordCountOutput.txt")
	defer f.Close()
	for _, v := range sorted {
		f.WriteString(v.x + " : " + fmt.Sprint(v.y) + "\n")
	}

}

func main() {
	fmt.Print("Enter file path: ")
	var inputFile string
	fmt.Scanln(&inputFile)

	file, _ := ioutil.ReadFile(inputFile)
	text := string(file)
	words := strings.Split(strings.ToLower(strings.ReplaceAll(text, "\r\n", " ")), " ")

	var sharedMap map[string]int
	sharedMap = make(map[string]int)
	slice1 := int(math.Ceil(float64(len(words)) / 5.0))
	c := make(chan int)

	go countWords(words[:slice1], c, sharedMap)
	go countWords(words[slice1:slice1*2], c, sharedMap)
	go countWords(words[slice1*2:slice1*3], c, sharedMap)
	go countWords(words[slice1*3:slice1*4], c, sharedMap)
	go countWords(words[slice1*4:], c, sharedMap)

	reducer(words, sharedMap, c)

}
