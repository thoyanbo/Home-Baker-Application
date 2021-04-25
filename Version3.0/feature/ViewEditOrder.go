package feature

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	ds "HomeBakerAppGIA2/datastruct"
	ses "HomeBakerAppGIA2/session"

	log "github.com/sirupsen/logrus"
)

//ViewOrEditOrder function allows user to edit orders if they are logged in and have existing orders. Details entered will be edited and if empty would be un-edited.
func ViewOrEditOrder(res http.ResponseWriter, req *http.Request) {
	defer func() { //to handle potential panic situation
		if err := recover(); err != nil {
			log.Panic("Panic occured at view or edit order:", err)
		}
	}()

	myUser := ses.GetUser(res, req)
	if !ses.AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	ptr := ds.OrderList.Head
	var tempUserName string // this is used to check if no order username matches current user
	for i := 0; i < ds.OrderList.Size; i++ {
		if ptr.Order.Username == myUser.Username {
			tempUserName = ptr.Order.Username
		}
		ptr = ptr.Next
	}
	if tempUserName == "" { // if current logged user has no orders, an error message will be issued
		http.Error(res, "You have no current orders", http.StatusBadRequest)
		return
	}

	var newShoppingCart []ds.Order
	var oNum64 int64
	var name string
	var add string
	var dday int

	if req.Method == http.MethodPost {
		orderNumRegExp := regexp.MustCompile(`^[0-9]{1,}$`) //order number reg exp to check for number pattern match
		if !orderNumRegExp.MatchString(req.FormValue("orderNum")) {
			http.Error(res, "Non-numbers detected, please enter a valid order number pattern.", http.StatusBadRequest)
			log.Warning("Invalid user input for order number field")
			return
		}
		orderNum, _ := strconv.Atoi(req.FormValue("orderNum"))
		oNum64 = int64(orderNum) // convert orderNum to type int64

		nameRegExp := regexp.MustCompile(`^([\w'\-,.\s]|)[^0-9_!¡?÷?¿/\\+=@#$%ˆ&*(){}|~<>;:[\]]{0,30}$`) //name regexp to check for name pattern match
		name = strings.TrimSpace(req.FormValue("name"))
		if !nameRegExp.MatchString(name) {
			http.Error(res, "You have entered an invalid name field.", http.StatusBadRequest)
			log.Warning("Invalid user input for name field")
			return
		}
		name = pol.Sanitize(name)

		addRegExp := regexp.MustCompile(`^([\w'\-,.]|)[^_!¡?÷?¿/\\+=$%ˆ&*(){}|~<>;:[\]]{0,100}$`) //address regexp to check for address pattern match
		add = strings.TrimSpace(req.FormValue("address"))
		if !addRegExp.MatchString(add) {
			http.Error(res, "You have entered an invalid address.", http.StatusBadRequest)
			log.Warning("Invalid user input for address field")
			return
		}
		add = pol.Sanitize(add)

		sday := req.FormValue("dday") //sday is string day
		dayRegExp := regexp.MustCompile(`^[1-7]?$`)
		if !dayRegExp.MatchString(sday) {
			http.Error(res, "You have entered an invalid delivery day.", http.StatusBadRequest)
			log.Warning("Invalid user input for delivery day field")
			return
		}

		dday, _ = strconv.Atoi(sday)

		if dday > 0 && dday < 8 {
			availableDay := ds.IsDayAvailable(dday)
			if availableDay == false { //this checks if the order was placed on an unavailable day
				errorString := "Sorry! There are no more available delivery slots for " + ds.IntToDay(dday)
				http.Error(res, errorString, http.StatusBadRequest)
				log.Warning("There are no more available delivery slots for ", ds.IntToDay(dday))
				return
			}
		}

		orderQtyRegExp := regexp.MustCompile(`^[0-9]{0,2}$`) //order quantity reg exp to check for quantity pattern match

		for i := 0; i < len(ds.Items); i++ {
			if !orderQtyRegExp.MatchString(req.FormValue(ds.Items[i].Name)) {
				errorString := "You have entered an invalid order quantity for " + ds.Items[i].Name + "."
				http.Error(res, errorString, http.StatusBadRequest)
				log.Error("Invalid order quantity entered into html form.")
				return
			}
			quantity, _ := strconv.Atoi(req.FormValue(ds.Items[i].Name))
			quantity64 := float64(quantity)

			if quantity64 > 0 {
				itemAvailable := availableItem(ds.Items[i].Name)
				if itemAvailable == false { // this checks if the current item is in stock
					errorString := "Oops, " + ds.Items[i].Name + " is no longer available for ordering."
					http.Error(res, errorString, http.StatusBadRequest)
					log.Warning("Item unavailable: ", ds.Items[i].Name)
					return
				}

				availableBalance := isBalanceEnough(ds.Items[i].Name, quantity64)
				if availableBalance == false { //this checks if the user over ordered on the item
					errorString := "Oops, there is no sufficient balance of" + ds.Items[i].Name + " for ordering.."
					http.Error(res, errorString, http.StatusBadRequest)
					log.Warning("Insufficient balance for ", ds.Items[i].Name)
					return
				}

				singleCart := ds.Order{
					OrderItem: ds.Items[i].Name,
					Qty:       quantity64,
				}
				newShoppingCart = append(newShoppingCart, singleCart)
			}
		}

	}

	if oNum64 != 0 {

		ptr := ds.OrderList.Head

		if ptr == nil {
			http.Error(res, "The booking list is empty, nothing to edit", http.StatusBadRequest)
			log.Warning("Empty booking list, nothing to edit")
			return
		}

		for i := 0; i < ds.OrderList.Size; i++ {
			if ptr.Order.OrderNum == oNum64 {
				if ptr.Order.Username != myUser.Username {
					http.Error(res, "You are not authorized to edit other's bookings!", http.StatusUnauthorized)
					log.Warning("User unauthorized attempt to edit other bookings")
					return
				}
				if name == "" {
					name := ptr.Order.Name
					ptr.Order.Name = name
				} else {
					ptr.Order.Name = name
				}

				if add == "" {
					add = ptr.Order.Address
					ptr.Order.Address = add
				} else {
					ptr.Order.Address = add
				}

				if dday == 0 {
					dday = ptr.Order.DeliveryDay
					ptr.Order.DeliveryDay = dday
				} else {
					ptr.Order.DeliveryDay = dday
				}

				if len(newShoppingCart) == 0 {
					newShoppingCart = ptr.Order.ShoppingCart
					ptr.Order.ShoppingCart = newShoppingCart
				} else {
					ptr.Order.ShoppingCart = newShoppingCart
				}
				amt := ds.CalculateAmount(newShoppingCart)
				ptr.Order.Amount = amt
				//orderList.printAllOrderNodes()

				UpdateWeeklySchedule(ds.OrderList)
				UpdateWeeklyOrder(ds.OrderList)

				break
			}

			if ptr.Next == nil && ptr.Order.OrderNum != oNum64 {
				http.Error(res, "The order number cannot be found, nothing to view or edit", http.StatusBadRequest)
				log.Warning("Order number cannot be found")
				return
			}
			ptr = ptr.Next

		}

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

	ptr = ds.OrderList.Head
	var userOrder []ds.OrderInfo

	for i := 0; i < ds.OrderList.Size; i++ {
		if ptr.Order.Username == myUser.Username {
			userOrder = append(userOrder, ptr.Order) //this slice of order info will be shown to client on edit page
		}
		ptr = ptr.Next
	}

	//orderList.printAllOrderNodes()
	tpl.ExecuteTemplate(res, "viewOrEdit.gohtml", itemData)
	tpl.ExecuteTemplate(res, "yourOrder.gohtml", userOrder) // prints user order detail onto client side
	ViewAvailableDays(res, req)
	showRemainingBalance(res, req)

}
