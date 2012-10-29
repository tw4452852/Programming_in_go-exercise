package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"log"
	"runtime"
	"path/filepath"
	"regexp"
	"strings"
	"image"
	"sync"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var workers = runtime.NumCPU()
const (
	widthAttr	= "width="
	heightAttr	= "height="
)

var (
	imageRx	*regexp.Regexp
	srcRx	*regexp.Regexp
)

func init() {
	imageRx = regexp.MustCompile(`<[iI][mM][gG][^>]+>`)
	srcRx	= regexp.MustCompile(`src=["']([^"']+)['"]`)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s <htmlFile1>...<htmlFileN>\n",
			filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	files := parseCmdFiles(os.Args[1:])
	jobs := make(chan string, workers * 16)
	go addJobs(jobs, files)
	waiter := &sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		waiter.Add(1)
		go imageSize(jobs, func() { waiter.Done() })
	}
	waiter.Wait()
}

func parseCmdFiles(files []string) []string {
	if runtime.GOOS == "windows" {
		args := make([]string, 0, len(files))
		for _, file := range files {
			if matches, err := filepath.Glob(file); err != nil {
				args = append(args, file)
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

func imageSize(jobs <-chan string, done func()) {
	defer done()
	for filename := range jobs {
		if info, err := os.Stat(filename); err != nil ||
			(info.Mode() & os.ModeType == 1) {
			log.Println("ignoring:", filename)
			continue
		}
		raw, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Println("failed to read:", err)
			continue
		}
		html := string(raw)
		dir, _ := filepath.Split(filename)
		newHtml := imageRx.ReplaceAllStringFunc(html, makeSizerFunc(dir))
		if len(html) != len(newHtml) {
			file, err := os.Create(filename)
			if err != nil {
				log.Printf("couldn't update %s: %v\n", filename, err)
				continue
			}
			defer file.Close()
			if _, err := file.WriteString(newHtml); err != nil {
				fmt.Printf("error when updating %s: %v\n", filename, err)
			}
		}
	}
}

func makeSizerFunc(dir string) func(string) string {
	return func(originalTag string) string {
		tag := originalTag
		if strings.Contains(tag, widthAttr) &&
			strings.Contains(tag, heightAttr) {
				return tag
		}
		match := srcRx.FindStringSubmatch(tag)
		if match == nil {
			fmt.Println("can't find <img>'s src attribute", tag)
			return tag
		}
		filename := match[1]
		if !filepath.IsAbs(filename) {
			filename = filepath.Join(dir, filename)
		}
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println("can't open image to read its size:", err)
			return tag
		}
		defer file.Close()
		config, _, err := image.DecodeConfig(file)
		if err != nil {
			fmt.Println("can't ascertain the image's size:", err)
			return tag
		}
		tag, end := tagEnd(tag)
		if !strings.Contains(tag, widthAttr) {
			tag += fmt.Sprintf(` %s"%d"`, widthAttr, config.Width)
		}
		if !strings.Contains(tag, heightAttr) {
			tag += fmt.Sprintf(` %s"%d"`, heightAttr, config.Height)
		}
		tag += end
		return tag
	}
}

func tagEnd(originalTag string) (tag, end string) {
	end = ">"
	tag = originalTag[:len(originalTag) - 1]
	if tag[len(tag) - 1] == '/' {
		end = " />"
		tag = tag[:len(tag) - 1]
	}
	return strings.TrimSpace(tag), end
}
