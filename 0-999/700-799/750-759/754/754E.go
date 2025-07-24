package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007
const base1 int64 = 911382323
const base2 int64 = 972663749

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(in, &n, &m)
	board := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &board[i])
	}
	var r, c int
	fmt.Fscan(in, &r, &c)
	pattern := make([]string, r)
	for i := 0; i < r; i++ {
		fmt.Fscan(in, &pattern[i])
	}

	N := n + r - 1
	M := m + c - 1
	pow1 := make([]int64, N+1)
	pow2 := make([]int64, M+1)
	pow1[0] = 1
	for i := 1; i <= N; i++ {
		pow1[i] = pow1[i-1] * base1 % mod
	}
	pow2[0] = 1
	for i := 1; i <= M; i++ {
		pow2[i] = pow2[i-1] * base2 % mod
	}

	ans := make([][]bool, n)
	for i := 0; i < n; i++ {
		ans[i] = make([]bool, m)
		for j := range ans[i] {
			ans[i][j] = true
		}
	}

	prefix := make([][]int64, N+1)
	for i := range prefix {
		prefix[i] = make([]int64, M+1)
	}

	for ch := byte('a'); ch <= 'z'; ch++ {
		// compute pattern hash for this letter
		var patHash int64
		for i := 0; i < r; i++ {
			row := pattern[i]
			for j := 0; j < c; j++ {
				if row[j] == ch {
					patHash = (patHash + pow1[i]*pow2[j]) % mod
				}
			}
		}
		if patHash == 0 {
			continue
		}

		// build prefix table for this letter
		for i := 0; i <= N; i++ {
			for j := range prefix[i] {
				prefix[i][j] = 0
			}
		}
		for i := 0; i < N; i++ {
			rowIdx := i % n
			prevRow := prefix[i]
			curRow := prefix[i+1]
			bRow := board[rowIdx]
			for j := 0; j < M; j++ {
				colIdx := j % m
				val := int64(0)
				if bRow[colIdx] == ch {
					val = pow1[i] * pow2[j] % mod
				}
				curRow[j+1] = (curRow[j] + prevRow[j+1] - prevRow[j] + val) % mod
				if curRow[j+1] < 0 {
					curRow[j+1] += mod
				}
			}
		}

		for i := 0; i < n; i++ {
			if patHash == 0 {
				break
			}
			base1Shift := pow1[i]
			for j := 0; j < m; j++ {
				if !ans[i][j] {
					continue
				}
				sub := prefix[i+r][j+c]
				sub -= prefix[i][j+c]
				if sub < 0 {
					sub += mod
				}
				sub -= prefix[i+r][j]
				if sub < 0 {
					sub += mod
				}
				sub += prefix[i][j]
				sub %= mod
				if sub < 0 {
					sub += mod
				}
				expected := patHash * base1Shift % mod * pow2[j] % mod
				if sub != expected {
					ans[i][j] = false
				}
			}
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for i := 0; i < n; i++ {
		buf := make([]byte, m)
		for j := 0; j < m; j++ {
			if ans[i][j] {
				buf[j] = '1'
			} else {
				buf[j] = '0'
			}
		}
		fmt.Fprintln(writer, string(buf))
	}
}
