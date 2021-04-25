package main

import (
	"errors"
	"fmt"
)

func scanAvailableDay() int { //this function returns an available day booking
	updateWeeklySchedule(orderList)
	var day int
	fmt.Println("Which is the delivery day?")
	fmt.Println("Enter 1 to 7. Eg. 1 for Monday, 2 for Tuesday and etc..")
	fmt.Scanln(&day)
	if day == 0 {
		return 0
	}
	for isDayAvailable(day) == false {
		fmt.Println("You have made an invalid selection, please enter another day.")
		day = scanAvailableDay()
		break
	}
	return day
}

//this function checks if input day is still available for booking
func isDayAvailable(day int) bool {
	for day > 0 && day < 8 {
		if weeklySchedule[day-1].currentOrders < weeklySchedule[day-1].maxOrders {
			return true
		}
		fmt.Println("There is no more available booking for the day.")
		return false

	}
	return false
}

func updateWeeklySchedule(p *orderLL) { //this function iterates over the order linked list and stores the number of daily orders into weekly schedule's current order
	ptr := p.head
	dayCount := make([]int, 7)
	for i := 0; i < p.size; i++ {
		//fmt.Println(ptr.order.deliveryDay)
		dayCount[ptr.order.deliveryDay-1]++
		ptr = ptr.next
	}

	for j := 0; j < len(weeklySchedule); j++ {
		weeklySchedule[j].currentOrders = dayCount[j]
	}
}

func updateWeeklyOrder(p *orderLL) { //this function iterates over the order linked list and tabulates the total count of orders for each sales item
	for j := 0; j < len(items); j++ {
		items[j].weeklyOrder = 0
	} // this code resets count of weekly order to 0 so that we will not use previously counted weeklyOrder
	ptr := p.head
	for i := 0; i < p.size; i++ {
		//fmt.Println(ptr.order.shoppingCart)
		for _, v := range ptr.order.shoppingCart {
			for n := 0; n < len(items); n++ {
				if v.orderItem == items[n].name {
					items[n].weeklyOrder += v.qty
					break
				}
			}
		}
		ptr = ptr.next
	}

}

func (p *orderLL) printDailyOrder(day int) { // this function iterates over the length of the linked list and displays the order if passed in argument day matches order day
	ptr := p.head
	var gotOrder bool = false
	for i := 0; i < p.size; i++ {
		if ptr.order.deliveryDay == day {
			gotOrder = true
			fmt.Println("Name:", ptr.order.name)
			fmt.Println("Address:", ptr.order.address)
			fmt.Println("Order Number:", ptr.order.orderNum)
			fmt.Println("Shopping Cart:", ptr.order.shoppingCart)
			fmt.Println("Amount:", ptr.order.amount)
			fmt.Println("=====================================")
		}
		ptr = ptr.next
	}

	if gotOrder == false {
		fmt.Println("There are no bookings on", intToDay(day))
	}
}

func manageOrder(p *orderLL) error { //manageOrder shows summary of orders. Admin has option to view the orders by day or return to main menu
	userAuthenticated := authenticateAdmin()
	if userAuthenticated == false {
		fmt.Println("Wrong password. Access Denied. returning to main menu")
		return nil
	}
	fmt.Println("This is the weekly summary:")
	fmt.Println("==================================================================================")
	if p.head == nil {
		return errors.New("empty order list")
	}
	var choice int
	var weeklyRevenue float64 = 0

	for i := 0; i < len(items); i++ {
		itemRevenue := items[i].weeklyOrder * items[i].unitPrice
		weeklyRevenue = weeklyRevenue + itemRevenue
	}
	for i := 0; i < len(weeklySchedule); i++ {
		if weeklySchedule[i].currentOrders != 0 {
			fmt.Printf("There are a total of %v orders for %v\n", weeklySchedule[i].currentOrders, intToDay(i+1))
		}
	}
	for i := 0; i < len(items); i++ {
		if items[i].weeklyOrder != 0 {
			fmt.Printf("Quantity %v of item %v has been ordered.\n", items[i].weeklyOrder, items[i].name)
		}
	}

	fmt.Printf("The upcoming week's anticipated revenue is %v\n", weeklyRevenue)
	fmt.Println("Please enter 1 to 7 to view orders by day. Eg: 1 for Monday, 2 for Tuesday..")
	fmt.Println("Or enter '0' to return to main program.")
	fmt.Println("==================================================================================")
	fmt.Scanln(&choice)
	if choice == 0 {
		fmt.Println("Returning to main menu.")
	} else if choice > 0 && choice < 8 {
		fmt.Println("Showing orders for", intToDay(choice))
		orderList.printDailyOrder(choice)
	} else {
		return errors.New("Invalid selection made. Input not between 1 to 7. Returning to main menu")
	}

	return nil
}
