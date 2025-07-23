package main

import (
	"fmt"
)

func main() {
	var n int64
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	// (9n - 3) * 4^{n-3}
	pow := int64(1)
	for i := int64(0); i < n-3; i++ {
		pow *= 4
	}
	ans := (9*n - 3) * pow
	fmt.Print(ans)
}
