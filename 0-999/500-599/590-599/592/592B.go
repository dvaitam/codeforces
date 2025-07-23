package main

import (
	"fmt"
)

func main() {
	var n int64
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	if n < 2 {
		fmt.Println(0)
		return
	}
	ans := (n - 2) * (n - 2)
	fmt.Println(ans)
}
