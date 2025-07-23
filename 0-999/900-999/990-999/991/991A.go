package main

import "fmt"

func main() {
	var a, b, c, n int
	if _, err := fmt.Scan(&a, &b, &c, &n); err != nil {
		return
	}
	visited := a + b - c
	if a < c || b < c || visited >= n || a > n || b > n {
		fmt.Println(-1)
		return
	}
	fmt.Println(n - visited)
}
