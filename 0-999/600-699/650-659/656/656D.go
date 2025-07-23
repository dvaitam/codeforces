package main

import "fmt"

func main() {
	var a int
	if _, err := fmt.Scan(&a); err != nil {
		return
	}
	fmt.Println(1)
}
