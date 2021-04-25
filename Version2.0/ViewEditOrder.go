package main

import (
	"net/http"
	"strconv"
)

func (p *orderLL) viewOrEditOrder(res http.ResponseWriter, req *http.Request) { //view or edit function that allows user to edit orders if they are logged in and have current orders
	myUser := getUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	ptr := p.head
	var tempUserName string // this is used to check if no order username matches current user
	for i := 0; i < p.size; i++ {
		if ptr.order.Username == myUser.Username {
			tempUserName = ptr.order.Username
		}
		ptr = ptr.next
	}
	if tempUserName == "" { // if current logged user has no orders, an error message will be issued
		http.Error(res, "You have no current orders", http.StatusBadRequest)
		return
	}

	var newShoppingCart []order
	var oNum64 int64
	var name string
	var add string
	var dday int

	if req.Method == http.MethodPost {
		orderNum, _ := strconv.Atoi(req.FormValue("orderNum"))
		oNum64 = int64(orderNum) // convert orderNum to type int64
		name = req.FormValue("name")
		add = req.FormValue("address")
		dday, _ = strconv.Atoi(req.FormValue("dday"))

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

	}
	if oNum64 != 0 {

		ptr := p.head

		if ptr == nil {
			http.Error(res, "The booking list is empty, nothing to edit", http.StatusBadRequest)
			return
		}

		for i := 0; i < p.size; i++ {
			if ptr.order.OrderNum == oNum64 {
				if ptr.order.Username != myUser.Username {
					http.Error(res, "You are not authorized to edit other's bookings!", http.StatusUnauthorized)
					return
				}
				if name == "" {
					name := ptr.order.Name
					ptr.order.Name = name
				} else {
					ptr.order.Name = name
				}

				if add == "" {
					add = ptr.order.Address
					ptr.order.Address = add
				} else {
					ptr.order.Address = add
				}

				if dday == 0 {
					dday = ptr.order.DeliveryDay
					ptr.order.DeliveryDay = dday
				} else {
					ptr.order.DeliveryDay = dday
				}

				if len(newShoppingCart) == 0 {
					newShoppingCart = ptr.order.ShoppingCart
					ptr.order.ShoppingCart = newShoppingCart
				} else {
					ptr.order.ShoppingCart = newShoppingCart
				}
				amt := calculateAmount(newShoppingCart)
				ptr.order.Amount = amt
				//orderList.printAllOrderNodes()

				updateWeeklySchedule(orderList)
				updateWeeklyOrder(orderList)

				break
			}

			if ptr.next == nil && ptr.order.OrderNum != oNum64 {
				http.Error(res, "The order number cannot be found, nothing to view or edit", http.StatusBadRequest)
				return
			}
			ptr = ptr.next

		}

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

	ptr = p.head
	var userOrder []orderInfo

	for i := 0; i < p.size; i++ {
		if ptr.order.Username == myUser.Username {
			userOrder = append(userOrder, ptr.order) //this slice of order info will be shown to client on edit page
		}
		ptr = ptr.next
	}

	//orderList.printAllOrderNodes()
	tpl.ExecuteTemplate(res, "viewOrEdit.gohtml", itemData)
	tpl.ExecuteTemplate(res, "yourOrder.gohtml", userOrder) // prints user order detail onto client side
	viewAvailableDays(res, req)
	showRemainingBalance(res, req)

}
