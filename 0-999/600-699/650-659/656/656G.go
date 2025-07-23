package main

import (
	"fmt"
)

func main() {
	var F, I, T int
	if _, err := fmt.Scan(&F, &I, &T); err != nil {
		return
	}
	counts := make([]int, I)
	for i := 0; i < F; i++ {
		var s string
		fmt.Scan(&s)
		for j, ch := range s {
			if ch == 'Y' {
				counts[j]++
			}
		}
	}
	ans := 0
	for j := 0; j < I; j++ {
		if counts[j] >= T {
			ans++
		}
	}
	fmt.Println(ans)
}
