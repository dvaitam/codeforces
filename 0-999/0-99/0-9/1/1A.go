package main

import "fmt"

func main() {
	var n, m, a int64
	if _, err := fmt.Scan(&n, &m, &a); err != nil {
		return
	}
	// calculate number of flags along each dimension using ceiling division
	res := ((n-1)/a + 1) * ((m-1)/a + 1)
	fmt.Print(res)
}
