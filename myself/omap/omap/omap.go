package omap

import (
	"strings"
)

type Map struct {
	root	*node
	less	func(interface{}, interface{}) bool
	length	int
}

type node struct {
	key, value	interface{}
	red			bool
	left, right	*node
}

func New(less func(interface{}, interface{}) bool) *Map {
	return &Map{less:less}
}

func NewCaseFoldedKeyed() *Map {
	return &Map{less: func(a, b interface{}) bool {
		return strings.ToLower(a.(string)) < strings.ToLower(b.(string))
	}}
}

func NewIntKeyed() *Map {
	return &Map{less: func(a, b interface{}) bool {
		return a.(int) < b.(int)
	}}
}

func NewStringKeyed() *Map {
	return &Map{less: func(a, b interface{}) bool {
		return a.(string) < b.(string)
	}}
}

func NewFloat64Keyed() *Map {
	return &Map{less: func(a, b interface{}) bool {
		return a.(float64) < b.(float64)
	}}
}

func (m *Map) Insert(key, value interface{}) (inserted bool) {
	m.root, inserted = m.insert(m.root, key, value)
	m.root.red = false
	if inserted {
		m.length++
	}
	return inserted
}

func (m *Map) insert(root *node, key, value interface{}) (*node, bool) {
	inserted := false
	if root == nil {
		return &node{key: key, value: value, red: true}, true
	}
	if isRed(root.left) && isRed(root.right) {
		colorFlip(root)
	}
	if m.less(key, root.key) {
		root.left, inserted = m.insert(root.left, key, value)
	} else if m.less(root.key, key) {
		root.right, inserted = m.insert(root.right, key, value)
	} else {
		root.value = value
	}
	if isRed(root.right) && !isRed(root.left) {
		root = rotateLeft(root)
	}
	if isRed(root.left) && isRed(root.left.left) {
		root = rotateRight(root)
	}
	return root, inserted
}

func isRed(root *node) bool {
	return root != nil && root.red
}

func colorFlip(root *node) {
	root.red = !root.red
	if root.left != nil {
		root.left.red = !root.left.red
	}
	if root.right != nil {
		root.right.red = !root.right.red
	}
}

func rotateLeft(root *node) *node {
	x := root.right
	root.right = x.left
	x.left = root
	x.red = root.red
	root.red = true
	return x
}

func rotateRight(root *node) *node {
	x := root.left
	root.left = x.right
	x.right = root
	x.red = root.red
	root.red = true
	return x
}

func (m *Map) Find(key interface{}) (value interface{}, found bool) {
	root := m.root
	for root != nil {
		if m.less(key, root.key) {
			root = root.left
		} else if m.less(root.key, key) {
			root = root.right
		} else {
			return root.value, true
		}
	}
	return nil, false
}

func (m *Map) Delete(key interface{}) (deleted bool) {
	if m.root != nil {
		if m.root, deleted = m.remove(m.root, key); m.root != nil {
			m.root.red = false
		}
	}
	if deleted {
		m.length--
	}
	return deleted
}

func (m *Map) remove(root *node, key interface{}) (*node, bool) {
    deleted := false
    if m.less(key, root.key) {
        if root.left != nil {
            if !isRed(root.left) && !isRed(root.left.left) {
                root = moveRedLeft(root)
            }
            root.left, deleted = m.remove(root.left, key)
        }
    } else {
        if isRed(root.left) {
            root = rotateRight(root)
        }
        if !m.less(key, root.key) && !m.less(root.key, key) &&
            root.right == nil {
            return nil, true
        }
        if root.right != nil {
            if !isRed(root.right) && !isRed(root.right.left) {
                root = moveRedRight(root)
            }
            if !m.less(key, root.key) && !m.less(root.key, key) {
                smallest := first(root.right)
                root.key = smallest.key
                root.value = smallest.value
                root.right = deleteMinimum(root.right)
                deleted = true
            } else {
                root.right, deleted = m.remove(root.right, key)
            }
        }
    }
    return fixUp(root), deleted
}

func first(root *node) *node {
    for root.left != nil {
        root = root.left
    }
    return root
}

func moveRedLeft(root *node) *node {
    colorFlip(root)
    if root.right != nil && isRed(root.right.left) {
        root.right = rotateRight(root.right)
        root = rotateLeft(root)
        colorFlip(root)
    }
    return root
}

func moveRedRight(root *node) *node {
    colorFlip(root)
    if root.left != nil && isRed(root.left.left) {
        root = rotateRight(root)
        colorFlip(root)
    }
    return root
}

func deleteMinimum(root *node) *node {
    if root.left == nil {
        return nil
    }
    if !isRed(root.left) && !isRed(root.left.left) {
        root = moveRedLeft(root)
    }
    root.left = deleteMinimum(root.left)
    return fixUp(root)
}

func fixUp(root *node) *node {
    if isRed(root.right) {
        root = rotateLeft(root)
    }
    if isRed(root.left) && isRed(root.left.left) {
        root = rotateRight(root)
    }
    if isRed(root.left) && isRed(root.right) {
        colorFlip(root)
    }
    return root
}

func (m *Map) Do(function func(interface{}, interface{})) {
	do(m.root, function)
}

func do(root *node, function func(interface{}, interface{})) {
	if root != nil {
		do(root.left, function)
		function(root.key, root.value)
		do(root.right, function)
	}
}

func (m *Map) Len() int {
	return m.length
}
