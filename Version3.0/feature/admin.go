//Package feature contains the functions and variables that are used in Home Baker application.
//		The functions described in this package are mainly used for:
//		1. Admin features of Home Baker application.
//		2. Create new order
//		3. Display Items based on pre-determined patterns
//		4. Search item function
//		5. View and edit existing order.
package feature

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"sync"
	"text/template"

	ds "HomeBakerAppGIA2/datastruct"
	ses "HomeBakerAppGIA2/session"

	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
)

// OrderNumber is the global variable to track order number
var OrderNumber int64 = 0
var tpl = template.Must(template.ParseGlob("templates/*"))
var mutex sync.Mutex

// Do this once for each unique policy, and use the policy for the life of the program
// Policy creation/editing is not safe to use in multiple goroutines
var pol = bluemonday.UGCPolicy() //pol for policy

//ViewAllOrders shows an overview for all orders for the admin, and is only accessible by user "admin"
func ViewAllOrders(res http.ResponseWriter, req *http.Request) {

	myUser := ses.GetUser(res, req)
	if myUser.Username != "admin" {
		http.Error(res, "You are not an admin user", http.StatusUnauthorized)
		log.Warning("Unauthorized access by non admin user.")
		return
	}
	var weeklyRevenue float64 = 0

	for i := 0; i < len(ds.Items); i++ {
		itemRevenue := ds.Items[i].WeeklyOrder * ds.Items[i].UnitPrice
		weeklyRevenue = weeklyRevenue + itemRevenue
	}

	err := tpl.ExecuteTemplate(res, "overview.gohtml", weeklyRevenue)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Error("Error at overview template.", err)
	}

	showDailyOrder(res, req)
	showItemOrder(res, req)
	ViewDailyOrder(res, req)
}

//showItemOrder summarizes the weekly data of invididual items ordered and their order quantity. Executed as a html response.
func showItemOrder(res http.ResponseWriter, req *http.Request) {
	type itemStruct struct {
		ItemName string
		Quantity float64
	}

	var weeklyData []itemStruct

	for i := 0; i < len(ds.Items); i++ {
		if ds.Items[i].WeeklyOrder > 0 {
			weeklyItemData := itemStruct{
				ItemName: ds.Items[i].Name,
				Quantity: ds.Items[i].WeeklyOrder,
			}
			weeklyData = append(weeklyData, weeklyItemData)
		}
	}
	err := tpl.ExecuteTemplate(res, "showItemOrder.gohtml", weeklyData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Error("Error at showItemOrder template.", err)
		return
	}
}

//showDailyOrder summarizes the daily orders (day and number of orders for that day). Executed as a html response.
func showDailyOrder(res http.ResponseWriter, req *http.Request) {
	type orderStruct struct {
		Day        string
		DailyOrder int
	}
	var weeklyData []orderStruct

	for i := 0; i < len(ds.WeeklySchedule); i++ {
		if ds.WeeklySchedule[i].CurrentOrders > 0 {
			dailyData := orderStruct{
				Day:        ds.IntToDay(ds.WeeklySchedule[i].Weekday),
				DailyOrder: ds.WeeklySchedule[i].CurrentOrders,
			}
			weeklyData = append(weeklyData, dailyData)
		}
	}
	err := tpl.ExecuteTemplate(res, "showDailyOrder.gohtml", weeklyData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Error("Error at showDailyOrder template.", err)
		return
	}
}

//Function DelOrder takes in a valid order number and deletes it from the order linked list. Can only be implemented by user "admin".
func DelOrder(res http.ResponseWriter, req *http.Request) {
	defer func() { //to handle potential panic situation
		if err := recover(); err != nil {
			log.Panic("Panic occured at delete items:", err)
		}
	}()

	var delOrder ds.OrderInfo
	ptr := ds.OrderList.Head
	myUser := ses.GetUser(res, req)
	if myUser.Username != "admin" {
		http.Error(res, "You are not an admin user", http.StatusUnauthorized)
		log.Warning("Unauthorized access by non admin user.")
		return
	}
	if req.Method == http.MethodPost {
		bookingNumber, _ := strconv.Atoi(req.FormValue("bookingNumber"))
		bookingNum := int64(bookingNumber)

		if bookingNum != 0 {

			if ptr == nil {
				http.Error(res, "The booking list is empty, nothing to delete", http.StatusBadRequest)
				log.Warning("Booking list is empty.")
				return
			}
			for ptr.Next.Order.OrderNum != bookingNum {
				ptr = ptr.Next
				if ptr.Next == nil {
					http.Error(res, "This booking number does not exist", http.StatusBadRequest)
					log.Warning("Invalid user attempt to reach non-existing booking number.")
					return
				}
			}

			delOrder = ptr.Next.Order
			ptr.Next = ptr.Next.Next
			ds.OrderList.Size--
			ds.OrderList.PrintAllOrderNodes()
			UpdateWeeklySchedule(ds.OrderList)
			UpdateWeeklyOrder(ds.OrderList)

		}
	}

	err := tpl.ExecuteTemplate(res, "deleteOrder.gohtml", delOrder)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Error("Error at deleteOrder template.", err)
		return
	}

	ptr = ds.OrderList.Head
	for i := 0; i < ds.OrderList.Size; i++ {
		err := tpl.ExecuteTemplate(res, "viewAllOrders.gohtml", ptr.Order) // prints order detail onto client side
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			log.Error("Error at viewAllOrders:", err)
			return
		}
		ptr = ptr.Next
	}

}

//Function UpdateWeeklySchedule iterates over the order linked list and stores the number of daily orders into weekly schedule's current order.
func UpdateWeeklySchedule(p *ds.OrderLL) error {

	ptr := p.Head

	if p.Head == nil {
		log.Error("Empty weekly order, unable to update schedule")
		return errors.New("Empty weekly order, unable to update schedule")
	}

	dayCount := make([]int, 7)
	for i := 0; i < p.Size; i++ {
		//fmt.Println(ptr.order.deliveryDay)
		dayCount[ptr.Order.DeliveryDay-1]++
		ptr = ptr.Next
	}
	mutex.Lock() // we lock this portion of the code to prevent multiple clients from having concurrent access to global variable weeklySchedule
	{
		for j := 0; j < len(ds.WeeklySchedule); j++ {
			ds.WeeklySchedule[j].CurrentOrders = dayCount[j]
		}
	}
	mutex.Unlock()
	return nil
}

//Function UpdateWeeklyOrder iterates over the order linked list and tabulates the total count of orders for each sales item
func UpdateWeeklyOrder(p *ds.OrderLL) error {
	mutex.Lock()
	{
		for j := 0; j < len(ds.Items); j++ {
			ds.Items[j].WeeklyOrder = 0
		} // this code resets count of weekly order to 0 so that we will not use previously counted weeklyOrder

		ptr := p.Head

		if p.Head == nil {
			log.Error("Empty weekly order, unable to update orders")
			return errors.New("Empty weekly order, unable to update orders")
		}

		for i := 0; i < p.Size; i++ {
			//fmt.Println(ptr.order.shoppingCart)
			for _, v := range ptr.Order.ShoppingCart {
				for n := 0; n < len(ds.Items); n++ {
					if v.OrderItem == ds.Items[n].Name {
						ds.Items[n].WeeklyOrder += v.Qty
						break
					}
				}
			}
			ptr = ptr.Next
		}
	}
	mutex.Unlock()
	return nil
}

//Function viewDailyOrder takes in the html form value "day" and returns all the orders for the day, if any.
func ViewDailyOrder(res http.ResponseWriter, req *http.Request) {
	var day int
	var dayData []ds.OrderInfo

	if req.Method == http.MethodPost {
		var dayFound bool = false
		day, _ = strconv.Atoi(req.FormValue("day"))

		ptr := ds.OrderList.Head
		for i := 0; i < ds.OrderList.Size; i++ {
			if day == ptr.Order.DeliveryDay {
				dayFound = true
				dayData = append(dayData, ptr.Order)
			}
			ptr = ptr.Next
		}
		if dayFound == false {
			io.WriteString(res, `<h1> There are no orders for the selected day. </h1>`)

		}
	}

	err := tpl.ExecuteTemplate(res, "viewDailyOrder.gohtml", dayData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Error("Error at viewDailyOrder:", err)
		return
	}
}
