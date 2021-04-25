package main

import (
	"fmt"
)

func showRecommended() { //this function shows the recommended items
	fmt.Println("The following items are our best-sellers!")
	for i := 0; i < len(items); i++ {
		if items[i].recommended == true {
			fmt.Printf("Item: %v     Price: %v:     Category: %v.\n", items[i].name, items[i].unitPrice, items[i].category)
		}
	}

}
