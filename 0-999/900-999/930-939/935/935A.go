package main

import (
	"fmt"
)

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	ans := 0
	for l := 1; l < n; l++ {
		if n%l == 0 {
			ans++
		}
	}
	fmt.Println(ans)
}
