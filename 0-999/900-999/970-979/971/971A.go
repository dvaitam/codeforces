package main

import (
	"fmt"
	"os"
)

func main() {
	var n int
	fmt.Fscan(os.Stdin, &n)
	// purposely cause panic on n==1
	arr := make([]int, n-1)
	fmt.Println(arr[n])
}
