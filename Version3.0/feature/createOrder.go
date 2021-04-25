package feature

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"

	ds "HomeBakerAppGIA2/datastruct"
	ses "HomeBakerAppGIA2/session"

	log "github.com/sirupsen/logrus"
)

//CreateNewOrder takes in user information through a HTML form request to create a new order.
func CreateNewOrder(res http.ResponseWriter, req *http.Request) {
	defer func() { //to handle potential panic situation
		if err := recover(); err != nil {
			log.Panic("Panic occured at create order:", err)
		}
	}()
	myUser := ses.GetUser(res, req)
	if !ses.AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	//fmt.Println(Items)
	sortItems(ds.Items)

	var newShoppingCart = []ds.Order{}

	if req.Method == http.MethodPost {
		nameRegExp := regexp.MustCompile(`^[\w'\-,.][^0-9_!¡?÷?¿/\\+=@#$%ˆ&*(){}|~<>;:[\]]{2,30}$`) //name regexp to check for name pattern match
		name := strings.TrimSpace(req.FormValue("name"))
		if !nameRegExp.MatchString(name) {
			http.Error(res, "You have entered an invalid name field.", http.StatusBadRequest)
			log.Warning("Invalid user input for name field")
			return
		}
		name = pol.Sanitize(name) //pol.Sanitize is used to sanitize inputs

		addRegExp := regexp.MustCompile(`^[\w'\-,.][^_!¡?÷?¿/\\+=$%ˆ&*(){}|~<>;:[\]]{2,100}$`) ////address regexp to check for address pattern match
		add := strings.TrimSpace(req.FormValue("address"))
		if !addRegExp.MatchString(add) {
			http.Error(res, "You have entered an invalid address.", http.StatusBadRequest)
			log.Warning("Invalid user input for address field")
			return
		}
		add = pol.Sanitize(add) //pol.Sanitize is used to sanitize inputs

		sday := req.FormValue("dday") //sday is string day
		dayRegExp := regexp.MustCompile(`^[1-7]$`)
		if !dayRegExp.MatchString(sday) {
			http.Error(res, "You have entered an invalid delivery day.", http.StatusBadRequest)
			log.Warning("Invalid user input for delivery day")
			return
		}

		dday, _ := strconv.Atoi(sday)

		availableDay := ds.IsDayAvailable(dday)
		if availableDay == false { //this checks if the order was placed on an unavailable day
			errorString := "Sorry! There are no more available delivery slots for " + ds.IntToDay(dday)
			http.Error(res, errorString, http.StatusBadRequest)
			log.Warning("There are no more available delivery slots for " + ds.IntToDay(dday))
			return
		}

		orderQtyRegExp := regexp.MustCompile(`^[0-9]{1,2}$`) //order quantity reg exp to check for quantity pattern match

		for i := 0; i < len(ds.Items); i++ {
			if !orderQtyRegExp.MatchString(req.FormValue(ds.Items[i].Name)) {
				errorString := "You have entered an invalid order quantity for " + ds.Items[i].Name + "."
				http.Error(res, errorString, http.StatusBadRequest)
				log.Warning("Invalid user input for order quantity")
				return
			}
			quantity, _ := strconv.Atoi(req.FormValue(ds.Items[i].Name)) //label for the form input is the item name, but returns a quantity of that item
			quantity64 := float64(quantity)

			if quantity64 > 0 {
				itemAvailable := availableItem(ds.Items[i].Name)
				if itemAvailable == false { // this checks if the current item is in stock
					errorString := "Oops, " + ds.Items[i].Name + " is no longer available for ordering."
					http.Error(res, errorString, http.StatusBadRequest)
					log.Warning("User overordered on item:", ds.Items[i].Name)
					return
				}
				availableBalance := isBalanceEnough(ds.Items[i].Name, quantity64)
				if availableBalance == false { //this checks if the user over ordered on the item
					errorString := "Oops, there is no sufficient balance of" + ds.Items[i].Name + " for ordering.."
					http.Error(res, errorString, http.StatusBadRequest)
					log.Warning("User overordered on item:", ds.Items[i].Name)
					return
				}
				singleCart := ds.Order{
					OrderItem: ds.Items[i].Name,
					Qty:       quantity64,
				}
				newShoppingCart = append(newShoppingCart, singleCart)
			}
		}

		if len(newShoppingCart) == 0 {
			http.Error(res, "Error: You cannot submit an empty shopping cart.", http.StatusBadRequest)
			log.Warning("User entered empty shopping cart.")
			return
		}

		on := atomic.AddInt64(&OrderNumber, 1) // use of atomic function to prevent multiple clients from possibly creating identical order number
		amt := ds.CalculateAmount(newShoppingCart)
		newOrder := ds.OrderInfo{
			Username:     myUser.Username,
			Name:         name,
			Address:      add,
			DeliveryDay:  dday,
			OrderNum:     on,
			ShoppingCart: newShoppingCart,
			Amount:       amt,
		}

		ds.OrderList.AddOrder(newOrder)
		err := UpdateWeeklySchedule(ds.OrderList)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			log.Error(err)
			return
		}

		err = UpdateWeeklyOrder(ds.OrderList)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			log.Error(err)
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

	for i := 0; i < len(ds.Items); i++ {
		remainingQuantity := ds.Items[i].WeeklyCapacity - ds.Items[i].WeeklyOrder
		d := balanceStruct{
			Item:     ds.Items[i].Name,
			Quantity: remainingQuantity}
		itemData = append(itemData, d)
	}

	err := tpl.ExecuteTemplate(res, "createOrder.gohtml", itemData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Fatalln(err)
	}

	ViewAvailableDays(res, req)
	showRemainingBalance(res, req)
}

//Function ViewAvailableDays to show remaining available days for booking to user during create/edit order page
func ViewAvailableDays(res http.ResponseWriter, req *http.Request) {
	type dayData struct {
		Day          string
		RemainingOrd int
	}
	var AvailableDaysData []dayData

	for n := 0; n < len(ds.WeeklySchedule); n++ {
		if ds.WeeklySchedule[n].CurrentOrders < ds.WeeklySchedule[n].MaxOrders {
			remainingOrders := ds.WeeklySchedule[n].MaxOrders - ds.WeeklySchedule[n].CurrentOrders
			if remainingOrders > 0 {
				tempData := dayData{
					Day:          ds.IntToDay(n + 1),
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

//Function to show remaining item balance to user during create/edit order page
func showRemainingBalance(res http.ResponseWriter, req *http.Request) {
	type balanceStruct struct {
		Item     string
		Quantity float64
	}

	var itemData []balanceStruct

	for i := 0; i < len(ds.Items); i++ {
		remainingQuantity := ds.Items[i].WeeklyCapacity - ds.Items[i].WeeklyOrder
		d := balanceStruct{
			Item:     ds.Items[i].Name,
			Quantity: remainingQuantity}
		itemData = append(itemData, d)
	}

	err := tpl.ExecuteTemplate(res, "itemBalance.gohtml", itemData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Error("Error at itemBalance template.", err)
		return
	}
}

//Function to check if input item still has weekly stock for ordering
func availableItem(nam string) (exist bool) {
	for n := 0; n < len(ds.Items); n++ {
		if nam == ds.Items[n].Name {
			if ds.Items[n].WeeklyCapacity > ds.Items[n].WeeklyOrder {
				return true
			}
			return false
		}
	}
	return false
}

//Function checks if the order quantity exceeds the balance quantity, and returns false if true
func isBalanceEnough(nam string, qty float64) bool {
	var balanceQty float64
	for n := 0; n < len(ds.Items); n++ {
		if nam == ds.Items[n].Name {
			balanceQty = ds.Items[n].WeeklyCapacity - ds.Items[n].WeeklyOrder
			break
		}
	}

	if qty <= balanceQty {
		return true
	}
	return false
}
