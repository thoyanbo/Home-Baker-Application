package main

import (
	"errors"
	"fmt"
)

//itemInfo stores information about the item
type itemInfo struct {
	Name           string
	UnitPrice      float64
	Category       string
	Recommended    bool
	WeeklyCapacity float64
	WeeklyOrder    float64
}

//items is the global declared slice of sales items and stores relevant information.
var (
	items = []itemInfo{
		{Name: "Pineapple Tarts", UnitPrice: 18.00, Category: "Snacks", Recommended: true, WeeklyCapacity: 35, WeeklyOrder: 0},
		{Name: "Coffee Cookies", UnitPrice: 12.00, Category: "Snacks", Recommended: false, WeeklyCapacity: 35, WeeklyOrder: 0},
		{Name: "Kaya Buns", UnitPrice: 5.00, Category: "Bread", Recommended: true, WeeklyCapacity: 70, WeeklyOrder: 0},
		{Name: "Mochi Buns", UnitPrice: 6.00, Category: "Bread", Recommended: true, WeeklyCapacity: 70, WeeklyOrder: 0},
		{Name: "Walnut Raisin Buns", UnitPrice: 5.00, Category: "Bread", Recommended: false, WeeklyCapacity: 70, WeeklyOrder: 0},
		{Name: "4 Cheese Rolls", UnitPrice: 7.50, Category: "Bread", Recommended: false, WeeklyCapacity: 70, WeeklyOrder: 0},
		{Name: "Blackforest Cake", UnitPrice: 10.00, Category: "Pastries", Recommended: true, WeeklyCapacity: 35, WeeklyOrder: 0},
		{Name: "Burnt Cheesecake", UnitPrice: 12.00, Category: "Pastries", Recommended: false, WeeklyCapacity: 35, WeeklyOrder: 0},
		{Name: "Red Velvet Carrot Cake", UnitPrice: 8.00, Category: "Pastries", Recommended: true, WeeklyCapacity: 35, WeeklyOrder: 0},
		{Name: "Baileys Chocolate Cake", UnitPrice: 8.00, Category: "Pastries", Recommended: false, WeeklyCapacity: 35, WeeklyOrder: 0},
	}
)

var orderList = &orderLL{ //orderList is set up as a global variable to contain head and size of order linked list
	head: nil,
	size: 0}

//order is used to store singular item orders into shopping cart
type order struct {
	orderItem string
	qty       float64
}

//orderInfo is used to store order information
type orderInfo struct {
	Username     string
	Name         string
	Address      string
	DeliveryDay  int
	OrderNum     int64
	ShoppingCart []order
	Amount       float64
}

//OrderNode is the linked list node that stores order data and pointer to next node
type orderNode struct {
	order orderInfo
	next  *orderNode
}

//OrderLL stores head and size of order linked list
type orderLL struct {
	head *orderNode
	size int
}

//printAllOrderNodes prints all nodes in order linked list and is used for checking purpose
func (p *orderLL) printAllOrderNodes() error {
	if p.head == nil {
		return errors.New("the list is empty")
	}
	currentNode := p.head
	for currentNode != nil {
		fmt.Println(currentNode.order)
		currentNode = currentNode.next
	}
	return nil

}

//bookingDay keeps track of max and current orders status in the days of the week
type bookingDay struct {
	weekday       int
	maxOrders     int
	currentOrders int
}

//weeklySchedule is the global variable that tracks the max orders of the day in a week.
var weeklySchedule = []bookingDay{
	{weekday: 1, maxOrders: 8, currentOrders: 0},
	{weekday: 2, maxOrders: 8, currentOrders: 0},
	{weekday: 3, maxOrders: 8, currentOrders: 0},
	{weekday: 4, maxOrders: 8, currentOrders: 0},
	{weekday: 5, maxOrders: 8, currentOrders: 0},
	{weekday: 6, maxOrders: 5, currentOrders: 0},
	{weekday: 7, maxOrders: 5, currentOrders: 0},
}

//IntToDay takes in an integer and returns the corresponding day string value of the week.
func intToDay(day int) string {
	if day == 1 {
		return "Monday"
	} else if day == 2 {
		return "Tuesday"
	} else if day == 3 {
		return "Wednesday"
	} else if day == 4 {
		return "Thursday"
	} else if day == 5 {
		return "Friday"
	} else if day == 6 {
		return "Saturday"
	} else if day == 7 {
		return "Sunday"
	} else {
		return "Error! You did not select a weekday."
	}
}
