package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func preprocessNext(s string) [][]int {
	n := len(s)
	next := make([][]int, n+2)
	for i := range next {
		next[i] = []int{n, n, n}
	}
	for i := n - 1; i >= 0; i-- {
		for c := 0; c < 3; c++ {
			next[i][c] = next[i+1][c]
		}
		ch := s[i]
		if ch == '0' || ch == '?' {
			next[i][0] = i
		}
		if ch == '1' || ch == '?' {
			next[i][1] = i
		}
		next[i][2] = i
	}
	return next
}

func preprocessPrev(s string) [][]int {
	n := len(s)
	prev := make([][]int, n+2)
	for i := range prev {
		prev[i] = []int{-1, -1, -1}
	}
	for i := 0; i < n; i++ {
		for c := 0; c < 3; c++ {
			prev[i+1][c] = prev[i][c]
		}
		ch := s[i]
		if ch == '0' || ch == '?' {
			prev[i+1][0] = i
		}
		if ch == '1' || ch == '?' {
			prev[i+1][1] = i
		}
		prev[i+1][2] = i
	}
	return prev
}

func countWays(s string) int64 {
	n := len(s)
	next := preprocessNext(s)
	prev := preprocessPrev(s)

	pow2 := make([]int64, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = pow2[i-1] * 2 % MOD
	}

	ans := int64(0)
	for L := 1; L < n; L++ {
		needs0 := next[0][0] < L && next[n-L][1] < n
		needs1 := next[0][1] < L && next[n-L][0] < n

		if !needs0 {
			switch s[0] {
			case '0':
				needs0 = true
			case '1':
				needs1 = true
			}
		}
		if needs0 && needs1 {
			return 0
		}

		var minSuffixZero int
		if needs1 {
			minSuffixZero = n
		} else {
			minSuffixZero = next[n-L][0]
		}

		var maxPrefixChoose int
		if needs0 {
			maxPrefixChoose = -1
		} else {
			maxPrefixChoose = prev[L][0]
		}

		var minPrefixChoose int
		if next[0][0] >= L {
			minPrefixChoose = -1
		} else {
			minPrefixChoose = next[0][0]
		}

		var maxSuffixChoose int
		if prev[n][0] < n-L {
			maxSuffixChoose = n
		} else {
			maxSuffixChoose = prev[n][0]
		}

		var zeroPositions []int
		for i := minPrefixChoose; i >= 0; i-- {
			if s[i] == '?' {
				zeroPositions = append(zeroPositions, i)
			}
		}

		var needZeros int
		if minSuffixZero < n && maxPrefixChoose < minSuffixZero {
			needZeros = minSuffixZero - maxPrefixChoose
		} else {
			needZeros = 0
		}

		if needZeros <= len(zeroPositions) {
			ans = (ans + pow2[len(zeroPositions)-needZeros]) % MOD
		}
	}
	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		fmt.Fprintln(writer, countWays(s))
	}
}

