package main

import (
	"errors"
	"fmt"
	"sync"
)

func (p *orderLL) deleteOrder() error { //function deletes order based on input booking number
	defer func() { //to handle possible panic situation
		if err := recover(); err != nil {
			fmt.Println("Oops, panic occurred:", err)
		}
	}()
	userAuthenticated := authenticateAdmin()
	if userAuthenticated == false {
		fmt.Println("Wrong password. Access Denied. returning to main menu")
		return nil
	}
	var bookingNumber int64
	fmt.Println("Which booking number would you like to delete? Press enter to cancel.")
	fmt.Scanln(&bookingNumber)
	if bookingNumber == 0 {
		fmt.Println("Cancelling delete order, returning to main menu.")
		return nil
	}
	ptr := p.head

	if ptr == nil {
		return errors.New("empty order list")
	}

	for ptr.next.order.orderNum != bookingNumber {
		ptr = ptr.next
		if ptr.next == nil {
			return errors.New("booking number does not exist")
		}
	}
	fmt.Println(ptr.next.order, "deleting in progress..")
	ptr.next = ptr.next.next
	p.size--
	fmt.Printf("Order %v has been deleted.\n", bookingNumber)
	var wg2 sync.WaitGroup
	func() {
		wg2.Add(2)
		go func() {
			defer wg2.Done()
			updateWeeklySchedule(orderList)
		}()
		go func() {
			defer wg2.Done()
			updateWeeklyOrder(orderList)
		}()
		wg2.Wait()
		fmt.Println(weeklySchedule)
	}()
	//orderList.printAllOrderNodes() // this is used for program check purpose

	return nil

}
