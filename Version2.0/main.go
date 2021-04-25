package main

import (
	"html/template"
	"net/http"
	"sync"
)

var tpl *template.Template
var mutex sync.Mutex

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	preLoadUser()                   //preload users
	preLoadOrders()                 // this is done to load some existing orders into the program
	updateWeeklySchedule(orderList) //update weekly schedule list
	updateWeeklyOrder(orderList)    //update weekly item list
}

func main() {
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/", index)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/menu", menu)
	http.HandleFunc("/display", display)
	http.HandleFunc("/search", searchMenu)
	http.HandleFunc("/viewPrice", viewPrice)
	http.HandleFunc("/viewCategory", viewCategory)
	http.HandleFunc("/recommended", viewRecommended)
	http.HandleFunc("/createOrder", createNewOrder)
	http.HandleFunc("/deleteOrder", orderList.delOrder)
	http.HandleFunc("/viewOrEdit", orderList.viewOrEditOrder)
	http.HandleFunc("/overview", orderList.viewAllOrders)
	http.ListenAndServe(":5221", nil)
}

func menu(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	tpl.ExecuteTemplate(res, "menu.gohtml", myUser)
}
