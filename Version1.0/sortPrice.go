package main

import (
	"fmt"
)

func sortPrice(ii []itemInfo) { // this function does a selection sort of the array by price.
	var n int
	n = len(items)

	for last := n - 1; last >= 1; last-- {
		// select most expensive item in array
		largest := indexOfLargest(ii, last+1)

		//swap largest item array[largest] with array[last]
		swap(&ii[largest], &ii[last])
	}

	for n := 0; n < len(items); n++ {
		fmt.Printf("Price: %v:      Item: %v    Category: %v.\n", items[n].unitPrice, items[n].name, items[n].category)
	}
}

func indexOfLargest(ii []itemInfo, n int) int {
	largestIndex := 0
	for i := 1; i < n; i++ {
		if ii[i].unitPrice > ii[largestIndex].unitPrice {
			largestIndex = i
		}
	}
	return largestIndex
}

func swap(x *itemInfo, y *itemInfo) {
	temp := *x
	*x = *y
	*y = temp
}
