package main

import (
	"errors"
	"io"
	"net/http"
	"strconv"
)

func (p *orderLL) viewAllOrders(res http.ResponseWriter, req *http.Request) { //viewAllOrders shows an overview for all orders for the admin

	myUser := getUser(res, req)
	if myUser.Username != "admin" {
		http.Error(res, "You are not an admin user", http.StatusUnauthorized)
		return
	}
	var weeklyRevenue float64 = 0

	for i := 0; i < len(items); i++ {
		itemRevenue := items[i].WeeklyOrder * items[i].UnitPrice
		weeklyRevenue = weeklyRevenue + itemRevenue
	}

	err := tpl.ExecuteTemplate(res, "overview.gohtml", weeklyRevenue)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	showDailyOrder(res, req)
	showItemOrder(res, req)
	orderList.viewDailyOrder(res, req)
}

func showItemOrder(res http.ResponseWriter, req *http.Request) {
	type itemStruct struct {
		ItemName string
		Quantity float64
	}

	var weeklyData []itemStruct

	for i := 0; i < len(items); i++ {
		if items[i].WeeklyOrder > 0 {
			weeklyItemData := itemStruct{
				ItemName: items[i].Name,
				Quantity: items[i].WeeklyOrder,
			}
			weeklyData = append(weeklyData, weeklyItemData)
		}
	}
	err := tpl.ExecuteTemplate(res, "showItemOrder.gohtml", weeklyData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

func showDailyOrder(res http.ResponseWriter, req *http.Request) {
	type orderStruct struct {
		Day        string
		DailyOrder int
	}
	var weeklyData []orderStruct

	for i := 0; i < len(weeklySchedule); i++ {
		if weeklySchedule[i].currentOrders > 0 {
			dailyData := orderStruct{
				Day:        intToDay(weeklySchedule[i].weekday),
				DailyOrder: weeklySchedule[i].currentOrders,
			}
			weeklyData = append(weeklyData, dailyData)
		}
	}
	err := tpl.ExecuteTemplate(res, "showDailyOrder.gohtml", weeklyData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

func (p *orderLL) delOrder(res http.ResponseWriter, req *http.Request) {

	var delOrder orderInfo
	ptr := p.head
	myUser := getUser(res, req)
	if myUser.Username != "admin" {
		http.Error(res, "You are not an admin user", http.StatusUnauthorized)
		return
	}
	if req.Method == http.MethodPost {
		bookingNumber, _ := strconv.Atoi(req.FormValue("bookingNumber"))
		bookingNum := int64(bookingNumber)

		if bookingNum != 0 {

			if ptr == nil {
				http.Error(res, "The booking list is empty, nothing to delete", http.StatusBadRequest)
				return
			}
			for ptr.next.order.OrderNum != bookingNum {
				ptr = ptr.next
				if ptr.next == nil {
					http.Error(res, "This booking number does not exist", http.StatusBadRequest)
					return
				}
			}

			delOrder = ptr.next.order
			ptr.next = ptr.next.next
			p.size--
			orderList.printAllOrderNodes()
			updateWeeklySchedule(orderList)
			updateWeeklyOrder(orderList)

		}
	}

	err := tpl.ExecuteTemplate(res, "deleteOrder.gohtml", delOrder)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	ptr = p.head
	for i := 0; i < p.size; i++ {
		err := tpl.ExecuteTemplate(res, "viewAllOrders.gohtml", ptr.order) // prints order detail onto client side
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		ptr = ptr.next
	}

}

func updateWeeklySchedule(p *orderLL) error { //this function iterates over the order linked list and stores the number of daily orders into weekly schedule's current order

	ptr := p.head

	if p.head == nil {
		return errors.New("empty order list")
	}

	dayCount := make([]int, 7)
	for i := 0; i < p.size; i++ {
		//fmt.Println(ptr.order.deliveryDay)
		dayCount[ptr.order.DeliveryDay-1]++
		ptr = ptr.next
	}
	mutex.Lock() // we lock this portion of the code to prevent multiple clients from having concurrent access to global variable weeklySchedule
	{
		for j := 0; j < len(weeklySchedule); j++ {
			weeklySchedule[j].currentOrders = dayCount[j]
		}
	}
	mutex.Unlock()
	return nil
}

func updateWeeklyOrder(p *orderLL) error { //this function iterates over the order linked list and tabulates the total count of orders for each sales item
	mutex.Lock()
	{
		for j := 0; j < len(items); j++ {
			items[j].WeeklyOrder = 0
		} // this code resets count of weekly order to 0 so that we will not use previously counted weeklyOrder

		ptr := p.head

		if p.head == nil {
			return errors.New("empty order list")
		}

		for i := 0; i < p.size; i++ {
			//fmt.Println(ptr.order.shoppingCart)
			for _, v := range ptr.order.ShoppingCart {
				for n := 0; n < len(items); n++ {
					if v.orderItem == items[n].Name {
						items[n].WeeklyOrder += v.qty
						break
					}
				}
			}
			ptr = ptr.next
		}
	}
	mutex.Unlock()
	return nil
}

func (p *orderLL) viewDailyOrder(res http.ResponseWriter, req *http.Request) {
	var day int
	var dayData []orderInfo

	if req.Method == http.MethodPost {
		var dayFound bool = false
		day, _ = strconv.Atoi(req.FormValue("day"))

		ptr := p.head
		for i := 0; i < p.size; i++ {
			if day == ptr.order.DeliveryDay {
				dayFound = true
				dayData = append(dayData, ptr.order)
			}
			ptr = ptr.next
		}
		if dayFound == false {
			io.WriteString(res, `<h1> There are no orders for the selected day. </h1>`)

		}
	}
	err := tpl.ExecuteTemplate(res, "viewDailyOrder.gohtml", dayData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}
