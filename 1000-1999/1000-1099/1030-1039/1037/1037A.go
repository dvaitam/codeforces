package main

import "fmt"

func main() {
	var n int64
	// Read the input number
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	// Find the smallest i such that n < 2^i
	for i := 0; i < 32; i++ {
		if n < 1<<i {
			fmt.Println(i)
			return
		}
	}
}
