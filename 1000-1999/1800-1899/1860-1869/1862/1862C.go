package main

import (
	"bufio"
	"fmt"
	"os"
)

func isSymmetric(a []int) bool {
	n := len(a)
	if a[0] != n {
		return false
	}
	j := n - 1
	for i := 1; i <= n; i++ {
		for j >= 0 && a[j] < i {
			j--
		}
		if j+1 != a[i-1] {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		if isSymmetric(a) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
