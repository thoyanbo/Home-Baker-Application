package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

var orderNumber int64 = 0 // global variable to track order number

func addShoppingCart() order { //function to add type order data into order
	var name string
	var qty float64
	var toAdd order
	fmt.Println("Choose item to order.")
	for n := 0; n < len(items); n++ {
		fmt.Printf("Item: %v     Price: %v:    \n", items[n].name, items[n].unitPrice)
	}
	name = orderExistingInStockItem()
	qty = getOrderQuantity(name)
	toAdd = order{
		orderItem: name,
		qty:       qty,
	}
	return toAdd
}

func existingItem(nam string) (exist bool, err error) { //function to check if user enters an existing item from global variable items. case sensitive.
	for n := 0; n < len(items); n++ {
		if nam == items[n].name {
			return true, nil
		}
	}
	fmt.Println("Item does not exist! Please enter an existing")
	return false, nil
}

func availableItem(nam string) (exist bool, err error) { //function to check if input item still has weekly stock for ordering
	for n := 0; n < len(items); n++ {
		if nam == items[n].name {
			if items[n].weeklyCapacity > items[n].weeklyOrder {
				return true, nil
			}
			fmt.Println("Item is not in stock! Please make another selection")
			return false, nil
		}
	}
	return false, nil
}

func orderExistingInStockItem() string { // this function is recursive in nature and will only terminate when user enters a valid item selection that is case sensitive
	var exist bool = false
	var available bool = false
	var nam string
	for exist != true && available != true {
		fmt.Println("Please type name of item to order.")
		reader := bufio.NewReader(os.Stdin) //create new reader, assuming bufio imported
		nam, _ = reader.ReadString('\n')
		nam = strings.TrimRight(nam, "\n") // this removes the \n at end of scan function
		exist, _ = existingItem(nam)
		available, _ = availableItem(nam)
	}
	return nam
}

func getOrderQuantity(nam string) float64 { //this function iterates through the items list to obtain the balance quantity available through the week

	var qty float64
	var balanceQty float64
	for n := 0; n < len(items); n++ {
		if nam == items[n].name {
			balanceQty = items[n].weeklyCapacity - items[n].weeklyOrder
			fmt.Printf("Quantity %v of %v is still in stock. Please enter quantity to order.\n", balanceQty, nam)
			break
		}

	}

	for qty < balanceQty {
		fmt.Printf("How many %v would you like? Please enter quantity\n", nam)
		fmt.Scanln(&qty)
		if qty > balanceQty {
			fmt.Println("You have ordered too much, please enter an amount smaller than", balanceQty)
			getOrderQuantity(nam)
		}
		break
	}
	return qty
}

func doneShopping() bool { //function to track addition of ordering. will return false if done.
	var response string
	fmt.Println("Do you still want to add items to your order?")
	fmt.Println("Enter Y for Yes, N for No.")
	fmt.Scanln(&response)
	if response == "Y" {
		return false
	} else if response == "N" {
		return true
	} else {
		fmt.Println("You have not selected Y or N, please try again.")
		responseBool := doneShopping()
		return responseBool
	}
}

func generateShoppingCart(shoppingOrder chan []order) { //generates a shopping cart array
	var shoppingComplete bool = false
	var newShoppingCart = []order{} //creates an empty order
	for shoppingComplete != true {
		addShopping := addShoppingCart()
		newShoppingCart = append(newShoppingCart, addShopping)
		shoppingComplete = doneShopping()
	}
	shoppingOrder <- newShoppingCart
}

func calculateAmount(sc []order) float64 { //calculates amount based on shopping cart array
	var totalAmount float64
	var itemAmount float64 = 0
	for n := 0; n < len(sc); n++ {
		itemAmount = 0
		for i := 0; i < len(items); i++ {
			if sc[n].orderItem == items[i].name {
				itemAmount = (items[i].unitPrice * sc[n].qty)
				totalAmount = totalAmount + itemAmount
			}
		}
	}
	return totalAmount
}

func createOrder() error { //function to create a new order which will be added into order linked list
	defer func() { //to handle possible panic situation
		if err := recover(); err != nil {
			fmt.Println("Oops, panic occurred:", err)
		}
	}()
	var nam string
	var add string
	var dday int
	var on int64
	var amt float64
	fmt.Println("Enter Name:")
	reader := bufio.NewReader(os.Stdin) //create new reader, assuming bufio imported
	nam, _ = reader.ReadString('\n')
	nam = strings.TrimRight(nam, "\n") // this removes the \n at end of scan function
	if nam == "" {
		return errors.New("invalid empty input")
	}
	fmt.Println("Enter Address:")
	add, _ = reader.ReadString('\n')
	add = strings.TrimRight(add, "\n") // this removes the \n at end of scan function
	if add == "" {
		return errors.New("invalid empty input")
	}
	dday = scanAvailableDay()
	shoppingOrder := make(chan []order)
	go generateShoppingCart(shoppingOrder)
	sc := <-shoppingOrder
	on = atomic.AddInt64(&orderNumber, 1)

	amt = calculateAmount(sc)

	newOrder := orderInfo{
		name:         nam,
		address:      add,
		deliveryDay:  dday,
		orderNum:     on,
		shoppingCart: sc,
		amount:       amt,
	}

	var wg sync.WaitGroup

	wg.Add(1)

	func() {
		defer wg.Done()
		orderList.addOrder(newOrder)
	}()

	wg.Wait()

	//fmt.Println(newOrder, orderNumber) // this is used to check for order entry
	//orderList.printAllOrderNodes() //this is used to check for all orders
	var wg2 sync.WaitGroup
	func() {
		wg2.Add(2)
		go func() {
			defer wg2.Done()
			updateWeeklySchedule(orderList)
		}()
		go func() {
			defer wg2.Done()
			updateWeeklyOrder(orderList)
		}()
		wg2.Wait()
		//fmt.Println(weeklySchedule) // this is to check if weeklySchedule is updated
	}()
	fmt.Printf("Your order has been successfully created! Your order number is %v, please remember this number for future viewing or editing of order.\n", on)
	return nil
}
