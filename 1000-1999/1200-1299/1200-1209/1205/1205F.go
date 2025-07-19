package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	N = 100
	K = N * (N + 1) / 2
)

var dp [N + 1][K + 1]int

func initDP() {
	dp[1][1] = 1
	for n := 1; n <= N; n++ {
		maxK := n * (n + 1) / 2
		for k := n; k <= maxK; k++ {
			if dp[n][k] != 0 {
				// positive transitions
				for c := 2; n+c-1 <= N; c++ {
					n2 := n + c - 1
					k2 := k + c*(c+1)/2 - 1
					if dp[n2][k2] == 0 {
						dp[n2][k2] = c
					}
				}
				// negative transitions
				for c := 4; n+c-1 <= N; c++ {
					n2 := n + c - 1
					k2 := k + c
					if dp[n2][k2] == 0 {
						dp[n2][k2] = -c
					}
				}
			}
		}
	}
}

func trace(aa []int, n, k int) {
	if n == 1 {
		aa[0] = 1
		return
	}
	if dp[n][k] > 0 {
		c := dp[n][k]
		n2 := n - c + 1
		k2 := k - c*(c+1)/2 + 1
		trace(aa, n2, k2)
		// reverse first n2 elements
		for i, j := 0, n2-1; i < j; i, j = i+1, j-1 {
			aa[i], aa[j] = aa[j], aa[i]
		}
		// fill tail
		for i := n2; i < n; i++ {
			aa[i] = i + 1
		}
	} else {
		c := -dp[n][k]
		n2 := n - c + 1
		k2 := k - c
		trace(aa, n2, k2)
		// increment prefix
		for i := 0; i < n2; i++ {
			aa[i]++
		}
		idx := n2
		if c%2 == 0 {
			for a := n2 + 3; a <= n; a += 2 {
				aa[idx] = a
				idx++
			}
			aa[idx] = 1
			idx++
			for a := n2 + 2; a <= n; a += 2 {
				aa[idx] = a
				idx++
			}
		} else {
			aMid := n2 + c/2
			for a := n2 + 3; a < n; a += 2 {
				v := a
				if a >= aMid {
					v++
				}
				aa[idx] = v
				idx++
			}
			aa[idx] = aMid
			idx++
			aa[idx] = 1
			idx++
			for a := n2 + 2; a < n; a += 2 {
				v := a
				if a >= aMid {
					v++
				}
				aa[idx] = v
				idx++
			}
		}
	}
}

func main() {
	initDP()
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}
	aa := make([]int, N)
	for ; q > 0; q-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		if n <= N && k <= K && dp[n][k] != 0 {
			trace(aa, n, k)
			fmt.Fprintln(writer, "YES")
			for i := 0; i < n; i++ {
				if i > 0 {
					writer.WriteByte(' ')
				}
				fmt.Fprint(writer, aa[i])
			}
			fmt.Fprintln(writer)
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
