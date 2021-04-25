package main

import (
	"fmt"
)

func search() error {
	var choice int
	fmt.Println("The following options are available:")
	fmt.Println("1. Search by name.")
	fmt.Println("2. View by price.")
	fmt.Println("3. View by category.")
	fmt.Println("4. Recommended.")
	fmt.Println("5. See available delivery days")
	fmt.Println("Please make a selection from 1 to 5.")
	fmt.Scanln(&choice)

	if choice == 1 {
		err := searchName()
		if err != nil {
			fmt.Println(err)
		}
	} else if choice == 2 {
		sortPrice(items)
	} else if choice == 3 {
		err := searchCategory()
		if err != nil {
			fmt.Println(err)
		}
	} else if choice == 4 {
		showRecommended()
	} else if choice == 5 {
		SeeAvailableDays()
	} else {
		fmt.Println("You did not make a valid selection.")
	}
	return nil
}
