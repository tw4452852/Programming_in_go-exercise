package main

import (
	"fmt"
	"os"
	"io"
	"log"
	"bufio"
	"strings"
	"path/filepath"
)

func main() {
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s file\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	separators := []string{"\t", "*", "|", " "}

	linesRead, lines := readUpToNLine(os.Args[1], 5)
	counts := createCounts(lines, separators, linesRead)
	separator := guessSep(counts, separators, linesRead)
	report(separator)
}

func readUpToNLine(filename string, maxLines int) (int, []string) {
	var (
		file	*os.File
		err		error
	)
	if file, err = os.Open(filename); err != nil {
		log.Fatal("failed to open the file: ", err)
	}
	defer file.Close()

	lines := make([]string, maxLines)
	reader := bufio.NewReader(file)
	i := 0
	for ; i < maxLines; i++ {
		line, err := reader.ReadString('\n')
		if line != "" {
			lines[i] = line
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal("failed to finish reading the file: ", err)
		}
	}
	return i, lines[:i]
}

func createCounts(lines, seperators []string, linesRead int) [][]int {
	counts := make([][]int, len(seperators))
	for sepIndex := range seperators {
		counts[sepIndex] = make([]int, linesRead)
		for lineIndex, line := range lines {
			counts[sepIndex][lineIndex] = strings.Count(line, seperators[sepIndex])
		}
	}
	return counts
}

func guessSep(counts [][]int, seperators []string, linesRead int) string {
	for sepIndex := range seperators {
		same := true
		count := counts[sepIndex][0]
		for lineIndex := 1; lineIndex < linesRead; lineIndex++ {
			if counts[sepIndex][lineIndex] != count {
				same = false
				break
			}
		}
		if count > 0 && same == true {
			return seperators[sepIndex]
		}
	}
	return ""
}

func report(separator string) {
	switch separator {
	case "":
		fmt.Println("whitespace-separated or not separated at all")
	case "\t":
		fmt.Println("tab-separated")
	default:
		fmt.Printf("%s-separated\n", separator)
	}
}
