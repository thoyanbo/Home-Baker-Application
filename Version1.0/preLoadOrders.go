package main

import (
	"sync"
	"sync/atomic"
)

//These functions created for pre-loading of existing data into the data structure so that data structure is not empty.

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
		name:        "Jack Neo",
		address:     "Caldecott Apartment 1",
		deliveryDay: 6,
		orderNum:    atomic.AddInt64(&orderNumber, 1),
		shoppingCart: []order{
			{orderItem: "Blackforest Cake", qty: 7},
		},
		amount: calculateAmount([]order{{orderItem: "Blackforest Cake", qty: 7}}),
	}
	orderList.addOrder(O1)
}
func preLoadOrders2() {

	O2 := orderInfo{
		name:        "Mark Lee",
		address:     "Caldecott Apartment 2",
		deliveryDay: 6,
		orderNum:    atomic.AddInt64(&orderNumber, 1),
		shoppingCart: []order{
			{orderItem: "Burnt Cheesecake", qty: 4},
			{orderItem: "Pineapple Tarts", qty: 4},
		},
		amount: calculateAmount([]order{{orderItem: "Burnt Cheesecake", qty: 7}, {orderItem: "Pineapple Tarts", qty: 4}}),
	}
	orderList.addOrder(O2)
}
func preLoadOrders3() {
	O3 := orderInfo{
		name:        "Tony Stark",
		address:     "Stark Towers",
		deliveryDay: 6,
		orderNum:    atomic.AddInt64(&orderNumber, 1),
		shoppingCart: []order{
			{orderItem: "Mochi Buns", qty: 3},
			{orderItem: "Pineapple Tarts", qty: 3},
		},
		amount: calculateAmount([]order{{orderItem: "Mochi Buns", qty: 3}, {orderItem: "Pineapple Tarts", qty: 3}}),
	}
	orderList.addOrder(O3)
}
func preLoadOrders4() {
	O4 := orderInfo{
		name:        "Bruce Wayne",
		address:     "Gotham City",
		deliveryDay: 6,
		orderNum:    atomic.AddInt64(&orderNumber, 1),
		shoppingCart: []order{
			{orderItem: "Red Velvet Carrot Cake", qty: 2},
			{orderItem: "Kaya Buns", qty: 7},
		},
		amount: calculateAmount([]order{{orderItem: "Red Velvet Carrot Cake", qty: 2}, {orderItem: "Kaya Buns", qty: 7}}),
	}
	orderList.addOrder(O4)
}
func preLoadOrders5() {

	O5 := orderInfo{
		name:        "PM Lee",
		address:     "Istana",
		deliveryDay: 6,
		orderNum:    atomic.AddInt64(&orderNumber, 1),
		shoppingCart: []order{
			{orderItem: "Mochi Buns", qty: 10},
			{orderItem: "Burnt Cheesecake", qty: 10},
			{orderItem: "Kaya Buns", qty: 5},
		},
		amount: calculateAmount([]order{{orderItem: "Mochi Buns", qty: 10}, {orderItem: "Burnt Cheesecake", qty: 10}, {orderItem: "Kaya Buns", qty: 5}}),
	}
	orderList.addOrder(O5)
}

func preLoadOrders6() {

	O6 := orderInfo{
		name:        "Jackie Chan",
		address:     "Hollywood Mansion",
		deliveryDay: 3,
		orderNum:    atomic.AddInt64(&orderNumber, 1),
		shoppingCart: []order{
			{orderItem: "Red Velvet Carrot Cake", qty: 10},
			{orderItem: "Kaya Buns", qty: 5},
		},
		amount: calculateAmount([]order{{orderItem: "Red Velvet Carrot Cake", qty: 10}, {orderItem: "Kaya Buns", qty: 5}}),
	}
	orderList.addOrder(O6)
}

func preLoadOrders7() {

	O7 := orderInfo{
		name:        "Low Thia Khiang",
		address:     "Hougang",
		deliveryDay: 4,
		orderNum:    atomic.AddInt64(&orderNumber, 1),
		shoppingCart: []order{
			{orderItem: "Mochi Buns", qty: 10},
			{orderItem: "Kaya Buns", qty: 5},
		},
		amount: calculateAmount([]order{{orderItem: "Mochi Buns", qty: 10}, {orderItem: "Kaya Buns", qty: 5}}),
	}
	orderList.addOrder(O7)
}
