package main

import (
	"fmt"
)

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	if n <= 3 {
		fmt.Println("NO")
		return
	}
	fmt.Println("YES")
	if n%2 == 0 {
		// Use pairs to make 1
		for i := n; i > 4; i -= 2 {
			fmt.Printf("%d - %d = 1\n", i, i-1)
		}
		// Multiply ones
		for i := n; i > 4; i -= 2 {
			fmt.Println("1 * 1 = 1")
		}
		fmt.Println("2 * 3 = 6")
		fmt.Println("6 * 4 = 24")
		fmt.Println("24 * 1 = 24")
	} else {
		// n is odd
		fmt.Println("5 - 1 = 4")
		fmt.Println("4 - 2 = 2")
		fmt.Println("2 * 3 = 6")
		fmt.Println("6 * 4 = 24")
		for i := n; i > 5; i -= 2 {
			fmt.Printf("%d - %d = 1\n", i, i-1)
			fmt.Println("24 * 1 = 24")
		}
	}
}
