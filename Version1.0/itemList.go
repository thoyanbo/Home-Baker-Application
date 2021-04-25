package main

import "fmt"

//itemInfo stores information about the item
type itemInfo struct {
	name           string
	unitPrice      float64
	category       string
	recommended    bool
	weeklyCapacity float64
	weeklyOrder    float64
}

//items is the global declared slice of sales items and stores relevant information.
var (
	items = []itemInfo{
		{name: "Pineapple Tarts", unitPrice: 18.00, category: "Snacks", recommended: true, weeklyCapacity: 35, weeklyOrder: 0},
		{name: "Coffee Cookies", unitPrice: 12.00, category: "Snacks", recommended: false, weeklyCapacity: 35, weeklyOrder: 0},
		{name: "Kaya Buns", unitPrice: 5.00, category: "Bread", recommended: true, weeklyCapacity: 70, weeklyOrder: 0},
		{name: "Mochi Buns", unitPrice: 6.00, category: "Bread", recommended: true, weeklyCapacity: 70, weeklyOrder: 0},
		{name: "Walnut Raisin Buns", unitPrice: 5.00, category: "Bread", recommended: false, weeklyCapacity: 70, weeklyOrder: 0},
		{name: "4 Cheese Rolls", unitPrice: 7.50, category: "Bread", recommended: false, weeklyCapacity: 70, weeklyOrder: 0},
		{name: "Blackforest Cake", unitPrice: 10.00, category: "Pastries", recommended: true, weeklyCapacity: 35, weeklyOrder: 0},
		{name: "Burnt Cheesecake", unitPrice: 12.00, category: "Pastries", recommended: false, weeklyCapacity: 35, weeklyOrder: 0},
		{name: "Red Velvet Carrot Cake", unitPrice: 8.00, category: "Pastries", recommended: true, weeklyCapacity: 35, weeklyOrder: 0},
		{name: "Baileys Chocolate Cake", unitPrice: 8.00, category: "Pastries", recommended: false, weeklyCapacity: 35, weeklyOrder: 0},
	}
)

//printItemList prints list of items
func printItemList() {
	fmt.Println(items)
}
