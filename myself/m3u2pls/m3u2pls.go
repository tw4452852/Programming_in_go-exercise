package main

import (/*{{{*/
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
	"strings"
	"strconv"
	"log"
)/*}}}*/

type Song struct {/*{{{*/
	Title		string
	Filename	string
	Seconds		int
}/*}}}*/

func readM3uPlaylist(data string) (songs []Song) {/*{{{*/
	var song Song
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#EXTM3U") {
			continue
		}
		if strings.HasPrefix(line, "#EXTINF:") {
			song.Title, song.Seconds = parseExtinfLine(line)
		} else {
			song.Filename = strings.Map(mapPlatformDirSeparator, line)
		}
		if song.Filename != "" && song.Title != "" && song.Seconds != 0 {
			songs = append(songs, song)
			song = Song{}
		}
	}
	return songs
}/*}}}*/

func parseExtinfLine(line string) (title string, seconds int) {/*{{{*/
	if i := strings.IndexAny(line, "-0123456789"); i > -1 {
		const separator = ","
		line = line[i:]
		if j := strings.Index(line, separator); j > -1 {
			title = line[j + len(separator):]
			var err error
			if seconds, err = strconv.Atoi(line[:j]); err != nil {
				log.Printf("failed to read the duration for '%s': %v\n", title, err)
				seconds = -1
			}
		}
	}
	return title, seconds
}/*}}}*/

func mapPlatformDirSeparator(char rune) rune {/*{{{*/
	if char == '/' || char == '\\' {
		return filepath.Separator
	}
	return char
}/*}}}*/

func writePlsPlaylist(songs []Song) {/*{{{*/
	fmt.Println("[playlist]")
	for i, song := range songs {
		i++
		fmt.Printf("File%d=%s\n", i, song.Filename)
		fmt.Printf("Tile%d=%s\n", i, song.Title)
		fmt.Printf("Length%d=%d\n", i, song.Seconds)
	}
	fmt.Printf("NumberOfEntries=%d\nVersion=2\n", len(songs))
}/*}}}*/

func readPlsPlaylist(data string) (songs []Song) {/*{{{*/
	var song Song
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "[playlist]") || strings.HasPrefix(line, "NumberOfEntries") || strings.HasPrefix(line, "Version") {
			continue
		}
		if strings.HasPrefix(line, "Title") {
			song.Title = parseTitleLine(line)
		} else if strings.HasPrefix(line, "Length") {
			song.Seconds = parseLengthLine(line)
		} else if strings.HasPrefix(line, "File") {
			song.Filename = parseFilenameLine(line)
		}
		if song.Filename != "" && song.Title != "" && song.Seconds != 0 {
			songs = append(songs, song)
			song = Song{}
		}
	}
	return songs
}/*}}}*/

func parseFilenameLine(line string) (filename string) {/*{{{*/
	const equ = "="
	if i := strings.Index(line, equ); i > -1 {
		filename = line[i + len(equ):]
	}
	if filename != "" {
		filename = strings.Map(mapPlatformDirSeparator, filename)
	}
	return filename
}/*}}}*/

func parseTitleLine(line string) (title string) {/*{{{*/
	const equ = "="
	if i := strings.Index(line, equ); i > -1 {
		title = line[i + len(equ):]
	}
	return title
}/*}}}*/

func parseLengthLine(line string) (seconds int) {/*{{{*/
	const equ = "="
	if i := strings.Index(line, equ); i > -1 {
		var err error
		if seconds, err = strconv.Atoi(line[i + len(equ):]); err != nil {
			log.Printf("failed to read durating '%s' : %v\n", line, err)
			seconds = -1
		}
	}
	return seconds
}/*}}}*/

func writeM3uPlaylist(songs []Song) {/*{{{*/
	fmt.Println("#EXTM3U")
	for _, song := range songs {
		fmt.Printf("#EXTINF:%d,%s\n", song.Seconds, song.Title)
		fmt.Println(song.Filename)
	}
}/*}}}*/

func main() {/*{{{*/
	if len(os.Args) == 1 || !(strings.HasSuffix(os.Args[1], ".m3u") || strings.HasSuffix(os.Args[1], ".pls")) {
		fmt.Printf("usage: %s <file.[m3u|pls]>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	if rawBytes, err := ioutil.ReadFile(os.Args[1]); err != nil {
		log.Fatal(err)
	} else {
		if strings.HasSuffix(os.Args[1], ".m3u") {
			songs := readM3uPlaylist(string(rawBytes))
			writePlsPlaylist(songs)
		} else {
			songs := readPlsPlaylist(string(rawBytes))
			writeM3uPlaylist(songs)
		}
	}
}/*}}}*/
