package main

import "fmt"

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
func main() {
	var a, b int64
	if _, err := fmt.Scan(&a, &b); err != nil {
		return
	}
	fmt.Println(gcd(a, b))
}
