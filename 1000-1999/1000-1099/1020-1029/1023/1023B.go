package main

import "fmt"

func main() {
	var n, k int64
	fmt.Scan(&n, &k)
	if k > n {
		c := 2*n - k + 1
		if c < 0 {
			fmt.Println(0)
		} else {
			fmt.Println(c / 2)
		}
	} else {
		c := k - 1
		fmt.Println(c / 2)
	}
}
