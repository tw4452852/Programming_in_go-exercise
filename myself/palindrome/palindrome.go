package main

import (
	"fmt"
	"os"
	"path/filepath"
	"unicode/utf8"
)

var IsPalindrome func(string) bool

func init() {
	if len(os.Args) > 1 && (os.Args[1] == "-a" || os.Args[1] == "--ascii") {
		os.Args = append(os.Args[:1], os.Args[2:]...)
		IsPalindrome = func(str string) bool {
			j := len(str) - 1
			for i := 0; i < len(str) / 2; i++ {
				if str[i] != str[j] {
					return false
				}
				j--
			}
			return true
		}
	} else {
		IsPalindrome = func(str string) bool {
			for len(str) > 0 {
				first, sizeOfFirst := utf8.DecodeRuneInString(str)
				if sizeOfFirst == len(str) {
					break
				}
				last, sizeOfLast := utf8.DecodeLastRuneInString(str)
				if first != last {
					return false
				}
				str = str[sizeOfLast : len(str) - sizeOfLast]
			}
			return true
		}
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("usage: %s [-a|--ascii] word1 [word2] [...wordN]\n",
			filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	words := os.Args[1:]
	for _, word := range words {
		fmt.Printf("%-5t %q\n", IsPalindrome(word), word)
	}
}
