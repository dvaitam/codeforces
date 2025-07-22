package main

import "fmt"

func main() {
	var a int
	if _, err := fmt.Scan(&a); err != nil {
		return
	}
	if a%2 == 1 {
		fmt.Print(1)
	} else {
		fmt.Print(0)
	}
}
