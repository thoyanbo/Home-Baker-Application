package main

import (
	"errors"
	"fmt"
)

func searchCategory() error { //this function narrows the search result by displaying items of the category type.
	fmt.Println("Which Category are you searching for?")
	fmt.Println("The categories are 1. Bread, 2. Snacks, 3. Cakes")
	fmt.Println("Please make a selection from 1 to 3.")
	var choice int
	fmt.Scanln(&choice)
	var cat string
	if choice == 1 {
		cat = "Bread"
	} else if choice == 2 {
		cat = "Snacks"
	} else if choice == 3 {
		cat = "Cakes"
	} else {
		return errors.New("you did not make a selection between 1 to 3")
	}
	for i := 0; i < len(items); i++ {
		if cat == items[i].category {
			fmt.Printf("Item: %v     Price: %v:     Category: %v.\n", items[i].name, items[i].unitPrice, items[i].category)
		}
	}
	return nil
}
