package main

import "fmt"

func main() {
	var n, k int64
	if _, err := fmt.Scan(&n, &k); err != nil {
		return
	}
	for i := int64(0); i < k; i++ {
		if n == 1 {
			break
		}
		res := n
		temp := n
		for j := int64(2); j*j <= temp; j++ {
			if temp%j == 0 {
				for temp%j == 0 {
					temp /= j
				}
				res -= res / j
			}
		}
		if temp > 1 {
			res -= res / temp
		}
		n = res
	}
	fmt.Println(n % 1000000007)
}