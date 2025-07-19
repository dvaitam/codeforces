package main

import (
	"fmt"
)

func main() {
	var n, m, k float64
	for {
		if _, err := fmt.Scan(&n, &m, &k); err != nil {
			break
		}
		if n+k < m {
			fmt.Println("0")
		} else {
			ans := 1.0
			// compute product for i from 0 to k inclusive
			for i := 0; i <= int(k); i++ {
				ans *= (m - float64(i)) / (n + 1 + float64(i))
			}
			res := 1.0 - ans
			fmt.Printf("%.6f\n", res)
		}
	}
}
