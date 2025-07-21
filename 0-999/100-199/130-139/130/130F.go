package main

import "fmt"

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	first := true
	for i := 2; i <= n; i++ {
		for n%i == 0 {
			if !first {
				fmt.Print(" ")
			}
			fmt.Print(i)
			first = false
			n /= i
		}
	}
	fmt.Println()
}
