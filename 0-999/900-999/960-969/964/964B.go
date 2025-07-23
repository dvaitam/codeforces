package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, A, B, C, T int
	if _, err := fmt.Fscan(reader, &n, &A, &B, &C, &T); err != nil {
		return
	}
	times := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &times[i])
	}

	if C <= B {
		// Best to read messages immediately when they arrive
		fmt.Println(int64(n) * int64(A))
		return
	}

	var sumDiff int64
	for _, t := range times {
		sumDiff += int64(T - t)
	}
	result := int64(n)*int64(A) + int64(C-B)*sumDiff
	fmt.Println(result)
}
