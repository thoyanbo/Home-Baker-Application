package main

import "fmt"

//bookingDay keeps track of max and current orders status in the days of the week
type bookingDay struct {
	weekday       int
	maxOrders     int
	currentOrders int
}

//weeklySchedule is the global variable that tracks the max orders of the day in a week.
var weeklySchedule = []bookingDay{
	{weekday: 1, maxOrders: 8, currentOrders: 0},
	{weekday: 2, maxOrders: 8, currentOrders: 0},
	{weekday: 3, maxOrders: 8, currentOrders: 0},
	{weekday: 4, maxOrders: 8, currentOrders: 0},
	{weekday: 5, maxOrders: 8, currentOrders: 0},
	{weekday: 6, maxOrders: 5, currentOrders: 0},
	{weekday: 7, maxOrders: 5, currentOrders: 0},
}

//BookingSlice is created to create method
type bookingSlice []bookingDay

//IntToDay takes in an integer and returns the corresponding day string value of the week.
func intToDay(day int) string {
	if day == 1 {
		return "Monday"
	} else if day == 2 {
		return "Tuesday"
	} else if day == 3 {
		return "Wednesday"
	} else if day == 4 {
		return "Thursday"
	} else if day == 5 {
		return "Friday"
	} else if day == 6 {
		return "Saturday"
	} else if day == 7 {
		return "Sunday"
	} else {
		return "Error! You did not select a weekday."
	}
}

//PrintCurrentSchedule prints the current schedule
func printCurrentSchedule() {
	fmt.Println(weeklySchedule)
}

//SeeAvailableDays is a function to see days available for booking
func SeeAvailableDays() {
	for n := 0; n < len(weeklySchedule); n++ {
		if weeklySchedule[n].currentOrders < weeklySchedule[n].maxOrders {
			remainingOrders := weeklySchedule[n].maxOrders - weeklySchedule[n].currentOrders
			fmt.Println(intToDay(n+1), "is available for booking.", remainingOrders, "orders still available for taking.")
		}
	}
}
