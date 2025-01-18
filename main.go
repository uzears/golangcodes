package main

import (
	"fmt"
)

// MYSQL

const LoginId string = "System"

func main() {

	var MenuCase int
	//var name string
	//var age int

	fmt.Println("Welcome to Salary management application\n")
	fmt.Println("1.Display Salary \n2.Calculate Salary \n3.Get Employee Details \n4.Add New Employee")
	fmt.Println("\nEnter Case Number")
	fmt.Scanln(&MenuCase)

	switch MenuCase {
	case 1:
		fmt.Println("Yet to be developed")
	case 2:
		fmt.Println("Yet to be developed")
	case 3:
		fmt.Println("Yet to be developed")
	case 4:
		fmt.Println("Yet to be developed")
	default:
		fmt.Println("Invalid input")
	}
}
