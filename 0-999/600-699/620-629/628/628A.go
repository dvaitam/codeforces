package main

import "fmt"

func main() {
	var n, b, p int
	if _, err := fmt.Scan(&n, &b, &p); err != nil {
		return
	}
	bottles := (n - 1) * (2*b + 1)
	towels := n * p
	fmt.Printf("%d %d\n", bottles, towels)
}
