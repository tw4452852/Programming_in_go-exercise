package main

import (
	"fmt"
	"os"
	"log"
	"path/filepath"
	"image"
	"sync"
	"runtime"
    _ "image/gif"
    _ "image/jpeg"
    _ "image/png"
)

var workers = runtime.NumCPU()

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s <file>...\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	files := parseCommandline(os.Args[1:])
	jobs := make(chan string, len(files))
	results := make(chan string, workers * 16)
	go addJobs(jobs, files)
	go processImage(results, jobs)
	printTags(results)
}

func parseCommandline(files []string) []string {
	if runtime.GOOS == "windows" {
		args := make([]string, 0, len(files))
		for _, name := range files {
			if matches, err := filepath.Glob(name); err != nil {
				args = append(args, name)
			} else if matches != nil {
				args = append(args, matches...)
			}
		}
		return args
	}
	return files
}

func addJobs(jobs chan<- string, files []string) {
	for _, file := range files {
		jobs <- file
	}
	close(jobs)
}

func processImage(results chan<- string, jobs <-chan string) {
	waiter := &sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		waiter.Add(1)
		go getImage(results, jobs, func() { waiter.Done() })
	}
	waiter.Wait()
	close(results)
}

func getImage(results chan<- string, jobs <-chan string, done func()) {
	defer done()
	for filename := range jobs {
		file, err := os.Open(filename)
		if err != nil {
			log.Println("error:", err)
			continue
		}
		defer file.Close()
		config, _, err := image.DecodeConfig(file)
		if err != nil {
			log.Println("error:", err)
			continue
		}
		results <- fmt.Sprintf(`<img src=%q width="%d" heigh="%d" />`,
			filepath.Base(filename), config.Width, config.Height)
	}
}

func printTags(results <-chan string) {
	for tag := range results {
		fmt.Println(tag)
	}
}
