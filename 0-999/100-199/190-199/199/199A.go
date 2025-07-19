package main

import "fmt"

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	// output n, 0, 0
	fmt.Printf("%d %d %d", n, 0, 0)
}
