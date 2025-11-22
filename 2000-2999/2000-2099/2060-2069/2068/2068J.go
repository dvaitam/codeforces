package main

import (
	"bufio"
	"fmt"
	"os"
)

// interleave checks if string s can be formed by interleaving sequences a and b while
// preserving the order inside each sequence.
func interleave(s string, a, b []byte) bool {
	na, nb := len(a), len(b)
	dp := make([][]bool, na+1)
	for i := range dp {
		dp[i] = make([]bool, nb+1)
	}
	dp[0][0] = true

	for i := 0; i <= na; i++ {
		for j := 0; j <= nb; j++ {
			if !dp[i][j] {
				continue
			}
			idx := i + j
			if idx >= len(s) {
				continue
			}
			if i < na && s[idx] == a[i] {
				dp[i+1][j] = true
			}
			if j < nb && s[idx] == b[j] {
				dp[i][j+1] = true
			}
		}
	}
	return dp[na][nb]
}

func canArrange(s string, n int) bool {
	// x = number of positions assigned to the first subset among the first n bottles.
	for x := 0; x <= n; x++ {
		need1 := make([]byte, n) // required colors for positions of the first subset (p1_i)
		need2 := make([]byte, n) // required colors for positions of the second subset (p2_i)

		for i := 1; i <= n; i++ {
			switch {
			case i <= x && i <= n-x: // both elements of pair i lie in the prefix
				need1[i-1] = 'W'
				need2[i-1] = 'W'
			case i > x && i > n-x: // both elements lie in the suffix
				need1[i-1] = 'R'
				need2[i-1] = 'R'
			case i <= x: // p1 in prefix, p2 in suffix
				need1[i-1] = 'R' // moves to suffix
				need2[i-1] = 'W' // moves to prefix
			default: // p1 in suffix, p2 in prefix
				need1[i-1] = 'W'
				need2[i-1] = 'R'
			}
		}

		// First n characters must interleave need1[:x] and need2[:n-x];
		// remaining n characters must interleave the rest.
		if !interleave(s[:n], need1[:x], need2[:n-x]) {
			continue
		}
		if !interleave(s[n:], need1[x:], need2[n-x:]) {
			continue
		}
		return true
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)

		if canArrange(s, n) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
