package main

import (
	"bufio"
	"fmt"
	"os"
)

func checkAscending(a []int64, c int64) bool {
	sum := a[0]
	for i := 1; i < len(a); i++ {
		need := int64(i+1) * c
		if sum+a[i] < need {
			return false
		}
		sum += a[i]
	}
	return true
}

func checkDescending(a []int64, c int64) bool {
	n := len(a)
	sum := a[n-1]
	for i := n - 2; i >= 0; i-- {
		need := int64(i+1) * int64(i+2) * c
		if sum+a[i] < need {
			return false
		}
		sum += a[i]
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(reader, &t)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for ; t > 0; t-- {
		var n int
		var c int64
		fmt.Fscan(reader, &n, &c)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		if checkAscending(a, c) || checkDescending(a, c) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
