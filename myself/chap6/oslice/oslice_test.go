package oslice_test

import (
    "fmt"
	"testing"
	"github.com/tw4452852/Programming_in_go-exercise/myself/chap6/oslice"
)

func TestOslice_string(t *testing.T) {
	s := oslice.NewStringSlice()
	s.Add("hello")
	s.Add("tw")
	if s.Len() != 2 {
		t.Fatal("string Len(): failed")
	}
	if s.At(0) != "hello" {
		t.Fatal("string At(): failed")
	}
	if s.Remove("tw") != 1 {
		t.Fatal("string Remove(): failed")
	}
	if s.Remove("tw") != -1 {
		t.Fatal("string Remove(): failed")
	}
	s.Clear()
	if s.Len() != 0 {
		t.Fatal("string Clear(): failed")
	}
	s.Add("b")
	s.Add("a")
	s.Add("a")
	if s.Index("a") != 0 {
		t.Fatal("string index(): failed")
	}
	fmt.Println(s)
}

func TestOslice_int(t *testing.T) {
	i := oslice.NewIntSlice()
	i.Add(10)
	i.Add(3)
	if i.Len() != 2 {
		t.Fatal("int Len(): failed")
	}
	if i.At(0) != 3 {
		t.Fatal("int At(): failed")
	}
	if i.Remove(10) != 1 {
		t.Fatal("int Remove(): failed")
	}
	if i.Remove(10) != -1 {
		t.Fatal("int Remove(): failed")
	}
	i.Clear()
	if i.Len() != 0 {
		t.Fatal("int Clear(): failed")
	}
	i.Add(2)
	i.Add(1)
	i.Add(1)
	if i.Index(1) != 0 {
		t.Fatal("int Index(): failed")
	}

	fmt.Println(i)
}

type mys struct {
	v	int
}

func TestOslice_general(t *testing.T) {
	g := oslice.New(func(a, b oslice.Item) bool {
		return a.(mys).v < b.(mys).v
	})

	g.Add(mys{10})
	g.Add(mys{3})
	if g.Len() != 2 {
		t.Fatal("general Len(): failed")
	}
	if g.At(0).(mys).v != 3 {
		t.Fatal("general At(): failed")
	}
	if g.Remove(mys{10}) != 1 {
		t.Fatal("general Remove(): failed")
	}
	if g.Remove(mys{10}) != -1 {
		t.Fatal("general Remove(): failed")
	}
	g.Clear()
	if g.Len() != 0 {
		t.Fatal("general Clear(): failed")
	}
	g.Add(mys{2})
	g.Add(mys{1})
	g.Add(mys{1})
	if g.Index(mys{1}) != 0 {
		t.Fatal("general Index(): failed")
	}
	fmt.Println(g)
}
