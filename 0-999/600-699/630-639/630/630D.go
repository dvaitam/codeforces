package main

import "fmt"

func main() {
	var n int64
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	result := 1 + 3*n*(n+1)
	fmt.Println(result)
}
