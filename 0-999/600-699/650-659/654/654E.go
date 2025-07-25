package main

import "fmt"

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	max := int64(-1 << 63)
	for i := 0; i < n; i++ {
		var x int64
		fmt.Scan(&x)
		if x > max {
			max = x
		}
	}
	fmt.Println(max)
}
