package main

import "fmt"

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	result := 3*n*(n+1) + 1
	fmt.Println(result)
}
