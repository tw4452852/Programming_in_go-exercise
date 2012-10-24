package oslice

import (
	"strings"
	"fmt"
)

type Item interface{}

type lessfunc func(Item, Item) bool

type Slice struct {
	items	[]Item
	less	lessfunc
}

func New(f lessfunc) *Slice {
	return &Slice{items: make([]Item, 0), less: f}
}

func NewStringSlice() *Slice {
	return &Slice{items: make([]Item, 0), less: func(a, b Item) bool {
		return strings.ToLower(a.(string)) < strings.ToLower(b.(string))
	}}
}

func NewIntSlice() *Slice{
	return &Slice{items: make([]Item, 0), less: func(a, b Item) bool {
		return a.(int) < b.(int)
	}}
}

func (s *Slice) Clear() {
	s.items = s.items[:0]
}

func (s *Slice) Add(a Item) {
	index := bisectLeft(s.items, s.less, a)
	if index >= len(s.items) {
		s.items = append(s.items, a)
	} else {
		s.items = insert(s.items, index, a)
	}
}

func (s *Slice) Remove(a Item) int {
	index := bisectLeft(s.items, s.less, a)
	if index >= len(s.items) || s.less(s.items[index], a) {
		return -1
	}
	s.items = remove(s.items, index)
	return index
}

func (s *Slice) Index(a Item) int {
	index := bisectLeft(s.items, s.less, a)
	if index >= len(s.items) || s.less(s.items[index], a) {
		return -1
	}
	return index
}

func (s *Slice) At(index int) Item {
	if index >= len(s.items) {
		panic("out of range")
	}
	return s.items[index]
}

func (s *Slice) Len() int {
	return len(s.items)
}

func (s *Slice) String() string {
	result := ""
	for i, value := range s.items {
		result += fmt.Sprintf("[%d : ", i)
		switch value.(type) {
		case int, int64:
			result += fmt.Sprintf("%d", value.(int))
		default:
			result += fmt.Sprintf("%v", value)
		}
		result += "] "
	}
	return result
}

func remove(items []Item, index int) []Item {
	result := items[:index]
	result = append(result, items[index + 1:]...)
	return result
}

func insert(items []Item, index int, a Item) []Item {
	result := make([]Item, 0, len(items) + 1)
	result = append(result, items[:index]...)
	result = append(result, a)
	result = append(result, items[index:]...)
	return result
}

func bisectLeft(slice []Item, less lessfunc, x Item) int {
	left, right := 0, len(slice)
	for left < right {
		middle := int((left + right) / 2)
		if less(slice[middle], x) {
			left = middle + 1
		} else {
			right = middle
		}
	}
	return left
}
