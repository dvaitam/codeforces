package main

import "fmt"

func main() {
	var n, k int64
	if _, err := fmt.Scan(&n, &k); err != nil {
		return
	}
	result := (n/k + 1) * k
	fmt.Println(result)
}
