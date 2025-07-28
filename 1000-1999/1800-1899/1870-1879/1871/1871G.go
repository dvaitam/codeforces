package main

import "fmt"

func reverse(x int) int {
	neg := false
	if x < 0 {
		neg = true
		x = -x
	}
	res := 0
	for x > 0 {
		res = res*10 + x%10
		x /= 10
	}
	if neg {
		res = -res
	}
	return res
}

func main() {
	var x int
	fmt.Scan(&x)
	fmt.Println(reverse(x))
}
