package main

import (
	"fmt"
	"strings"
	"log"
	"sort"
)

func cleanDuplicate(data []int) (result []int) {
	seen := map[int]bool{}
	for _, value := range data {
		if _, found := seen[value]; !found {
			result = append(result, value)
			seen[value] = true
		}
	}
	return result
}

func flatten(matrix [][]int) (result []int) {
	for _, row := range matrix {
		result = append(result, row...)
	}
	return result
}

func make2D(slice []int, column int) [][]int {
	cnt := len(slice) / column
	if len(slice) % column != 0 {
		cnt++
	}
	result := make([][]int, cnt)
	for i := 0; i < cnt; i++ {
		result[i] = makeRow(slice, i * column, column)
	}
	return result
}

func makeRow(slice []int, start, length int) []int {
	result := make([]int, 0, length)
	for i := 0; i < length; i++ {
		if i + start >= len(slice) {
			break
		}
		result = append(result, slice[i + start])
	}
	return result
}

func parseIni(ini []string) map[string]map[string]string {
	const separator = "="
	result := make(map[string]map[string]string)
	group := "General"
	for _, line := range ini {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, ";") || line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			group = line[1 : len(line) - 1]
		} else if i := strings.Index(line, separator); i > -1 {
			key := line[:i]
			value := line[i + len(separator) :]
			if _, found := result[group]; !found {
				result[group] = make(map[string]string)
			}
			result[group][key] = value
		} else {
			log.Print("failed to parse line: ", line)
		}
	}
	return result
}

func printInt(ini map[string]map[string]string) {
	groups := make([]string, 0, len(ini))
	for group := range ini {
		groups = append(groups, group)
	}
	sort.Strings(groups)
	for i, group := range groups {
		fmt.Printf("[%s]\n", group)
		keys := make([]string, 0, len(ini[group]))
		for key := range ini[group] {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			fmt.Printf("%s=%s\n", key, ini[group][key])
		}
		if i + 1 < len(groups) {
			fmt.Println()
		}
	}
}

func main() {
	fmt.Println(cleanDuplicate([]int{9, 1, 9, 5, 4, 4, 2, 1, 5, 4, 8, 8, 4, 3, 6, 9, 5, 7, 5}))

	irrMatrix := [][]int{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11},
		{12, 13, 14, 15},
	}
	slice := flatten(irrMatrix)
	fmt.Printf("1x%d: %v\n", len(slice), slice)

	var column int = 5
	fmt.Printf("%d %v\n", column, make2D(slice, column))

    iniData := []string{
        "; Cut down copy of Mozilla application.ini file",
        "",
        "[App]",
        "Vendor=Mozilla",
        "Name=Iceweasel",
        "Profile=mozilla/firefox",
        "Version=3.5.16",
        "[Gecko]",
        "MinVersion=1.9.1",
        "MaxVersion=1.9.1.*",
        "[XRE]",
        "EnableProfileMigrator=0",
        "EnableExtensionManager=1",
    }
	result := parseIni(iniData)
	printInt(result)
}
