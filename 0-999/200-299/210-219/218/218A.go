package main

import "fmt"

func main() {
	var n, k int
	if _, err := fmt.Scan(&n, &k); err != nil {
		return
	}
	total := 2*n + 1
	a := make([]int, total)
	for i := 0; i < total; i++ {
		fmt.Scan(&a[i])
	}
	for i := 0; i < total; i++ {
		if i%2 == 1 && k > 0 && a[i]-1 > a[i-1] && a[i]-1 > a[i+1] {
			a[i]--
			k--
		}
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(a[i])
	}
	fmt.Println()
}
