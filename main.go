package main

import "fmt"

func main() {
	sl := NewSkipList[int](3)

	// Insert values
	sl.Insert(3, "three")
	sl.Insert(6, "six")
	sl.Insert(7, "seven")
	sl.Insert(9, "nine")
	sl.Insert(12, "twelve")
	sl.Insert(7, "APPLE")

	// Search for a value

	value := sl.Search(9)
	fmt.Printf("Search for 9: %v, %+v\n", value != nil, value)
	value = sl.Search(7)
	fmt.Printf("Search for 7: %v, %+v\n", value != nil, value)

	// Delete a value
	sl.Delete(6)

	// Try searching for deleted value
	value = sl.Search(6)
	fmt.Printf("Search for 6: %v, %+v\n", value != nil, value)
	fmt.Printf("%#v", sl)
	value = sl.Search(0)
	fmt.Printf("Search for 0: %v, %+v\n", value != nil, value)
}
