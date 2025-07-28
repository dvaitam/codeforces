package main

import "fmt"

func main() {
	var t int
	fmt.Scan(&t)
	for ; t > 0; t-- {
		var n int
		fmt.Scan(&n)
		arr := []int{}
		fmt.Println(arr[n]) // panic index out of range
	}
}
