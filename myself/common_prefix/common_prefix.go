package main

import (
	"fmt"
	"bytes"
	"strings"
	"path/filepath"
)

func main() {
    testData := [][]string{
        {"/home/user/goeg", "/home/user/goeg/prefix",
            "/home/user/goeg/prefix/extra"},
        {"/home/user/goeg", "/home/user/goeg/prefix",
            "/home/user/prefix/extra"},
        {"/pecan/π/goeg", "/pecan/π/goeg/prefix",
            "/pecan/π/prefix/extra"},
        {"/pecan/π/circle", "/pecan/π/circle/prefix",
            "/pecan/π/circle/prefix/extra"},
        {"/home/user/goeg", "/home/users/goeg",
            "/home/userspace/goeg"},
        {"/home/user/goeg", "/tmp/user", "/var/log"},
        {"/home/mark/goeg", "/home/user/goeg"},
        {"home/user/goeg", "/tmp/user", "/var/log"},
    }
	for _, slice := range testData {
		fmt.Printf("[")
		gap := ""
		for _, data := range slice {
			fmt.Printf("%s%q", gap, data)
			gap = " "
		}
		fmt.Println("]")
		cp := CommonPrefix(slice)
		cpp := CommonPathPrefix(slice)
		equal := "=="
		if cpp != cp {
			equal = "!="
		}
		fmt.Printf("char x path prefix: %q %s %q\n\n", cp, equal, cpp)
	}
}

func CommonPrefix(slice []string) string {
	components := make([][]rune, len(slice))
	for i, text := range slice {
		components[i] = []rune(text)
	}
	if len(components) == 0 || len(components[0]) == 0 {
		return ""
	}
	var common bytes.Buffer
FINISH:
	for column := 0; column < len(components[0]); column++ {
		char := components[0][column]
		for row := 1; row < len(components); row++ {
			if column >= len(components[row]) || components[row][column] != char {
				break FINISH
			}
		}
		common.WriteRune(char)
	}
	return common.String()
}

func CommonPathPrefix(slice []string) string {
	const separator = string(filepath.Separator)
	components := make([][]string, len(slice))
	for i, path := range slice {
		components[i] = strings.Split(path, separator)
		if strings.HasPrefix(path, separator) {
			components[i] = append([]string{separator}, components[i]...)
		}
	}
	if len(components) == 0 || len(components[0]) == 0 {
		return ""
	}
	var common []string
FINISH:
	for colum := range components[0] {
		part := components[0][colum]
		for row := 1; row < len(components); row++ {
			if len(components[row]) == 0 ||
				colum >= len(components[row]) ||
				components[row][colum] != part {
					break FINISH
			}
		}
		common = append(common, part)
	}
	return filepath.Join(common...)
}
