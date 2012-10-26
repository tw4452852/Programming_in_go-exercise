package main

import (
	"runtime"
	"path/filepath"
	"os"
	"io"
	"fmt"
	"log"
	"bufio"
	"regexp"
	"github.com/tw4452852/Programming_in_go-exercise/myself/smap"
)

var workers = runtime.NumCPU()

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(os.Args) != 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s <file.log>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	lines := make(chan string, workers * 4)
	done := make(chan struct{}, workers)
	pageMap := smap.New()
	go readLines(os.Args[1], lines)
	processLines(done, pageMap, lines)
	waitUntil(done)
	showResults(pageMap)
}

func readLines(filename string, lines chan<- string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("failed to open the file:", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if line != "" {
			lines <- line
		}
		if err != nil {
			if err != io.EOF {
				log.Println("failed to finish reading the file:", err)
			}
			break
		}
	}
	close(lines)
}

func processLines(done chan<- struct{}, pageMap smap.SafeMap, lines <-chan string) {
	getRx := regexp.MustCompile(`GET[ \t]+([^ \t\n]+[.]html?)`)
	incrementer := func(value interface{}, found bool) interface{} {
		if found {
			return value.(int) + 1
		}
		return 1
	}
	for i := 0; i < workers; i++ {
		go func() {
			for line := range lines {
				if matches := getRx.FindStringSubmatch(line); matches != nil {
					pageMap.Update(matches[1], incrementer)
				}
			}
			done <- struct{}{}
		}()
	}
}

func waitUntil(done <-chan struct{}) {
	for i := 0; i < workers; i++ {
		<-done
	}
}

func showResults(pageMap smap.SafeMap) {
	pages := pageMap.Close()
	for page, count := range pages {
		fmt.Printf("%8d %s\n", count, page)
	}
}
