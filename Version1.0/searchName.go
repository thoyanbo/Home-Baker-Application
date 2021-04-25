package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func searchName() error {
	defer func() { //to handle possible panic situation
		if err := recover(); err != nil {
			fmt.Println("Oops, panic occurred:", err)
		}
	}()
	fmt.Println("Which item are you searching for?")
	reader := bufio.NewReader(os.Stdin) //create new reader, assuming bufio imported
	item, _ := reader.ReadString('\n')
	item = strings.TrimRight(item, "\n") // this removes the \n on end of scan function for comparison later
	if item == "" {
		return errors.New("You did not enter an item name")
	}
	var i int
	var itemFound bool = false
	for i = 0; i < len(items); i++ {
		if item == items[i].name {
			itemFound = true
			fmt.Println("Item found!")
			fmt.Printf("Item: %v     Price: %v:     Category: %v.\n", items[i].name, items[i].unitPrice, items[i].category)
			break
		}
	}
	if itemFound == false {
		fmt.Println("Error!", item, "cannot be found.")
	}
	return nil
}
