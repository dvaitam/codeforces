package main

import "fmt"

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	var a, b int
	if _, err := fmt.Scan(&a, &b); err != nil {
		return
	}
	fmt.Println(gcd(a, b))
}
