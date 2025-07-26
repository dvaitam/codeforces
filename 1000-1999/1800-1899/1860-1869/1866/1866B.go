package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int = 998244353

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	A := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &A[i])
	}
	B := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &B[i])
	}

	var m int
	fmt.Fscan(reader, &m)
	C := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &C[i])
	}
	D := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &D[i])
	}

	ans := 1
	i, j := 0, 0
	for i < n || j < m {
		if i < n && (j >= m || A[i] < C[j]) {
			ans = ans * 2 % mod
			i++
		} else if j < m && (i >= n || C[j] < A[i]) {
			ans = 0
			break
		} else {
			if D[j] > B[i] {
				ans = 0
				break
			}
			if D[j] < B[i] {
				ans = ans * 2 % mod
			}
			i++
			j++
		}
	}

	fmt.Fprintln(writer, ans)
}
