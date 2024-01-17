package main

import (
	"cmp"
	"math"
	"math/rand"
)

type Node[T cmp.Ordered] struct {
	value   T
	forward []*Node[T]
}

type SkipList[T cmp.Ordered] struct {
	head     *Node[T]
	level    uint
	maxLevel uint
}

func NewSkipList[T cmp.Ordered](maxElements uint) *SkipList[T] {
	maxLevel := uint(math.Ceil(float64(maxElements)))
	head := &Node[T]{
		forward: make([]*Node[T], maxLevel),
	}
	list := SkipList[T]{
		head:     head,
		maxLevel: maxLevel,
		level:    0, // base level is the only level to begin with
	}
	return &list
}

func (sl *SkipList[T]) randomLevel() uint {
	level := uint(0) // start at base level
	for rand.Intn(2) < 1 && level < sl.maxLevel {
		level++
	}
	return level
}

func (sl *SkipList[T]) Insert(value T) {
	update := make([]*Node[T], sl.maxLevel)
	current := sl.head

	// Find insertion position
	for i := sl.level; i >= 0; i-- {
		for current.forward[i] != nil {
		}
	}
}
