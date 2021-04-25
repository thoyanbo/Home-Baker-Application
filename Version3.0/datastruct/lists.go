//Package datastructs contains the data structs, variables and associated functions used in Home Baker application.
package datastruct

import (
	"errors"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

//ItemInfo stores information about the item
type ItemInfo struct {
	Name           string
	UnitPrice      float64
	Category       string
	Recommended    bool
	WeeklyCapacity float64
	WeeklyOrder    float64
}

var mutex sync.Mutex

var (
	//Items is the global declared slice of sales items and stores relevant information.
	Items = []ItemInfo{
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

//OrderList is set up as a global variable to contain head and size of order linked list
var OrderList = &OrderLL{
	Head: nil,
	Size: 0}

//Order is used to store singular item orders into shopping cart
type Order struct {
	OrderItem string
	Qty       float64
}

//OrderInfo is used to store order information
type OrderInfo struct {
	Username     string
	Name         string
	Address      string
	DeliveryDay  int
	OrderNum     int64
	ShoppingCart []Order
	Amount       float64
}

//OrderNode is the linked list node that stores order data and pointer to next node
type OrderNode struct {
	Order OrderInfo
	Next  *OrderNode
}

//OrderLL stores head and size of order linked list
type OrderLL struct {
	Head *OrderNode
	Size int
}

//PrintAllOrderNodes prints all nodes in order linked list and is used for checking purpose
func (p *OrderLL) PrintAllOrderNodes() error {
	if p.Head == nil {
		log.Error("Empty order list")
		return errors.New("the list is empty")
	}
	currentNode := p.Head
	for currentNode != nil {
		fmt.Println(currentNode.Order)
		currentNode = currentNode.Next
	}
	return nil

}

//bookingDay keeps track of max and current orders status in the days of the week
type bookingDay struct {
	Weekday       int
	MaxOrders     int
	CurrentOrders int
}

//WeeklySchedule is the global variable that tracks the max orders of the day in a week.
var WeeklySchedule = []bookingDay{
	{Weekday: 1, MaxOrders: 8, CurrentOrders: 0},
	{Weekday: 2, MaxOrders: 8, CurrentOrders: 0},
	{Weekday: 3, MaxOrders: 8, CurrentOrders: 0},
	{Weekday: 4, MaxOrders: 8, CurrentOrders: 0},
	{Weekday: 5, MaxOrders: 8, CurrentOrders: 0},
	{Weekday: 6, MaxOrders: 5, CurrentOrders: 0},
	{Weekday: 7, MaxOrders: 5, CurrentOrders: 0},
}

//IntToDay takes in an integer and returns the corresponding day string value of the week.
func IntToDay(day int) string {
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
		return ""
	}
}

//AddOrder allows addition of orders into order linkedlist
func (p *OrderLL) AddOrder(oi OrderInfo) {
	mutex.Lock()
	{
		newNode := &OrderNode{
			Order: oi,
			Next:  nil,
		}
		if p.Head == nil {
			p.Head = newNode
		} else {
			currentNode := p.Head
			for currentNode.Next != nil {
				currentNode = currentNode.Next
			}
			currentNode.Next = newNode
		}
		p.Size++
	}
	mutex.Unlock()
	return
}

//CalculateAmount calculates total amount based on shopping cart array
func CalculateAmount(sc []Order) float64 {
	var totalAmount float64
	var itemAmount float64 = 0
	for n := 0; n < len(sc); n++ {
		itemAmount = 0
		for i := 0; i < len(Items); i++ {
			if sc[n].OrderItem == Items[i].Name {
				itemAmount = (Items[i].UnitPrice * sc[n].Qty)
				totalAmount = totalAmount + itemAmount
			}
		}
	}
	return totalAmount
}

//IsDayAvailable returns a boolean value if an input day is available for taking orders.
func IsDayAvailable(day int) bool {
	for day > 0 && day < 8 {
		if WeeklySchedule[day-1].CurrentOrders < WeeklySchedule[day-1].MaxOrders {
			return true
		}
		return false

	}
	return false
}
