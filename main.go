package main

import "fmt"

func main() {
	sl := NewSkipList[int](3)

	// Insert values
	sl.Insert(3)
	sl.Insert(6)
	sl.Insert(7)
	sl.Insert(9)
	sl.Insert(12)

	// Search for a value
	fmt.Println("Search for 9:", sl.Search(9) != nil)

	// Delete a value
	sl.Delete(6)

	// Try searching for deleted value
	fmt.Println("Search for 6:", sl.Search(6) != nil)
}
