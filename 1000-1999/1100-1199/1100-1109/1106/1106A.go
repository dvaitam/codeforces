package main

import (
	"fmt"
)

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	s := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&s[i])
	}
	ans := 0
	for i := 1; i < n-1; i++ {
		for j := 1; j < n-1; j++ {
			if s[i][j] == 'X' &&
				s[i-1][j-1] == 'X' &&
				s[i-1][j+1] == 'X' &&
				s[i+1][j-1] == 'X' &&
				s[i+1][j+1] == 'X' {
				ans++
			}
		}
	}
	fmt.Println(ans)
}
