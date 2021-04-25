package main

import (
	"errors"
	"fmt"
)

//order is used to store singular item orders into shopping cart
type order struct {
	orderItem string
	qty       float64
}

//orderInfo is used to store order information
type orderInfo struct {
	name         string
	address      string
	deliveryDay  int
	orderNum     int64
	shoppingCart []order
	amount       float64
}

//OrderNode is the linked list node that stores order data and pointer to next node
type orderNode struct {
	order orderInfo
	next  *orderNode
}

//OrderLL stores head and size of order linked list
type orderLL struct {
	head *orderNode
	size int
}

//addOrder allows addition of orders into order linkedlist
func (p *orderLL) addOrder(oi orderInfo) error {
	newNode := &orderNode{
		order: oi,
		next:  nil,
	}
	if p.head == nil {
		p.head = newNode
	} else {
		currentNode := p.head
		for currentNode.next != nil {
			currentNode = currentNode.next
		}
		currentNode.next = newNode
	}
	p.size++
	return nil
}

//printAllOrderNodes prints all nodes in order linked list
func (p *orderLL) printAllOrderNodes() error {
	if p.head == nil {
		return errors.New("the list is empty")
	}
	currentNode := p.head
	for currentNode != nil {
		fmt.Println(currentNode.order)
		currentNode = currentNode.next
	}
	return nil

}
