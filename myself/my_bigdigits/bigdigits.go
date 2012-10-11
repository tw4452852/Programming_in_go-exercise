package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"flag"
	"strings"
)

var bigDigits = [][]string{
    {"  000  ",
     " 0   0 ",
     "0     0",
     "0     0",
     "0     0",
     " 0   0 ",
     "  000  "},
    {" 1 ", "11 ", " 1 ", " 1 ", " 1 ", " 1 ", "111"},
    {" 222 ", "2   2", "   2 ", "  2  ", " 2   ", "2    ", "22222"},
    {" 333 ", "3   3", "    3", "  33 ", "    3", "3   3", " 333 "},
    {"   4  ", "  44  ", " 4 4  ", "4  4  ", "444444", "   4  ",
        "   4  "},
    {"55555", "5    ", "5    ", " 555 ", "    5", "5   5", " 555 "},
    {" 666 ", "6    ", "6    ", "6666 ", "6   6", "6   6", " 666 "},
    {"77777", "    7", "   7 ", "  7  ", " 7   ", "7    ", "7    "},
    {" 888 ", "8   8", "8   8", " 888 ", "8   8", "8   8", " 888 "},
    {" 9999", "9   9", "9   9", " 9999", "    9", "    9", "    9"},
}

var (
	help	*bool
	bar		*bool
)

func init() {
	help = flag.Bool("h", false, "show usage help")
	bar = flag.Bool("b", false, "show bar")
}

func main() {
	flag.Parse()

	if len(os.Args) == 1 || *help {
		fmt.Printf("usage: %s [-h] [-b] <whole-number>\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(1)
	}

	stringOfDigits := flag.Arg(0)
	for row := range bigDigits[0] {
		line := ""
		for column := range stringOfDigits {
			digit := stringOfDigits[column] - '0'
			if 0 <= digit && digit <= 9 {
				line += bigDigits[digit][row] + " "
			} else {
				log.Fatal("invalid whole number")
			}
		}
		if *bar && row == 0 {
			fmt.Println(strings.Repeat("*", len(line)))
		}
		fmt.Println(line)
		if *bar && row + 1 == len(bigDigits[0]) {
			fmt.Println(strings.Repeat("*", len(line)))
		}
	}
}
