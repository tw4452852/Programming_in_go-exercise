package font

import (
	"log"
	"fmt"
	"unicode/utf8"
)

type Font struct {
	family	string
	size	int
}

func New(family string, size int) *Font {
	return &Font{validFamily("sans-family", family), validSize(10, size)}
}

func (f *Font) SetFamily(family string) {
	f.family = validFamily(f.family, family)
}

func (font *Font) Size() int { return font.size }

func (font *Font) Family() string { return font.family }

func (f *Font) SetSize(size int) {
	f.size = validSize(f.size, size)
}

func (f *Font) String() string {
	return fmt.Sprintf("{font-family: %q; font-size: %dpt;}", f.family, f.size)
}

func validSize(oldSize, newSize int) int {
	if newSize < 5 || 144 < newSize {
		log.Printf("font.validSize(): ignored invalied size %d", newSize)
		return oldSize
	}
	return newSize
}

func validFamily(oldFamily, newFamily string) string {
	if len(newFamily) < utf8.UTFMax &&
		utf8.RuneCountInString(newFamily) < 1 {
		log.Printf("font.validFamily(): ignored invalied family %s", newFamily)
		return oldFamily
	}
	return newFamily
}
