package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func authenticateAdmin() bool {
	var password string = "BestBakeryInSingapore"
	fmt.Println("Please verify Admin user by keying in the admin password.")
	var userPass string
	reader := bufio.NewReader(os.Stdin) //create new reader, assuming bufio imported
	userPass, _ = reader.ReadString('\n')
	userPass = strings.TrimRight(userPass, "\n") // this removes the \n at end of scan function
	if userPass == password {
		return true
	}
	return false
}
