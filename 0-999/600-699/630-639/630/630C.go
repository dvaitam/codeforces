package main

import (
	"fmt"
)

func main() {
	var n int64
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	ans := (int64(1) << (n + 1)) - 2
	fmt.Println(ans)
}
