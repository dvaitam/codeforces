package main

import "fmt"

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	if n == 2 {
		fmt.Println(2)
	} else {
		fmt.Println(1)
	}
}
