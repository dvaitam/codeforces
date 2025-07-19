package main

import (
	"fmt"
)

func main() {
	var n, d int
	if _, err := fmt.Scan(&n, &d); err != nil {
		return
	}
	if n == 1 && d == 10 {
		fmt.Print(-1)
		return
	}
	if d == 10 {
		// One followed by n-1 zeros
		fmt.Print("1")
		for i := 1; i < n; i++ {
			fmt.Print("0")
		}
	} else {
		// Repeat digit d n times
		for i := 0; i < n; i++ {
			fmt.Print(d)
		}
	}
}
