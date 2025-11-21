package main

import "fmt"

func main() {
	var n, m int
	fmt.Scan(&n, &m)
	ans := make([]int64, n)
	ans[0] = 1
	for i := 1; i < n; i++ {
		ans[i] = ans[i-1] + int64(i-1)
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(ans[i])
	}
	fmt.Println()
}
