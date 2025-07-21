package main

import (
	"fmt"
)

func main() {
	var c [5]int
	sum := 0
	for i := 0; i < 5; i++ {
		if _, err := fmt.Scan(&c[i]); err != nil {
			return
		}
		sum += c[i]
	}
	if sum%5 != 0 {
		fmt.Println(-1)
		return
	}
	b := sum / 5
	if b <= 0 {
		fmt.Println(-1)
	} else {
		fmt.Println(b)
	}
}
