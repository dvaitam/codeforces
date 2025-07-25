package main

import "fmt"

func main() {
	var n int64
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	fmt.Println(n * (n + 1) / 2)
}
