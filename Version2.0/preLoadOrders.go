package main

import (
	"sync"
	"sync/atomic"

	"golang.org/x/crypto/bcrypt"
)

//These functions created for pre-loading of orders and users into the application so that it does not start empty.

func preLoadUser() {
	bPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	mapUsers["admin"] = user{"admin", bPassword, "admin", "admin"}
	bPassword1, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	mapUsers["yanbo"] = user{"yanbo", bPassword1, "yanbo", "tho"}
	bPassword2, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	mapUsers["ironman"] = user{"ironman", bPassword2, "tony", "stark"}
}

func preLoadOrders() {
	var wg sync.WaitGroup
	wg.Add(7)
	go func() {
		preLoadOrders1()
		wg.Done()
	}()
	go func() {
		preLoadOrders2()
		wg.Done()
	}()
	go func() {
		preLoadOrders3()
		wg.Done()
	}()
	go func() {
		preLoadOrders4()
		wg.Done()
	}()
	go func() {
		preLoadOrders5()
		wg.Done()
	}()
	go func() {
		preLoadOrders6()
		wg.Done()
	}()
	go func() {
		preLoadOrders7()
		wg.Done()
	}()
	wg.Wait()
}

func preLoadOrders1() {
	O1 := orderInfo{
		Username:    "JackNeo",
		Name:        "Jack Neo",
		Address:     "Caldecott Apartment 1",
		DeliveryDay: 6,
		OrderNum:    atomic.AddInt64(&orderNumber, 1),
		ShoppingCart: []order{
			{orderItem: "Blackforest Cake", qty: 7},
		},
		Amount: calculateAmount([]order{{orderItem: "Blackforest Cake", qty: 7}}),
	}
	orderList.addOrder(O1)
}
func preLoadOrders2() {

	O2 := orderInfo{
		Username:    "MarkLee",
		Name:        "Mark Lee",
		Address:     "Caldecott Apartment 2",
		DeliveryDay: 6,
		OrderNum:    atomic.AddInt64(&orderNumber, 1),
		ShoppingCart: []order{
			{orderItem: "Burnt Cheesecake", qty: 4},
			{orderItem: "Pineapple Tarts", qty: 4},
		},
		Amount: calculateAmount([]order{{orderItem: "Burnt Cheesecake", qty: 7}, {orderItem: "Pineapple Tarts", qty: 4}}),
	}
	orderList.addOrder(O2)
}
func preLoadOrders3() {
	O3 := orderInfo{
		Username:    "ironman",
		Name:        "Tony Stark",
		Address:     "Stark Towers",
		DeliveryDay: 6,
		OrderNum:    atomic.AddInt64(&orderNumber, 1),
		ShoppingCart: []order{
			{orderItem: "Mochi Buns", qty: 3},
			{orderItem: "Pineapple Tarts", qty: 3},
		},
		Amount: calculateAmount([]order{{orderItem: "Mochi Buns", qty: 3}, {orderItem: "Pineapple Tarts", qty: 3}}),
	}
	orderList.addOrder(O3)
}
func preLoadOrders4() {
	O4 := orderInfo{
		Username:    "Batman",
		Name:        "Bruce Wayne",
		Address:     "Gotham City",
		DeliveryDay: 6,
		OrderNum:    atomic.AddInt64(&orderNumber, 1),
		ShoppingCart: []order{
			{orderItem: "Red Velvet Carrot Cake", qty: 2},
			{orderItem: "Kaya Buns", qty: 7},
		},
		Amount: calculateAmount([]order{{orderItem: "Red Velvet Carrot Cake", qty: 2}, {orderItem: "Kaya Buns", qty: 7}}),
	}
	orderList.addOrder(O4)
}
func preLoadOrders5() {

	O5 := orderInfo{
		Username:    "yanbo",
		Name:        "Yan Bo",
		Address:     "Choa Chu Kang",
		DeliveryDay: 6,
		OrderNum:    atomic.AddInt64(&orderNumber, 1),
		ShoppingCart: []order{
			{orderItem: "Mochi Buns", qty: 10},
			{orderItem: "Burnt Cheesecake", qty: 10},
			{orderItem: "Kaya Buns", qty: 5},
		},
		Amount: calculateAmount([]order{{orderItem: "Mochi Buns", qty: 10}, {orderItem: "Burnt Cheesecake", qty: 10}, {orderItem: "Kaya Buns", qty: 5}}),
	}
	orderList.addOrder(O5)
}

func preLoadOrders6() {

	O6 := orderInfo{
		Username:    "jackie",
		Name:        "Jackie Chan",
		Address:     "Hollywood Mansion",
		DeliveryDay: 3,
		OrderNum:    atomic.AddInt64(&orderNumber, 1),
		ShoppingCart: []order{
			{orderItem: "Red Velvet Carrot Cake", qty: 10},
			{orderItem: "Kaya Buns", qty: 5},
		},
		Amount: calculateAmount([]order{{orderItem: "Red Velvet Carrot Cake", qty: 10}, {orderItem: "Kaya Buns", qty: 5}}),
	}
	orderList.addOrder(O6)
}

func preLoadOrders7() {

	O7 := orderInfo{
		Username:    "LTK",
		Name:        "Low Thia Khiang",
		Address:     "Hougang",
		DeliveryDay: 4,
		OrderNum:    atomic.AddInt64(&orderNumber, 1),
		ShoppingCart: []order{
			{orderItem: "Mochi Buns", qty: 10},
			{orderItem: "Kaya Buns", qty: 5},
		},
		Amount: calculateAmount([]order{{orderItem: "Mochi Buns", qty: 10}, {orderItem: "Kaya Buns", qty: 5}}),
	}
	orderList.addOrder(O7)
}
