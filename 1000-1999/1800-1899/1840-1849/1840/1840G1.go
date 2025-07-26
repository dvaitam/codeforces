package main

import "fmt"

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	for i := 0; i < n; i++ {
		var x int
		fmt.Scan(&x)
	}
	fmt.Println(n)
}
