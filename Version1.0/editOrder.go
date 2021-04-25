package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

func viewOrEditMenu() error {
	fmt.Println("Would you like to view or edit a current booking?")
	fmt.Println("To view an existing booking, please enter 1.")
	fmt.Println("To edit an existing booking, please enter 2.")
	var choice int
	fmt.Scanln(&choice)
	if choice == 1 {
		orderList.viewOrder()
	} else if choice == 2 {
		orderList.editOrder()
	} else {
		return errors.New("error: you did not select options 1 or 2")
	}
	return nil
}

func (p *orderLL) viewOrder() error {
	var bookingNumber int64
	fmt.Println("Enter booking number")
	fmt.Scanln(&bookingNumber)

	if bookingNumber > orderNumber {
		return errors.New("booking number not yet created")
	}

	ptr := p.head
	for i := 0; i < p.size; i++ {
		if ptr.order.orderNum == bookingNumber {
			fmt.Printf("These are the order details for booking number %v: \n", ptr.order.orderNum)
			fmt.Printf("Recipient name: %v\n", ptr.order.name)
			fmt.Printf("Recipient address: %v\n", ptr.order.address)
			fmt.Printf("Delivery day: %v\n", intToDay(ptr.order.deliveryDay))
			for _, s := range ptr.order.shoppingCart {
				fmt.Printf("Quantity %v of %v was ordered.\n", s.qty, s.orderItem)
			}
			break
		}
		ptr = ptr.next
	}
	if bookingNumber > orderNumber {
		fmt.Println("Booking number not found!")
	}
	return nil
}

//EditOrder allows editing of booking
func (p *orderLL) editOrder() error {
	defer func() { //to handle possible panic situation
		if err := recover(); err != nil {
			fmt.Println("Oops, panic occurred:", err)
		}
	}()
	var bookingNumber int64
	var nam string
	var add string
	var dday int
	var amt float64
	var wg sync.WaitGroup
	fmt.Println("Enter booking number")
	fmt.Scanln(&bookingNumber)
	ptr := p.head

	if bookingNumber > orderNumber {
		return errors.New("booking number not yet created")
	}

	for i := 0; i < p.size; i++ {
		if ptr.order.orderNum == bookingNumber {
			fmt.Printf("These are the order details for %v: \n", ptr.order.orderNum)
			fmt.Printf("Name: %v\n", ptr.order.name)
			fmt.Printf("Address: %v\n", ptr.order.address)
			fmt.Printf("Delivery day: %v\n", intToDay(ptr.order.deliveryDay))
			for j := 0; j < len(ptr.order.shoppingCart); j++ {
				fmt.Printf("Quantity %v of item %v has been ordered.\n", ptr.order.shoppingCart[j].qty, ptr.order.shoppingCart[j].orderItem)
			}
			fmt.Println("Please re-enter details to edit your booking.")
			fmt.Println("==============================================")
			fmt.Println("Enter Name:")
			fmt.Println("Press enter for no change.")
			reader := bufio.NewReader(os.Stdin) //create new reader, assuming bufio imported
			nam, _ = reader.ReadString('\n')
			nam = strings.TrimRight(nam, "\n") // this removes the \n at end of scan function
			fmt.Println("Enter Address:")
			fmt.Println("Press for no change.")
			add, _ = reader.ReadString('\n')
			add = strings.TrimRight(add, "\n") // this removes the \n at end of scan function
			dday = scanAvailableDay()
			shoppingOrder := make(chan []order)
			go generateShoppingCart(shoppingOrder)
			sc := <-shoppingOrder
			if nam == "" {
				nam = ptr.order.name
				ptr.order.name = nam
			} else {
				ptr.order.name = nam
			}

			if add == "" {
				add = ptr.order.address
				ptr.order.address = add
			} else {
				ptr.order.address = add
			}

			if dday == 0 {
				dday = ptr.order.deliveryDay
				ptr.order.deliveryDay = dday
			} else {
				ptr.order.deliveryDay = dday
			}
			ptr.order.shoppingCart = sc
			amt = calculateAmount(sc)
			ptr.order.amount = amt
			break
		}
		ptr = ptr.next

	}
	//orderList.printAllOrderNodes()

	func() {
		wg.Add(2)
		go func() {
			defer wg.Done()
			updateWeeklySchedule(orderList)
		}()
		go func() {
			defer wg.Done()
			updateWeeklyOrder(orderList)
		}()
	}()
	wg.Wait()
	//fmt.Println(weeklySchedule)
	return nil
}
