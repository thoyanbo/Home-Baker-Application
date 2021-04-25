package main

import (
	"encoding/json"
	"io/ioutil"
	"sync"
	"sync/atomic"

	ds "HomeBakerAppGIA2/datastruct"
	feat "HomeBakerAppGIA2/feature"
	ses "HomeBakerAppGIA2/session"

	log "github.com/sirupsen/logrus"
)

//These functions created for pre-loading of orders and users into the application so that it does not start empty.
//For demo purpose, please login with any of the following credentials.
//username: admin, password: password
//username: yanbo, password: password
//username: ironman, password: password
func preLoadUser() {

	fileout, _ := ioutil.ReadFile("users.json")

	//var Roomout []roomInfo

	err := json.Unmarshal([]byte(fileout), &ses.MapUsers)
	if err != nil {
		log.Error(err)
	}

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
	O1 := ds.OrderInfo{
		Username:    "JackNeo",
		Name:        "Jack Neo",
		Address:     "Caldecott Apartment 1",
		DeliveryDay: 6,
		OrderNum:    atomic.AddInt64(&feat.OrderNumber, 1),
		ShoppingCart: []ds.Order{
			{OrderItem: "Blackforest Cake", Qty: 7},
		},
		Amount: ds.CalculateAmount([]ds.Order{{OrderItem: "Blackforest Cake", Qty: 7}}),
	}
	ds.OrderList.AddOrder(O1)
}
func preLoadOrders2() {

	O2 := ds.OrderInfo{
		Username:    "MarkLee",
		Name:        "Mark Lee",
		Address:     "Caldecott Apartment 2",
		DeliveryDay: 6,
		OrderNum:    atomic.AddInt64(&feat.OrderNumber, 1),
		ShoppingCart: []ds.Order{
			{OrderItem: "Burnt Cheesecake", Qty: 4},
			{OrderItem: "Pineapple Tarts", Qty: 4},
		},
		Amount: ds.CalculateAmount([]ds.Order{{OrderItem: "Burnt Cheesecake", Qty: 7}, {OrderItem: "Pineapple Tarts", Qty: 4}}),
	}
	ds.OrderList.AddOrder(O2)
}
func preLoadOrders3() {
	O3 := ds.OrderInfo{
		Username:    "ironman",
		Name:        "Tony Stark",
		Address:     "Stark Towers",
		DeliveryDay: 6,
		OrderNum:    atomic.AddInt64(&feat.OrderNumber, 1),
		ShoppingCart: []ds.Order{
			{OrderItem: "Mochi Buns", Qty: 3},
			{OrderItem: "Pineapple Tarts", Qty: 3},
		},
		Amount: ds.CalculateAmount([]ds.Order{{OrderItem: "Mochi Buns", Qty: 3}, {OrderItem: "Pineapple Tarts", Qty: 3}}),
	}
	ds.OrderList.AddOrder(O3)
}
func preLoadOrders4() {
	O4 := ds.OrderInfo{
		Username:    "Batman",
		Name:        "Bruce Wayne",
		Address:     "Gotham City",
		DeliveryDay: 6,
		OrderNum:    atomic.AddInt64(&feat.OrderNumber, 1),
		ShoppingCart: []ds.Order{
			{OrderItem: "Red Velvet Carrot Cake", Qty: 2},
			{OrderItem: "Kaya Buns", Qty: 7},
		},
		Amount: ds.CalculateAmount([]ds.Order{{OrderItem: "Red Velvet Carrot Cake", Qty: 2}, {OrderItem: "Kaya Buns", Qty: 7}}),
	}
	ds.OrderList.AddOrder(O4)
}
func preLoadOrders5() {

	O5 := ds.OrderInfo{
		Username:    "yanbo",
		Name:        "Yan Bo",
		Address:     "Choa Chu Kang",
		DeliveryDay: 6,
		OrderNum:    atomic.AddInt64(&feat.OrderNumber, 1),
		ShoppingCart: []ds.Order{
			{OrderItem: "Mochi Buns", Qty: 10},
			{OrderItem: "Burnt Cheesecake", Qty: 10},
			{OrderItem: "Kaya Buns", Qty: 5},
		},
		Amount: ds.CalculateAmount([]ds.Order{{OrderItem: "Mochi Buns", Qty: 10}, {OrderItem: "Burnt Cheesecake", Qty: 10}, {OrderItem: "Kaya Buns", Qty: 5}}),
	}
	ds.OrderList.AddOrder(O5)
}

func preLoadOrders6() {

	O6 := ds.OrderInfo{
		Username:    "jackie",
		Name:        "Jackie Chan",
		Address:     "Hollywood Mansion",
		DeliveryDay: 3,
		OrderNum:    atomic.AddInt64(&feat.OrderNumber, 1),
		ShoppingCart: []ds.Order{
			{OrderItem: "Red Velvet Carrot Cake", Qty: 10},
			{OrderItem: "Kaya Buns", Qty: 5},
		},
		Amount: ds.CalculateAmount([]ds.Order{{OrderItem: "Red Velvet Carrot Cake", Qty: 10}, {OrderItem: "Kaya Buns", Qty: 5}}),
	}
	ds.OrderList.AddOrder(O6)
}

func preLoadOrders7() {

	O7 := ds.OrderInfo{
		Username:    "LTK",
		Name:        "Low Thia Khiang",
		Address:     "Hougang",
		DeliveryDay: 4,
		OrderNum:    atomic.AddInt64(&feat.OrderNumber, 1),
		ShoppingCart: []ds.Order{
			{OrderItem: "Mochi Buns", Qty: 10},
			{OrderItem: "Kaya Buns", Qty: 5},
		},
		Amount: ds.CalculateAmount([]ds.Order{{OrderItem: "Mochi Buns", Qty: 10}, {OrderItem: "Kaya Buns", Qty: 5}}),
	}
	ds.OrderList.AddOrder(O7)
}
