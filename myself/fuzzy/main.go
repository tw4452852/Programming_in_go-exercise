package main

import (
	"github.com/tw4452852/Programming_in_go-exercise/myself/fuzzy/fuzzybool"
	"fmt"
)

func main() {
	a, _ := fuzzybool.New(0)
	b, _ := fuzzybool.New(.25)
	c, _ := fuzzybool.New(.75)
	d := c.Copy()
	if err := d.Set(1); err != nil {
		fmt.Println(err)
	}
	process(a, b, c, d)
	s := []*fuzzybool.FuzzyBool{a, b, c, d}
	fmt.Println(s)
}

func process(a, b, c, d *fuzzybool.FuzzyBool) {
    fmt.Println("Original:", a, b, c, d)
    fmt.Println("Not:     ", a.Not(), b.Not(), c.Not(), d.Not())
    fmt.Println("Not Not: ", a.Not().Not(), b.Not().Not(), c.Not().Not(),
        d.Not().Not())
    fmt.Print("0.And(.25)→", a.And(b), "• .25.And(.75)→", b.And(c),
        "• .75.And(1)→", c.And(d), " • .25.And(.75,1)→", b.And(c, d), "\n")
    fmt.Print("0.Or(.25)→", a.Or(b), "• .25.Or(.75)→", b.Or(c),
        "• .75.Or(1)→", c.Or(d), " • .25.Or(.75,1)→", b.Or(c, d), "\n")
    fmt.Println("a < c, a == c, a > c:", a.Less(c), a.Equal(c), c.Less(a))
    fmt.Println("Bool:    ", a.Bool(), b.Bool(), c.Bool(), d.Bool())
    fmt.Println("Float:   ", a.Float(), b.Float(), c.Float(), d.Float())
}
