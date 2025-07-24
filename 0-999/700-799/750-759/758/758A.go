package main

import "fmt"

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	arr := make([]int, n)
	maxVal := 0
	for i := 0; i < n; i++ {
		fmt.Scan(&arr[i])
		if arr[i] > maxVal {
			maxVal = arr[i]
		}
	}
	sum := 0
	for _, v := range arr {
		sum += maxVal - v
	}
	fmt.Println(sum)
}
