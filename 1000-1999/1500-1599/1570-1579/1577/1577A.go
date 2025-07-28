package main

import (
	"bufio"
	"fmt"
	"os"
)

// Intentional buggy solution for testing verifier runtime error detection.
func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	arr := make([]int, n)
	// Off-by-one bug: will panic when n > 0
	for i := 0; i <= n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	sum := 0
	for _, v := range arr {
		sum += v
	}
	fmt.Println(sum)
}
