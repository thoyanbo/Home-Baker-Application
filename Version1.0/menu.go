package main

import (
	"fmt"
)

var orderList = &orderLL{ //orderList is set up as a global variable to contain head and size of order linked list
	head: nil,
	size: 0}

//printMenu shows the menu page
func printMenu() {
	var choice int
	for choice != 10 {
		fmt.Println("Welcome to Home Baker App Beta.")
		fmt.Println("===============================")
		fmt.Println("Menu Selection")
		fmt.Println("===============================")
		fmt.Println("1. Display sales items")
		fmt.Println("2. Search")
		fmt.Println("3. Create order")
		fmt.Println("4. View/Edit order")
		fmt.Println("5. Delete order (Admin Feature)")
		fmt.Println("6. Weekly Overview(Admin Feature)")
		fmt.Println("10. Exit program")
		fmt.Println("Please make your selection:")
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			err := displayItems(items)
			if err != nil {
				fmt.Println(err)
			}
		case 2:
			err := search()
			if err != nil {
				fmt.Println(err)
			}
		case 3:
			err := createOrder()
			if err != nil {
				fmt.Println(err)
			}
		case 4:
			err := viewOrEditMenu()
			if err != nil {
				fmt.Println(err)
			}
		case 5:
			err := orderList.deleteOrder()
			if err != nil {
				fmt.Println(err)
			}
		case 6:
			err := manageOrder(orderList)
			if err != nil {
				fmt.Println(err)
			}
			/*} else if choice == 6 {
			addNewSalesItem()*/
		case 10:
			fmt.Println("Exiting program..")
		default:
			fmt.Println("You have not selected made a valid selection, please try again.")
		}
	}
}
