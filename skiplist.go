package main

import (
	"cmp"
	"fmt"
	"math/rand"
	"strings"
)

type forwarder[T cmp.Ordered] interface {
	next(level int) *Node[T]
	setNext(level int, node *Node[T])
}

type Node[T cmp.Ordered] struct {
	key     T
	forward []*Node[T]
	value   any
}

func (n *Node[T]) next(level int) *Node[T] {
	if level < len(n.forward) {
		return n.forward[level]
	}
	return nil
}

func (n *Node[T]) setNext(level int, node *Node[T]) {
	if level < len(n.forward) {
		n.forward[level] = node
	}
}

func (n *Node[T]) String() string {
	return fmt.Sprintf("(%v)", n.key)
}

type Header[T cmp.Ordered] struct {
	forward []*Node[T]
}

func (h *Header[T]) next(level int) *Node[T] {
	if level < len(h.forward) {
		return h.forward[level]
	}
	return nil
}

func (h *Header[T]) setNext(level int, node *Node[T]) {
	if level < len(h.forward) {
		h.forward[level] = node
	}
}

type SkipList[T cmp.Ordered] struct {
	head  *Header[T]
	level int
	// TODO: make p specifiable
	p        float64 // Probability, default is P=(1/2)
	maxLevel int
	seed     int64
}

func NewSkipList[T cmp.Ordered](maxLevel uint) *SkipList[T] {
	head := &Header[T]{
		forward: make([]*Node[T], maxLevel),
	}
	list := SkipList[T]{
		head:     head,
		level:    0,   // L0 is the only list initially
		p:        0.5, // P=(1/2)
		maxLevel: int(maxLevel),
	}
	return &list
}

func (sl *SkipList[T]) randomLevel() int {
	// if list is empty, use 0
	if sl.head.forward[0] == nil {
		return 0
	}
	level := 0 // start at base level
	for level < sl.level+1 && level < sl.maxLevel-1 {
		if rand.Float64() < sl.p {
			break // tails, stop
		}
		level++ // heads, increment level
	}
	return level
}

func (sl *SkipList[T]) Insert(key T, value any) {
	// During the traversal from the head at the current SkipList.level
	// to find the insertion position on the base level (L0), we must keep track
	// of the node at every stage where traversal requires 'stepping down' (e.g. L3->L2)
	// to know where to updateList pointers at each level (up to the random insertion level)
	updateList := make([]forwarder[T], sl.maxLevel)

	// STEP 1: Find insertion position
	var current forwarder[T]
	current = sl.head

	// Iterate down through levels (vertical movement)
	for lvl := sl.level; lvl >= 0; lvl-- {
		// Iterate through nodes at a given level (horizontal movement)
		for current.next(lvl) != nil && current.next(lvl).key < key {
			current = current.next(lvl)
		}
		// Store which nodes to update for each level
		updateList[lvl] = current
	}

	// STEP 1B: If existing key found, update value in-place
	if target, ok := current.(*Node[T]); ok {
		target = target.next(0)
		if target != nil && target.key == key {
			target.value = value
			return
		}
	}

	// STEP 2: Determine insertion level by random variable
	insertionLevel := sl.randomLevel()

	// Check if new level is greater than skip list's current level
	// promoteNewNode := false
	// if insertionLevel > sl.level {
	// 	firstElementAtCurrentLevel := sl.head.next(sl.level)
	// 	promoteNewNode = firstElementAtCurrentLevel == nil || key < firstElementAtCurrentLevel.key
	// }

	// If new level is greater than skip list's current level, modify update list
	if insertionLevel > sl.level {
		for i := sl.level + 1; i <= insertionLevel; i++ {
			updateList[i] = sl.head
		}
		sl.level = insertionLevel
	}

	// STEP 3: Create new node
	newNode := &Node[T]{
		key:     key,
		value:   value,
		forward: make([]*Node[T], insertionLevel+1),
	}

	// STEP 4: Insert new node into list by updating pointers for each level of insertion
	for i := 0; i <= insertionLevel && i <= sl.level; i++ {
		newNode.setNext(i, updateList[i].next(i))
		updateList[i].setNext(i, newNode)
	}

	// STEP 5: Update skip list level if new node level >
	if insertionLevel > sl.level {
		sl.level = insertionLevel
	}
}

func (sl *SkipList[T]) Delete(key T) {
	updateList := make([]forwarder[T], sl.maxLevel)

	// STEP 1: Find insertion position
	var current forwarder[T]
	current = sl.head

	// STEP 1: Find position of target node and store the previous nodes at each level
	for lvl := sl.level; lvl >= 0; lvl-- {
		for current.next(lvl) != nil && current.next(lvl).key < key {
			current = current.next(lvl)
		}
		updateList[lvl] = current
	}

	// Check if the target node exists
	target := current.next(0)
	if target != nil && target.key == key {
		// Update next pointers only up to the level of the target node
		for lvl := 0; lvl < len(target.forward); lvl++ {
			updateList[lvl].setNext(lvl, target.next(lvl))
		}

		// Remove levels that have no elements by adjusting the skip list level if empty
		for sl.level > 0 && sl.head.forward[sl.level] == nil {
			sl.level--
		}
	}
}

func (sl *SkipList[T]) Search(target T) any {
	defer func() {
		fmt.Println("end Search for:", target)
		fmt.Println("---------------------")
	}()
	var current forwarder[T]

	current = sl.head
	fmt.Println("beginning Search for:", target)
	fmt.Println("debugging print of structure:")
	fmt.Println(sl.String())
	for lvl := sl.level; lvl >= 0; lvl-- {
		fmt.Printf("\nL%v\t", lvl)
		for current.next(lvl) != nil && current.next(lvl).key <= target {
			current = current.next(lvl)
			node := current.(*Node[T])
			fmt.Printf("->(%v)\t", node.key)
			if node.key == target {
				fmt.Println("successfully found:", target)
				return node.value
			}
		}
	}

	fmt.Println("did not find:", target)
	return nil
}

func (sl *SkipList[T]) String() string {
	var levelStrings []string

	// Build the string for the base level (L0)
	var baseLevelBuilder strings.Builder
	current := sl.head.next(0)
	for current != nil {
		baseLevelBuilder.WriteString(fmt.Sprintf("(%v) -> ", current.key))
		current = current.next(0)
	}
	baseLevelStr := strings.TrimSuffix(baseLevelBuilder.String(), " -> ")
	levelStrings = append(levelStrings, baseLevelStr)

	// Construct strings for the higher levels based on the length of the base level string
	baseLen := len(baseLevelStr)
	for i := 1; i <= sl.level; i++ {
		var levelBuilder strings.Builder
		current := sl.head.next(i)

		pos := 0 // Starting position
		for current != nil {
			// Calculate the position where the next node should appear
			nextPos := strings.Index(baseLevelStr[pos:], fmt.Sprintf("(%v)", current.key)) + pos

			// Fill with dashes up to the next node, leaving space for "> "
			if nextPos > pos {
				levelBuilder.WriteString(strings.Repeat("-", nextPos-pos-2))
				levelBuilder.WriteString("> ")
			}

			// Append the node and a space after it
			nodeWithTrailingSpace := current.String() + " "
			levelBuilder.WriteString(nodeWithTrailingSpace)
			pos = nextPos + len(nodeWithTrailingSpace)

			// Move to the next node
			current = current.next(i)
		}

		// Fill the rest of the string with spaces to match the base level's length
		if pos < baseLen {
			levelBuilder.WriteString(strings.Repeat(" ", baseLen-pos))
		}

		levelStrings = append(levelStrings, levelBuilder.String())
	}

	// Reverse the slice and concatenate the strings for each level
	var builder strings.Builder
	for i := sl.level; i >= 0; i-- {
		builder.WriteString(fmt.Sprintf("L%d: %s\n", i, levelStrings[i]))
	}

	return builder.String()
}
