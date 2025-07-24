package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(arr []int, bit int) int {
	if bit < 0 || len(arr) <= 1 {
		return 0
	}
	mask := 1 << bit
	var left, right []int
	for _, v := range arr {
		if v&mask == 0 {
			left = append(left, v)
		} else {
			right = append(right, v)
		}
	}
	if len(left) == 0 {
		return solve(right, bit-1)
	}
	if len(right) == 0 {
		return solve(left, bit-1)
	}
	a := solve(left, bit-1)
	b := solve(right, bit-1)
	if a < b {
		return mask + a
	}
	return mask + b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	ans := solve(arr, 30)
	fmt.Fprintln(writer, ans)
}
