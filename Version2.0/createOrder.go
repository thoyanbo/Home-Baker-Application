package main

import (
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
)

var orderNumber int64 = 0 // global variable to track order number

func createNewOrder(res http.ResponseWriter, req *http.Request) {

	myUser := getUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	//fmt.Println(items)
	sortItems(items)

	var newShoppingCart = []order{}

	if req.Method == http.MethodPost {
		name := req.FormValue("name")
		if name == "" {
			http.Error(res, "You cannot leave name field empty.", http.StatusBadRequest)
			return
		}
		add := req.FormValue("address")
		if add == "" {
			http.Error(res, "You cannot leave address field empty.", http.StatusBadRequest)
			return
		}
		dday, _ := strconv.Atoi(req.FormValue("dday"))

		availableDay := isDayAvailable(dday)
		if availableDay == false { //this checks if the order was placed on an unavailable day
			http.Error(res, "There are no more available delivery slots for your selected booking day", http.StatusBadRequest)
			return
		}

		for i := 0; i < len(items); i++ {
			quantity, _ := strconv.Atoi(req.FormValue(items[i].Name))
			quantity64 := float64(quantity)

			if quantity64 > 0 {
				itemAvailable := availableItem(items[i].Name)
				if itemAvailable == false { // this checks if the current item is in stock
					http.Error(res, "One of the items ordered is not available", http.StatusBadRequest)
					return
				}
				availableBalance := isBalanceEnough(items[i].Name, quantity64)
				if availableBalance == false { //this checks if the user over ordered on the item
					http.Error(res, "There is not sufficient quantity for one of the items ordered", http.StatusBadRequest)
					return
				}
				singleCart := order{
					orderItem: items[i].Name,
					qty:       quantity64,
				}
				newShoppingCart = append(newShoppingCart, singleCart)
			}
		}

		if len(newShoppingCart) == 0 {
			http.Error(res, "Error: You cannot submit an empty shopping cart.", http.StatusBadRequest)
			return
		}

		on := atomic.AddInt64(&orderNumber, 1) // use of atomic function to prevent multiple clients from possibly creating identical order number
		amt := calculateAmount(newShoppingCart)
		newOrder := orderInfo{
			Username:     myUser.Username,
			Name:         name,
			Address:      add,
			DeliveryDay:  dday,
			OrderNum:     on,
			ShoppingCart: newShoppingCart,
			Amount:       amt,
		}

		orderList.addOrder(newOrder)
		err := updateWeeklySchedule(orderList)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		err = updateWeeklyOrder(orderList)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		//fmt.Println(weeklySchedule)
		//fmt.Println(items)
		//orderList.printAllOrderNodes()

		http.Redirect(res, req, "/menu", http.StatusSeeOther)
		return
	}

	type balanceStruct struct {
		Item     string
		Quantity float64
	}

	var itemData []balanceStruct

	for i := 0; i < len(items); i++ {
		remainingQuantity := items[i].WeeklyCapacity - items[i].WeeklyOrder
		d := balanceStruct{
			Item:     items[i].Name,
			Quantity: remainingQuantity}
		itemData = append(itemData, d)
	}

	err := tpl.ExecuteTemplate(res, "createOrder.gohtml", itemData)
	if err != nil {
		log.Fatalln(err)
	}

	viewAvailableDays(res, req)
	showRemainingBalance(res, req)
}

func viewAvailableDays(res http.ResponseWriter, req *http.Request) { //function to show remaining available days to user during create/edit order page
	type dayData struct {
		Day          string
		RemainingOrd int
	}
	var AvailableDaysData []dayData

	for n := 0; n < len(weeklySchedule); n++ {
		if weeklySchedule[n].currentOrders < weeklySchedule[n].maxOrders {
			remainingOrders := weeklySchedule[n].maxOrders - weeklySchedule[n].currentOrders
			if remainingOrders > 0 {
				tempData := dayData{
					Day:          intToDay(n + 1),
					RemainingOrd: remainingOrders}
				AvailableDaysData = append(AvailableDaysData, tempData)
			}
			//fmt.Println(intToDay(n+1), "is available for booking.", remainingOrders, "orders still available for taking.")
		}
	}
	err := tpl.ExecuteTemplate(res, "availableDays.gohtml", AvailableDaysData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

func showRemainingBalance(res http.ResponseWriter, req *http.Request) { //function to show remaining balance to user during create/edit order page
	type balanceStruct struct {
		Item     string
		Quantity float64
	}

	var itemData []balanceStruct

	for i := 0; i < len(items); i++ {
		remainingQuantity := items[i].WeeklyCapacity - items[i].WeeklyOrder
		d := balanceStruct{
			Item:     items[i].Name,
			Quantity: remainingQuantity}
		itemData = append(itemData, d)
	}

	err := tpl.ExecuteTemplate(res, "itemBalance.gohtml", itemData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

func availableItem(nam string) (exist bool) { //function to check if input item still has weekly stock for ordering
	for n := 0; n < len(items); n++ {
		if nam == items[n].Name {
			if items[n].WeeklyCapacity > items[n].WeeklyOrder {
				return true
			}
			return false
		}
	}
	return false
}

func isBalanceEnough(nam string, qty float64) bool { //this function checks if the order quantity exceeds the balance quantity, and returns false if true

	var balanceQty float64
	for n := 0; n < len(items); n++ {
		if nam == items[n].Name {
			balanceQty = items[n].WeeklyCapacity - items[n].WeeklyOrder
			break
		}
	}

	if qty <= balanceQty {
		return true
	}
	return false
}

func calculateAmount(sc []order) float64 { //calculates Amount based on shopping cart array
	var totalAmount float64
	var itemAmount float64 = 0
	for n := 0; n < len(sc); n++ {
		itemAmount = 0
		for i := 0; i < len(items); i++ {
			if sc[n].orderItem == items[i].Name {
				itemAmount = (items[i].UnitPrice * sc[n].qty)
				totalAmount = totalAmount + itemAmount
			}
		}
	}
	return totalAmount
}

func isDayAvailable(day int) bool { //this function checks if input day is still available for booking and returns true if available
	for day > 0 && day < 8 {
		if weeklySchedule[day-1].currentOrders < weeklySchedule[day-1].maxOrders {
			return true
		}
		return false

	}
	return false
}

func (p *orderLL) addOrder(oi orderInfo) { //addOrder allows addition of orders into order linkedlist
	mutex.Lock()
	{
		newNode := &orderNode{
			order: oi,
			next:  nil,
		}
		if p.head == nil {
			p.head = newNode
		} else {
			currentNode := p.head
			for currentNode.next != nil {
				currentNode = currentNode.next
			}
			currentNode.next = newNode
		}
		p.size++
	}
	mutex.Unlock()
	return
}
