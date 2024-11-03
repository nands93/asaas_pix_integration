package main

import (
	"fmt"
)

func main() {

	var operation int
	var program bool = true

	for program {
		fmt.Println("1 - GENERATE PIX QR CODE\n2 - REQUEST A CASH OUT\n3 - EXIT")
		fmt.Print("Select the operation: ")
		fmt.Scan(&operation)
		if operation == 1 {
			fmt.Println(create_qr_code())
		} else if operation == 2 {
			fmt.Println("OPERATION 2")
		} else if operation == 3 {
			program = false
		} else {
			fmt.Println("Please, select the right operation")
			continue
		}
	}
}
