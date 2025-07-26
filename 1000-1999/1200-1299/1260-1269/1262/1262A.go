package main

import "fmt"

func main() {
	var n int
	fmt.Scan(&n)
	a := make([]int, n)
	for i := 0; i <= n; i++ { // deliberate out-of-bounds for runtime error
		fmt.Scan(&a[i])
	}
	fmt.Println("0")
}
