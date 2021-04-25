package main

import (
	"fmt"
)

func displayItems(ii []itemInfo) error { //Note that display items by default will sort the items by category and display by category.
	defer func() { //to handle possible panic situation
		if err := recover(); err != nil {
			fmt.Println("Oops, panic occurred:", err)
		}
	}()
	var n int
	n = len(items)
	fmt.Println("These are our items for sale!")
	fmt.Println("==================================")

	for last := n - 1; last >= 1; last-- {
		// select the last alphabetical item in array
		largest := indexOfLast(ii, last+1)

		//swap largest item array[largest] with array[last]
		swap(&ii[largest], &ii[last])
	}
	for n := 0; n < len(items); n++ {
		fmt.Printf("Item: %v     Price: %v:     Category: %v.\n", items[n].name, items[n].unitPrice, items[n].category)
	}
	return nil
}

func indexOfLast(ii []itemInfo, n int) int { // this function returns index of last item in alphabetical order through string comparison
	lastIndex := 0
	for i := 1; i < n; i++ {
		if ii[i].name > ii[lastIndex].name {
			lastIndex = i
		}
	}
	return lastIndex
}
