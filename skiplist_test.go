package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"testing"

	"github.com/go-test/deep"
)

type insertionArgs struct {
	key   int
	value any
}

func runSkipListTest_StaticInsertSearchInt(t *testing.T, iteration int) {
	list := NewSkipList[int](10) // we take max level out as a factor by making it sufficiently large

	values := []insertionArgs{
		{5, "five"},
		{10, "ten"},
		{15, "fifteen"},
		{0, "zero"},
		{-3, "minus three"},
		{7, "seven"},
		{14, "fourteen"},
	}

	// build representation of entire data structure
	// Test insert
	for _, v := range values {
		list.Insert(v.key, v.value)
	}

	// Confirm that insertion order is correct by iterating through
	// base level (L0)
	baseLevel := 0
	expectedOrdering := []int{-3, 0, 5, 7, 10, 14, 15}
	actualOrdering := []int{}
	for elem := list.head.forward[baseLevel]; elem != nil; elem = elem.next(baseLevel) {
		actualOrdering = append(actualOrdering, elem.key)
	}
	if diff := deep.Equal(expectedOrdering, actualOrdering); diff != nil {
		t.Logf("Iteration %d:\n%s", iteration, list)
		t.Errorf("Elements not being correct inserted to list by ascending key value: %v", diff)
	}

	// Test search for all
	for _, v := range values {
		if got := list.Search(v.key); got != v.value {
			t.Logf("Iteration %d:\n%s", iteration, list)
			t.Errorf("Search(%d) = %v, want %v", v.key, got, v.value)
		}
	}

	// Test search for non-existent key
	if got := list.Search(20); got != nil {
		t.Errorf("Search(20) = %v, want nil", got)
	}
}

func TestSkipList_InsertSearchInt(t *testing.T) {
	tests := []struct {
		name       string
		f          func(*testing.T, int)
		iterations int
	}{
		{"defined set {-3,0,5,7,10,14,15}", runSkipListTest_StaticInsertSearchInt, 1000},
		{"randomly generated set of 2^16 values", runSkipListTest_RandomInsertSearchInt, 10},
	}

	for _, tt := range tests {
		for i := 0; i < tt.iterations; i++ {
			t.Run(fmt.Sprintf("%s_%04d", tt.name, i), func(t *testing.T) {
				tt.f(t, i)
			})
		}
	}
}

func runSkipListTest_RandomInsertSearchInt(t *testing.T, iteration int) {
	maxLevelHeight := uint(10)
	list := NewSkipList[int](maxLevelHeight) // Sufficiently large max level

	// Generate 2^16 random entries to the skip list
	scale := 65536
	randomInts := generateRandomInts(scale, -1000000, 1000000)
	var values []insertionArgs
	for _, key := range randomInts {
		values = append(values, insertionArgs{key: key, value: strconv.Itoa(key)})
	}

	// Insert values into the skip list
	for _, v := range values {
		list.Insert(v.key, v.value)
	}

	// check that we haven't exceeded headroom
	if list.level == list.maxLevel {
		t.Errorf("hit maxLevel (height) with %d elements", scale)
	}

	// Generate expected ordering by sorting the keys
	sort.Ints(randomInts) // Sort the keys to get expected ordering
	expectedOrdering := randomInts

	// Check base level ordering
	baseLevel := 0
	actualOrdering := []int{}
	for elem := list.head.forward[baseLevel]; elem != nil; elem = elem.next(baseLevel) {
		actualOrdering = append(actualOrdering, elem.key)
	}
	if diff := deep.Equal(expectedOrdering, actualOrdering); diff != nil {
		t.Logf("Iteration %d:\n%v\n%v\n", iteration, expectedOrdering, expectedOrdering)
		t.Errorf("Iteration: %d\t Incorrect insertion order: %v", iteration, diff)
	}

	// Check search functionality
	for _, v := range values {
		if got := list.Search(v.key); got != v.value {
			t.Logf("Iteration %d:\n%s", iteration, list)
			t.Errorf("Search(%d) = %v, want %v", v.key, got, v.value)
		}
	}
}

func generateRandomInts(n int, min int, max int) []int {
	ints := make(map[int]bool)
	var nums []int

	for len(nums) < n {
		num := rand.Intn(max-min+1) + min
		if !ints[num] {
			ints[num] = true
			nums = append(nums, num)
		}
	}

	return nums
}

func TestSkipList_UpdateInPlace(t *testing.T) {
	list := NewSkipList[int](10) // we take max level out as a factor by making it sufficiently large

	type insertionArgs struct {
		key   int
		value any
	}
	original_5 := insertionArgs{5, "five"}
	updated_5 := insertionArgs{5, "FIVE"}
	values := []insertionArgs{
		original_5,
		{10, "ten"},
		{15, "fifteen"},
	}

	// Test insert
	for _, v := range values {
		list.Insert(v.key, v.value)
	}

	// Test search for key (5) with expected value ("five")
	if got := list.Search(5); got != original_5.value {
		t.Errorf("Search(5) = %v, want %v", got, original_5.value)
	}

	// Update value associated with key (5)
	list.Insert(updated_5.key, updated_5.value)
	// Test search for key (5) with expected updated value ("FIVE")
	if got := list.Search(5); got != updated_5.value {
		t.Errorf("Search(5) = %v, want %v", got, updated_5.value)
	}
}

func TestSkipList_DeleteInt(t *testing.T) {
	list := NewSkipList[int](10)

	// Insert values
	list.Insert(5, "five")
	list.Insert(10, "ten")
	list.Insert(15, "fifteen")

	// Delete a value and test
	list.Delete(10)
	if got := list.Search(10); got != nil {
		t.Errorf("Search(10) after Delete(10) = %v, want nil", got)
	}

	// Test deletion of non-existing key
	list.Delete(20)
	if got := list.Search(15); got != "fifteen" {
		t.Errorf("Search(15) after Delete(20) = %v, want 'fifteen'", got)
	}
}

// TODO: Additional tests can be added to cover more types, edge cases, and complex scenarios.
