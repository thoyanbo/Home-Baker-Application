package main

func main() {
	preLoadOrders() // this is done to load some existing orders into the program
	updateWeeklySchedule(orderList)
	updateWeeklyOrder(orderList)
	printMenu()
	//orderList.printAllOrderNodes() // this is used to print and check all order nodes
}
