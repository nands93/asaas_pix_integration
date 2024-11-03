package main

import (
	"fmt"
)

func main() {

	var operation int

	fmt.Print("1 - GENERATE PIX QR CODE\n2 - REQUEST A CASH OUT\n")
	fmt.Print("Select the operation: ")
	fmt.Scan(&operation)

	if operation == 1 {
		fmt.Print(create_qr_code())
		fmt.Print("\n")
	}

}
